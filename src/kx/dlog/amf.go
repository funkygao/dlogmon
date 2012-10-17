// line reader
package dlog

import (
    "fmt"
    "io"
    "log"
    "os"
    "strings"
    "strconv"
    "sync"
)

var lineValidatorRegexes = [...][]string{
    {"AMF_SLOW", "PHP.CDlog"}, // must exists
    {"Q=DLog.log"}} // must not exists

// AMF_SLOW dlog analyzer
type AmfDlog struct {
    Dlog
}

// a single line meta info
type amfRequest struct {
    Request
    class, method, args string
    time int16
}

// parse a line into meta info
// better printable Request
func (this *amfRequest) String() string {
    return fmt.Sprintf("amfRequest{http:%s uri:%s rid:%s class:%s method:%s time:%d args:%s}",
        this.http_method, this.uri, this.rid, this.class, this.method, this.time, this.args)
}

// Constructor of AmfDlog
func NewAmfDlog(filename string, ch chan int, lock *sync.Mutex, options *Options) IDlogExecutor {
    this := new(AmfDlog)
    this.filename = filename
    this.chLines = ch
    this.lock = lock
    this.options = options

    var logWriter io.Writer
    if options.logfile == "" {
        logWriter = os.Stderr
    } else {
        f, e := os.OpenFile(options.logfile, os.O_APPEND | os.O_CREATE, 0666)
        if e != nil {
            panic(e)
        }
        logWriter = f
    }
    // notice how to access embedded types
    this.Logger = log.New(logWriter, "",  log.Ldate | log.Lshortfile | log.Ltime | log.Lmicroseconds)
    return this
    /*
    return &AmfDlog{
        Dlog{
            filename,
            ch,
            lock,
            options,
            log.New(os.Stderr, "",  log.Ldate | log.Llongfile | log.Ltime | log.Lmicroseconds),
            nil, nil}}
            */
}

func (this *amfRequest) parseLine(line string) {
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
func (this *AmfDlog) OperateLine(line string) Any {
    if x := this.Dlog.OperateLine(line); x != nil {
        if this.options.debug {
            this.Println(line)
        }
        return x
    }

    req := new(amfRequest)
    req.parseLine(line)

    this.lock.Lock()
    defer this.lock.Unlock()

    line = fmt.Sprintf("%d %s %s", req.time, req.class + "." + req.method, req.uri)
    if this.options.debug {
        this.Println(line)
    }

    return line
}

