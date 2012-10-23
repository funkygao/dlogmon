package dlog

import (
    "fmt"
    "kx/stats"
    t "kx/trace"
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
    this.executor = this

    // notice how to access embedded types
    this.Logger = this.manager.Logger

    // set the combiner
    this.combiner = stats.StatsSum

    return this
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
func (this *AmfWorker) Map(line string, out chan<- Any) {
    if x := this.Worker.ExtractLineInfo(line); x != nil {
        if this.manager.option.debug {
            this.Println(line)
        }
    }

    req := new(amfRequest)
    req.parseLine(line)

    d := newMapData()
    d.Set(1, req.class + "." + req.method, 1)
    d.Set(2, req.uri, 1)
    d.Set(3, req.rid, 1)

    out <- d
}
