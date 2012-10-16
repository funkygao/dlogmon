// line reader
package dlog

import (
    "fmt"
    "log"
    "strings"
    "strconv"
    "sync"
    
    "os/exec"
    "bufio"
    "io"
)

var lineValidatorRegexes = [][]string{
    {"AMF_SLOW", "100.123", "PHP.CDlog"}, // must exists
    {"Q=DLog.log"}}                       // must not exists

// AMF_SLOW dlog analyzer
type AmfDlog struct {
    Dlog
}

// the constructor
func NewAmfDlog(filename string, ch chan int, lock *sync.Mutex, options *Options) IDlogExecutor {
    this := new(AmfDlog)
    this.filename = filename
    this.options = options
    this.chLines = ch
    this.lock = lock

    return this
}

// a single line meta info
type amfRequest struct {
    Request
    class, method, args string
    time int16
}

// parse a line into meta info
// ret -> valid line?
func (this *amfRequest) ParseLine(line string) {
    // major parts seperated by space
    parts := strings.Split(line, " ")

    // uri related
    uriInfo := strings.Split(parts[5], "+")
    if len(uriInfo) < 3 {
        log.Fatal(line)
    }
    this.http_method, this.uri, this.rid = uriInfo[0], uriInfo[1], uriInfo[2]

    // class call and args related
    callRaw := strings.Replace(parts[6], "{", "", -1) 
    callRaw = strings.Replace(callRaw, "}", "", -1) 
    callRaw = strings.Replace(callRaw, "\"", "", -1) 
    callRaw = strings.Replace(callRaw, "[", "", -1) 
    callRaw = strings.Replace(callRaw, "]", "", -1) 
    callRaw = strings.Replace(callRaw, ",", ":", -1) 
    callInfo := strings.Split(callRaw, ":")
    time, err := strconv.Atoi(callInfo[1])
    if err != nil {
        log.Fatal(line, err)
    }
    this.time = int16(time)
    this.class = callInfo[3]
    if len(callInfo) > 10 {
        this.method = callInfo[10]
    }
}

// better printable Request
func (this *amfRequest) String() string {
    return fmt.Sprintf("amfRequest{http:%s uri:%s rid:%s class:%s method:%s time:%d args:%s}",
        this.http_method, this.uri, this.rid, this.class, this.method, this.time, this.args)
}

// TODO use dynamic base polymorphism
func (this *AmfDlog) ScanLines() {
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

        if !this.IsLineValid(line) {
            continue
        }

        // extract info from this line
        this.OperateLine(line)
    }

    if err := run.Wait(); err != nil {
        log.Fatal(err)
    }

    this.chLines <- lineCount
}

func (this *AmfDlog) IsLineValid(line string) bool {
    // super
    if !this.Dlog.IsLineValid(line) {
        return false
    }

    // must exists
    for _, regex := range lineValidatorRegexes[0] {
        if !strings.Contains(line, regex) {
            return false
        }
    }

    // must not exists
    for _, regex := range lineValidatorRegexes[1] {
        if strings.Contains(line, regex) {
            return false
        }
    }

    return true
}

// operate on a valid dlog line
func (this *AmfDlog) OperateLine(line string) {
    this.Dlog.OperateLine(line) // super

    req := new(amfRequest)
    req.ParseLine(line)

    this.lock.Lock()
    defer this.lock.Unlock()
    fmt.Printf("%6d%25s %35s   %s\n", req.time, req.class, req.method, req.uri)
}

