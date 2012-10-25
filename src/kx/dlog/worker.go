package dlog

import (
    "fmt"
    "io"
    "kx/mr"
    "kx/stream"
    T "kx/trace"
    "log"
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

// Is this dlog worker running?
func (this *Worker) Running() bool {
    return this.running
}

// How many lines in this worker file
// TODO accurate line count instead of const
func (this Worker) TotalLines() int {
    return LINES_PER_FILE
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
func (this *Worker) SafeRun(chOutProgress chan<- int, chOutMap chan<- interface{}, chOutWorker chan<- WorkerResult) {
    defer T.Un(T.Trace(""))

    // recover to make this worker safe for other workers
    defer func() {
        if err := recover(); err != nil {
            this.Println("panic recovered:", err)
        }
    }()

    if this.manager.option.debug {
        fmt.Println(this)
    }

    if mapper := this.initMapper(); mapper != nil {
        defer mapper.Close()
    }

    this.run(chOutProgress, chOutMap, chOutWorker)
}

func (this *Worker) run(chOutProgress chan<- int, chOutMap chan<- interface{}, chOutWorker chan<- WorkerResult) {
    this.running = true

    // invoke transform goroutine to transform k=>v into k=>[]v
    tranResult := make(chan mr.TransformData)
    tranIn := make(chan interface{})
    go this.transform(tranIn, tranResult)

    input := stream.NewStream(LZOP_CMD, LZOP_OPTION, this.filename)
    input.Open()
    defer input.Close()

    this.Println(this.BaseName(), this.name, LZOP_CMD, "exec opened")

    inputReader := input.Reader()
    var rawLines, validLines int
    for {
        line, err := inputReader.ReadString(EOL)
        if err != nil {
            if err != io.EOF {
                log.Fatal(err)
            }

            break
        }

        rawLines++
        if chOutProgress != nil && rawLines % PROGRESS_LINES_STEP == 0 {
            // report progress
            chOutProgress <- PROGRESS_LINES_STEP
        }

        if !this.self.IsLineValid(line) {
            continue
        }

        validLines++

        // run map for this line
        this.self.Map(line, tranIn)
    }

    chOutWorker <- newWorkerResult(rawLines, validLines)
    this.Printf("%s %s lines parsed: %d/%d\n", this.BaseName(), this.name, validLines, rawLines)

    // transform feed done, must close before get data from tranResult
    close(tranIn)

    var r mr.TransformData = <-tranResult
    this.Println(this.BaseName(), this.name, "transformed")

    if this.self.Combiner() != nil {
        // run combiner
        for k, v := range r {
            r[k] = []float64{this.self.Combiner()(v)}
        }

        this.Println(this.BaseName(), this.name, "combined")
    }

    // output the transform result
    chOutMap <- r

    this.running = false
}

func (this *Worker) transform(in <-chan interface{}, out chan<- mr.TransformData) {
    r := mr.NewTransformData()
    for x := range in {
        for k, v := range x.(mr.MapData) {
            r.Append(k, v)
        }
    }

    out <- r
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
