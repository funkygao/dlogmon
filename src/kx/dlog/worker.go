package dlog

import (
    "fmt"
    "io"
    "kx/stream"
    T "kx/trace"
    "log"
    "strings"
)

// Printable Worker
func (this *Worker) String() string {
    return fmt.Sprintf("Worker{filename: %s, option: %#v}", this.filename, this.manager.option)
}

func newWorkerResult(rawLines, validLines int) WorkerResult {
    return WorkerResult{rawLines, validLines}
}

// Is this dlog worker running?
func (this *Worker) Running() bool {
    return this.running
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
func (this *Worker) SafeRun(worker IWorker, chOutLine chan<- Any, chOutWorker chan<- WorkerResult) {
    defer T.Un(T.Trace(""))

    // recover to make this worker safe for other workers
    defer func() {
        if err := recover(); err != nil {
            this.manager.Println("panic recovered:", err)
        }
    }()

    this.Println(this.filename, "start scanning...")

    if this.manager.option.debug {
        fmt.Println(this)
    }

    if mapper := this.initMapper(); mapper != nil {
        defer mapper.Close()
    }

    this.run(worker, chOutLine, chOutWorker)
}

func (this *Worker) run(worker IWorker, chOutLine chan<- Any, chOutWorker chan<- WorkerResult) {
    this.running = true

    input := stream.NewStream(LZOP_CMD, LZOP_OPTION, this.filename)
    input.Open()
    defer input.Close()

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

        if !worker.IsLineValid(line) {
            continue
        }

        validLines++

        // extract parsed info from this line and report to manager
        if lineMeta := worker.ExtractLineInfo(line); lineMeta != nil {
            chOutLine <- lineMeta
        }
    }

    chOutWorker <- newWorkerResult(rawLines, validLines)

    this.running = false
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
func (this *Worker) ExtractLineInfo(line string) Any {
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
