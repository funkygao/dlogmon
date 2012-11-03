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
    "time"
)

func (this *Worker) String() string {
    return fmt.Sprintf("Worker{seq: %d, filename: %s, option: %#v}",
        this.seq, this.BaseFilename(), this.manager.option)
}

func (this Worker) Kind () string {
    return this.kind
}

func (this Worker) BaseFilename() string {
    return path.Base(this.filename)
}

// How many lines in this worker file
// TODO accurate line count instead of const
func (this Worker) TotalLines() int {
    return LINES_PER_FILE
}

// Common for worker constructors
func (this *Worker) init(manager *Manager, kind, filename string, seq uint16) {
    this.manager = manager
    this.kind = kind
    this.filename = filename
    this.seq = seq
    this.CreatedAt = time.Now()

    this.Logger = this.manager.Logger
}

func (this *Worker) initExternalMapper() *stream.Stream {
    defer T.Un(T.Trace(""))

    mapper := this.manager.option.mapper
    if mapper != "" {
        stream := stream.NewStream(stream.EXEC_PIPE, mapper)
        if err := stream.Open(); err != nil {
            this.Fatal(err)
        }

        this.Printf("external mapper stream opened: %s\n", mapper)

        this.mapReader = stream.Reader()
        this.mapWriter = stream.Writer()
        return stream
    }

    return nil
}

func (this *Worker) SafeRun(chOutProgress chan<- int, chOutMap chan<- mr.KeyValue, chOutWorker chan<- Worker) {
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

    if mapper := this.initExternalMapper(); mapper != nil {
        defer mapper.Close()
    }

    this.run(chOutProgress, chOutMap, chOutWorker)
}

// 每个worker向chan写入的次数：
// chOutProgress: N
// chOutMap: 1 for each parsed line
// chOutWorker: 1
func (this *Worker) run(chOutProgress chan<- int, chOutMap chan<- mr.KeyValue, chOutWorker chan<- Worker) {
    defer T.Un(T.Trace(""))

    this.StartAt = time.Now()

    var input *stream.Stream
    if this.manager.option.filemode {
        input = stream.NewStream(stream.PLAIN_FILE, this.filename)
    } else {
        input = stream.NewStream(stream.EXEC_PIPE, LZOP_CMD, LZOP_OPTION, this.filename)
    }
    input.Open()
    defer input.Close()

    this.Printf("%s worker[%d] opened %s, start to Map...\n", this.kind, this.seq, this.BaseFilename())

    for {
        line, err := input.Reader().ReadString(EOL)
        if err != nil {
            if err != io.EOF {
                log.Fatal(err)
            }

            break
        }

        this.RawLines ++
        if chOutProgress != nil && this.RawLines%PROGRESS_LINES_STEP == 0 {
            // report progress
            chOutProgress <- PROGRESS_LINES_STEP
        }

        if !this.self.IsLineValid(line) {
            continue
        }

        this.ValidLines++

        // run map for this line 
        // for pipe stream flush to work, we can't strip EOL
        this.self.Map(line, chOutMap)
    }
    this.EndAt = time.Now()

    chOutWorker <- *this
    this.Printf("%s worker[%d] %s done, parsed: %d/%d\n", this.kind, this.seq, this.BaseFilename(),
        this.ValidLines, this.RawLines)
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

// Only show top N final result
// Default show all(0)
func (this Worker) TopN() int {
    t, e := this.manager.Conf().Int(this.self.Kind(), "topN")
    if e == nil {
        return t
    }

    return 0
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
    if string(mapperLine) == "\n" {
        // empty result, just a EOL
        return nil
    }

    return mapperLine
}
