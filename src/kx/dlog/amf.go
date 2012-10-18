package dlog

import (
    "fmt"
    t "kx/trace"
    "strconv"
    "strings"
)

var lineValidatorRegexes = [...][]string{
    {"AMF_SLOW", "PHP.CDlog"}, // must exists
    {"Q=DLog.log"}}            // must not exists

// AMF_SLOW tag analyzer
type AmfDlog struct {
    Dlog
}

// a single line meta info
type amfRequest struct {
    Request
    class, method, args string
    time                int16
}

// Printable amfRequest 
func (this *amfRequest) String() string {
    return fmt.Sprintf("amfRequest{http:%s uri:%s rid:%s class:%s method:%s time:%d args:%s}",
        this.http_method, this.uri, this.rid, this.class, this.method, this.time, this.args)
}

// Constructor of AmfDlog
func NewAmfDlog(manager *Manager, filename string) IDlogExecutor {
    defer t.Un(t.Trace("NewAmfDlog"))

    this := new(AmfDlog)
    this.filename = filename
    this.manager = manager

    // notice how to access embedded types
    this.Logger = newLogger(this.manager.option)
    return this

    /*
       return &AmfDlog{
           Dlog{
               filename,
               ch,
               lock,
               option,
               newLogger(this.manager.option),
               nil, 
               nil}}
    */
}

func (this *amfRequest) parseLine(line string) {
    // major parts seperated by space
    parts := strings.Split(line, " ")

    // uri related
    uriInfo := strings.Split(parts[5], "+")
    if len(uriInfo) < 3 {
        panic(uriInfo)
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
        println(line)
        panic(err)
    }
    this.time = int16(time)
    this.class = callInfo[3]
    if len(callInfo) > 10 {
        this.method = callInfo[10]
    }
}

// Does a log line contain 'AMF_SLOW'?
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

// Extract meta info related to amf from a valid line
func (this *AmfDlog) ExtractLineInfo(line string) Any {
    if x := this.Dlog.ExtractLineInfo(line); x != nil {
        if this.manager.option.debug {
            this.Println(line)
        }
        return x
    }

    req := new(amfRequest)
    req.parseLine(line)

    this.manager.lock.Lock()
    defer this.manager.lock.Unlock()

    line = fmt.Sprintf("%d %s %s", req.time, req.class+"."+req.method, req.uri)
    if this.manager.option.debug {
        this.Println(line)
    }

    return line
}
