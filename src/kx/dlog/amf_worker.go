package dlog

import (
    "fmt"
    "kx/mr"
    "kx/stats"
    T "kx/trace"
    "strings"
)

// Constructor of AmfWorker
func NewAmfWorker(manager *Manager, name, filename string, seq uint16) IWorker {
    defer T.Un(T.Trace(""))

    this := new(AmfWorker)
    this.self = this // don't forget this
    this.init(manager, name, filename, seq)

    // set the combiner
    //this.combiner = stats.StatsSum

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
func (this *AmfWorker) Map(line string, out chan<- mr.KeyValue) {
    if x := this.Worker.ExtractLineInfo(line); x != nil {
        if this.manager.option.debug {
            this.Println(line)
        }
    }

    req := new(amfRequest)
    req.parseLine(line)

    kv := mr.NewKeyValue()
    kv[req.class+"."+req.method+"@"+req.uri] = 1

    // emit an intermediate data
    out <- kv
}

// Reduce
func (this *AmfWorker) Reduce(key interface{}, values []interface{}) (kv mr.KeyValue) {
    aggregate := stats.StatsSum(mr.ConvertAnySliceToFloat(values))
    if aggregate > 0 {
        kv = mr.NewKeyValue()
        kv[key] = aggregate
    }

    return
}

func (this AmfWorker) Printr(key interface{}, value interface{}) string {
    v := value.(mr.KeyValue)
    p := strings.Split(key.(string), "@")
    fmt.Printf("%58s    %-35s %4.0f\n", p[0], p[1], v[key])
    return ""
}
