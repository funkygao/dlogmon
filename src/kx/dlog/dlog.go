package dlog

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "os/exec"
    "sync"
)

const (
    LZOP_CMD = "lzop"
    LZOP_OPTION = "-dcf"
    EOL = '\n'
    DLOG_BASE_DIR = "/kx/dlog/"
)

// request object for a line
type Request struct {
    http_method, uri, rid string
}

// dlog interface
type IDlogExecutor interface {
    ScanLines(IDlogExecutor) // IDlogExecutor param for dynamic polymorphism
    IsLineValid(string) bool
    OperateLine(string)
}

// an executor for 1 dlog file
type Dlog struct {
    options *Options
    filename string // dlog filename
    chLines chan int // lines parsed channel
    lock *sync.Mutex
}

// java's toString()
func (this *Dlog) String() string {
    return fmt.Sprintf("Dlog{filename: %s, options: %#v}", this.filename, this.options)
}

// the main loop
func (this *Dlog) ScanLines(dlog IDlogExecutor) {
    run := exec.Command(LZOP_CMD, LZOP_OPTION, this.filename)
    out, err := run.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }
    if err := run.Start(); err != nil {
        log.Fatal(err)
    }

    if this.options.debug {
        fmt.Println(this)
    }

    if this.options.mapper != "" {
        //mapper := exec.Command(this.options.mapper)
    }

    inputReader := bufio.NewReader(out)
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

    if err := run.Wait(); err != nil {
        log.Fatal(err)
    }

    this.chLines <- lineCount
}

// base
func (this *Dlog) IsLineValid(line string) bool {
    return true
}

// base
func (this *Dlog) OperateLine(line string) {
    // do nothing
}

