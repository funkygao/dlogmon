package dlog

import (
    "bufio"
    "fmt"
    "io"
    "kx/stream"
    T "kx/trace"
    "log"
    "strings"
)

const (
    LZOP_CMD          = "lzop"
    LZOP_OPTION       = "-dcf"
    EOL               = '\n'
    DLOG_BASE_DIR     = "/kx/dlog/"
    SAMPLER_HOST      = "100.123"
    FLAG_TIMESPAN_SEP = "-"
)

// Any kind of things
type Any interface{}

// Worker struct method signatures
type IWorker interface {
    Run(IWorker) // IWorker param for dynamic polymorphism
    IsLineValid(string) bool
    ExtractLineInfo(string) Any
    Running() bool
}

// An executor for 1 dlog file
type Worker struct {
    running   bool
    filename  string // dlog filename
    mapReader *bufio.Reader
    mapWriter *bufio.Writer
    *log.Logger
    manager *Manager
}

// Worker constructor signature
type WorkerConstructor func(*Manager, string) IWorker

// Request object for a line
type Request struct {
    http_method, uri, rid string
}

// Scan result of raw lines and valid lines
type ScanResult struct {
    RawLines, ValidLines int
}

// Printable Worker
func (this *Worker) String() string {
    return fmt.Sprintf("Worker{filename: %s, option: %#v}", this.filename, this.manager.option)
}

// Is this dlog executor running?
func (this *Worker) Running() bool {
    return this.running
}

func (this *Worker) initMapper() *stream.Stream {
    defer T.Un(T.Trace("initMapper"))

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
func (this *Worker) Run(dlog IWorker) {
    defer T.Un(T.Trace("Run"))

    this.Println(this.filename, "start scanning...")

    if this.manager.option.debug {
        fmt.Println("\n", this, "\n")
    }

    if mapper := this.initMapper(); mapper != nil {
        defer mapper.Close()
    }

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

        if !dlog.IsLineValid(line) {
            continue
        }

        validLines++

        // extract parsed info from this line and report to manager
        if x := dlog.ExtractLineInfo(line); x != nil {
            this.manager.collectLineMeta(x)
        }
    }

    this.manager.collectExecutorSummary(rawLines, validLines)
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
