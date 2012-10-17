/*
Dlog base
*/
package dlog

import (
    "bufio"
    "fmt"
    "kx/stream"
    "io"
    "log"
    "strings"
    "sync"
)

const (
    LZOP_CMD = "lzop"
    LZOP_OPTION = "-dcf"
    EOL = '\n'
    DLOG_BASE_DIR = "/kx/dlog/"
    SAMPLER_HOST = "100.123"
)

type Any interface {}

type DlogConstructor func(string, chan int, *sync.Mutex, *Options) IDlogExecutor

// request object for a line
type Request struct {
    http_method, uri, rid string
}

// dlog interface
type IDlogExecutor interface {
    Run(IDlogExecutor) // IDlogExecutor param for dynamic polymorphism
    IsLineValid(string) bool
    OperateLine(string) Any
    Running() bool
    Progresser
}

type Progresser interface {
    Progress(int)
}

// an executor for 1 dlog file
type Dlog struct {
    running bool
    filename string // dlog filename
    chLines chan int // lines parsed channel
    lock *sync.Mutex
    options *Options
    *log.Logger
    mapReader *bufio.Reader
    mapWriter *bufio.Writer
}

// java's toString()
func (this *Dlog) String() string {
    return fmt.Sprintf("Dlog{filename: %s, options: %#v}", this.filename, this.options)
}

// Is this dlog executor running?
func (this *Dlog) Running() bool {
    return this.running
}

func (this *Dlog) init() {
    this.Println(this.filename, "start scanning...")

    if this.options.debug {
        fmt.Println("\n", this, "\n")
    }

    if this.options.mapper != "" {
        mapper := stream.NewStream(this.options.mapper)
        mapper.Open()
        defer mapper.Close()

        this.mapReader = mapper.Reader()
        this.mapWriter = mapper.Writer()
    }
}

// Scan each line and apply validator and parser
// Invoke mapper if neccessary
func (this *Dlog) Run(dlog IDlogExecutor) {
    this.init()

    this.running = true

    input := stream.NewStream(LZOP_CMD, LZOP_OPTION, this.filename)
    input.Open()
    defer input.Close()

    inputReader := input.Reader()
    lineCount := 0
    for {
        line, err := inputReader.ReadString(EOL)
        if err != nil {
            if err != io.EOF {
                log.Fatal(err)
            }

            break
        }

        lineCount += 1

        if !dlog.IsLineValid(line) {
            continue
        }

        // extract info from this line
        dlog.OperateLine(line)
    }

    this.chLines <- lineCount
    this.running = false
}

// base
func (this *Dlog) IsLineValid(line string) bool {
    if !strings.Contains(line, SAMPLER_HOST) {
        return false
    }
    return true
}

// base of valid line handler
func (this *Dlog) OperateLine(line string) Any {
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

// Mark this dlog executor done
func (this Dlog) Progress(finished int) {
    const BAR = "."

    total := len(this.options.files)
    fmt.Printf("[%*s%*s]\n", finished, BAR, total - finished, " ")
}
