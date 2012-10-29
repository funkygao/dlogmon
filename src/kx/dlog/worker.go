package dlog

import (
	"fmt"
	"io"
	"kx/mr"
	"kx/stream"
	T "kx/trace"
	"log"
	"os"
	"path"
	"strings"
)

// Printable Worker
func (this *Worker) String() string {
	return fmt.Sprintf("Worker{filename: %s, option: %#v}", this.BaseName(), this.manager.option)
}

// Base of my dlog filename
func (this Worker) BaseName() string {
	return path.Base(this.filename)
}

func newWorkerResult(rawLines, validLines int) WorkerResult {
	return WorkerResult{rawLines, validLines}
}

// How many lines in this worker file
// TODO accurate line count instead of const
func (this Worker) TotalLines() int {
	return LINES_PER_FILE
}

// Common for worker constructors
func (this *Worker) init(manager *Manager, name, filename string, seq uint16) {
	this.manager = manager
	this.name = name
	this.filename = filename
	this.seq = seq

	this.Logger = this.manager.Logger
}

func (this *Worker) initMapper() *stream.Stream {
	defer T.Un(T.Trace(""))

	option := this.manager.option
	if option.mapper != "" {
		mapper := stream.NewStream(option.mapper)
		mapper.Open()

		this.mapReader = mapper.Reader()
		this.mapWriter = mapper.Writer()
		return mapper
	}

	return nil
}

// Scan each line of a dlog file and apply validator and parser.
// Invoke mapper if neccessary
func (this *Worker) SafeRun(chOutProgress chan<- int, chOutMap chan<- mr.KeyValues, chOutWorker chan<- WorkerResult) {
	defer T.Un(T.Trace(""))

	// recover to make this worker safe for other workers
	defer func() {
		if err := recover(); err != nil {
			this.Println("panic recovered:", err)
			panic(err)
		}
	}()

	if this.manager.option.debug {
		fmt.Fprintln(os.Stderr, this)
	}

	if mapper := this.initMapper(); mapper != nil {
		defer mapper.Close()
	}

	this.run(chOutProgress, chOutMap, chOutWorker)
}

// 每个worker向chan写入的次数：
// chOutProgress: N
// chOutMap: 1
// chOutWorker: 1
func (this *Worker) run(chOutProgress chan<- int, chOutMap chan<- mr.KeyValues, chOutWorker chan<- WorkerResult) {
	// invoke transform goroutine to transform k=>v into k=>[]v
	chKvs := make(chan mr.KeyValues)
	chKv := make(chan mr.KeyValue, LINE_CHAN_BUF)
	go this.transform(chKv, chKvs)

	var input *stream.Stream
	if this.manager.option.filemode {
		input = stream.NewStream(this.filename)
	} else {
		input = stream.NewStream(LZOP_CMD, LZOP_OPTION, this.filename)
	}
	input.Open()
	defer input.Close()

	this.Printf("%s worker[%d] opened %s\n", this.name, this.seq, this.BaseName())

	var rawLines, validLines int
	for {
		line, err := input.Reader().ReadString(EOL)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}

			break
		}

		rawLines++
		if chOutProgress != nil && rawLines%PROGRESS_LINES_STEP == 0 {
			// report progress
			chOutProgress <- PROGRESS_LINES_STEP
		}

		if !this.self.IsLineValid(line) {
			continue
		}

		validLines++

		// run map for this line with EOL stripped
		this.self.Map(line[:len(line)-1], chKv)
	}

	chOutWorker <- newWorkerResult(rawLines, validLines)
	this.Printf("%s worker[%d] %s parsed: %d/%d\n", this.name, this.seq, this.BaseName(), validLines, rawLines)

	// transform feed done, must close before get data from tranResult
	close(chKv)

	var kvs mr.KeyValues = <-chKvs
	this.Printf("%s worker[%d] %s transformed\n", this.name, this.seq, this.BaseName())

	// after the work has done it's job, run it's combiner as a whole of this worker
	if this.self.Combiner() != nil {
		for k, v := range kvs {
			kvs[k] = []interface{}{this.self.Combiner()(mr.ConvertAnySliceToFloat(v))} // [1]float64
		}

		this.Printf("%s worker[%d] %s local combined\n", this.name, this.seq, this.BaseName())
	}

	// output the transform result
	chOutMap <- kvs
}

func (this *Worker) transform(in <-chan mr.KeyValue, out chan<- mr.KeyValues) {
	kvs := mr.NewKeyValues()
	for x := range in {
		for k, v := range x {
			kvs.Append(k, v)
		}
	}

	out <- kvs
}

// My combiner func pointer
func (this *Worker) Combiner() mr.CombinerFunc {
	return this.combiner
}

// Is a line valid?
// Only when log is from sampler host will it reuturn true
func (this *Worker) IsLineValid(line string) bool {
	if !strings.Contains(line, SAMPLER_HOST) {
		return false
	}
	return true
}

// Base to extract meta info from a valid line string.
// If mapper specified, return the mapper output, else return nil
func (this *Worker) ExtractLineInfo(line string) interface{} {
	if this.mapReader == nil || this.mapWriter == nil {
		return nil
	}

	_, err := this.mapWriter.WriteString(line)
	this.mapWriter.Flush() // must flush, else script will not get this line
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	mapperLine, _ := this.mapReader.ReadString(EOL)
	return mapperLine
}

func (this Worker) Name() string {
	return this.name
}
