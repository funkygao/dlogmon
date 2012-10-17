/*
Options is the CLI options object.

Dlog stands for a single dlog file executor.
Each Dlog will run in it's own goroutine.

Dlog analyzer has many kinds(such as amf), which is only interested in 
some specific kind of info. 
So Dlog has many sub structs, which should implement
'IsLineValid' and [map/reduct | ExtractLineInfo].

Attention:
    For performance issue, IsLineValid must be implemented in main go program, while
    map/reduce can be any runnable script file, e.g python/php/ruby/nodejs, etc.

Manager is the manager of all dlog goroutines.
There will be a single manager in runtime.

amf is a kind of Dlog, which just parse 'AMF_SLOW' related log lines.
*/
package dlog

import (
    "bufio"
    "fmt"
    "kx/stream"
    "io"
    "log"
    "strings"
)

const (
    LZOP_CMD = "lzop"
    LZOP_OPTION = "-dcf"
    EOL = '\n'
    DLOG_BASE_DIR = "/kx/dlog/"
    SAMPLER_HOST = "100.123"
    FLAG_TIMESPAN_SEP = "-"
)

// Any kind of things
type Any interface {}

// An executor for 1 dlog file
type Dlog struct {
    running bool
    filename string // dlog filename
    mapReader *bufio.Reader
    mapWriter *bufio.Writer
    *log.Logger
    manager *Manager
}

// Dlog constructor signature
type DlogConstructor func(*Manager, string) IDlogExecutor

// Request object for a line
type Request struct {
    http_method, uri, rid string
}

// Dlog struct method signatures
type IDlogExecutor interface {
    Run(IDlogExecutor) // IDlogExecutor param for dynamic polymorphism
    IsLineValid(string) bool
    ExtractLineInfo(string) Any
    Running() bool
    Progresser
}

type Progresser interface {
    Progress(int)
}

// Scan result of raw lines and valid lines
type ScanResult struct {
    RawLines, ValidLines int
}

// Printable Dlog
func (this *Dlog) String() string {
    return fmt.Sprintf("Dlog{filename: %s, options: %#v}", this.filename, this.manager.options)
}

// Is this dlog executor running?
func (this *Dlog) Running() bool {
    return this.running
}

func (this *Dlog) initMapper() *stream.Stream {
    options := this.manager.options
    if options.mapper != "" {
        mapper := stream.NewStream(options.mapper)
        mapper.Open()

        this.mapReader = mapper.Reader()
        this.mapWriter = mapper.Writer()
        return mapper
    } 

    return nil
}

// Scan each line of a dlog file and apply validator and parser.
// Invoke mapper if neccessary
func (this *Dlog) Run(dlog IDlogExecutor) {
    this.Println(this.filename, "start scanning...")

    if this.manager.options.debug {
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

        rawLines ++

        if !dlog.IsLineValid(line) {
            continue
        }

        validLines ++

        // extract parsed info from this line
        if x := dlog.ExtractLineInfo(line); x != nil {
        }
    }

    this.manager.chFileScanResult <- ScanResult{rawLines, validLines}
    this.running = false
}

// Is a line valid?
// Only when log is from sampler host will it reuturn true
func (this *Dlog) IsLineValid(line string) bool {
    if !strings.Contains(line, SAMPLER_HOST) {
        return false
    }
    return true
}

// Base to extract meta info from a valid line string.
// If mapper specified, return the mapper output, else return nil
func (this *Dlog) ExtractLineInfo(line string) Any {
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

func (this Dlog) Progress(finished int) {
    const BAR = "."

    total := len(this.manager.options.files)
    fmt.Printf("[%*s%*s]\n", finished, BAR, total - finished, " ")
}

