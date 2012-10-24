package dlog

import (
    "fmt"
    "kx/mr"
    "kx/stats"
    T "kx/trace"
    "strings"
)

// Printable amfRequest 
func (this *amfRequest) String() string {
    return fmt.Sprintf("amfRequest{http:%s uri:%s rid:%s class:%s method:%s time:%d args:%s}",
        this.http_method, this.uri, this.rid, this.class, this.method, this.time, this.args)
}

// Constructor of AmfWorker
func NewAmfWorker(manager *Manager, name, filename string) IWorker {
    defer T.Un(T.Trace(""))

    this := new(AmfWorker)
    this.name = name
    this.filename = filename
    this.manager = manager
    this.executor = this

    // notice how to access embedded types
    this.Logger = this.manager.Logger

    // set the combiner
//    this.combiner = stats.StatsSum

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
func (this *AmfWorker) Map(line string, out chan<- interface{}) {
    if x := this.Worker.ExtractLineInfo(line); x != nil {
        if this.manager.option.debug {
            this.Println(line)
        }
    }

    req := new(amfRequest)
    req.parseLine(line)

    d := mr.NewMapData()
    // keyType must starts with 0
    d.Set(0, req.class + "." + req.method, 1)
    d.Set(1, req.uri, 1)
    d.Set(2, req.rid, 1)

    out <- d
}

// Reduce
func (this *AmfWorker) Reduce(in mr.ReduceData) (r mr.ReduceResult) {
    defer T.Un(T.Trace(""))

    this.Println(this.name, "reduce")

    r = mr.NewReduceResult(len(in))
    for keyType, d := range in {
        for k, v := range d {
            // sum up
            r[keyType][k] = stats.StatsSum(v)
        }
    }

    return
}
