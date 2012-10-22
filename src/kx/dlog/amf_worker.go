package dlog

import (
    "fmt"
    t "kx/trace"
    "strconv"
    "strings"
)

// Printable amfRequest 
func (this *amfRequest) String() string {
    return fmt.Sprintf("amfRequest{http:%s uri:%s rid:%s class:%s method:%s time:%d args:%s}",
        this.http_method, this.uri, this.rid, this.class, this.method, this.time, this.args)
}

// Constructor of AmfWorker
func NewAmfWorker(manager *Manager, filename string) IWorker {
    defer t.Un(t.Trace(""))

    this := new(AmfWorker)
    this.filename = filename
    this.manager = manager

    // notice how to access embedded types
    this.Logger = this.manager.Logger

    return this
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
func (this *AmfWorker) IsLineValid(line string) bool {
    // super
    if !this.Worker.IsLineValid(line) {
        return false
    }

    // must exists
    for _, regex := range amfLineValidatorRegexes[0] {
        if !strings.Contains(line, regex) {
            return false
        }
    }

    // must not exists
    for _, regex := range amfLineValidatorRegexes[1] {
        if strings.Contains(line, regex) {
            return false
        }
    }

    return true
}

// Extract meta info related to amf from a valid line
func (this *AmfWorker) ExtractLineInfo(line string) Any {
    if x := this.Worker.ExtractLineInfo(line); x != nil {
        if this.manager.option.debug {
            this.Println(line)
        }
        return x
    }

    req := new(amfRequest)
    req.parseLine(line)

    this.manager.Lock()
    defer this.manager.Unlock()

    line = fmt.Sprintf("%d %s %s", req.time, req.class+"."+req.method, req.uri)
    if this.manager.option.debug {
        this.Println(line)
    }

    return line
}
