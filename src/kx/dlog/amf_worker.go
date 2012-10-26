package dlog

import (
    "fmt"
    "kx/mr"
    "kx/stats"
    T "kx/trace"
    "reflect"
    "strings"
)

// Printable amfRequest 
func (this *amfRequest) String() string {
    return fmt.Sprintf("amfRequest{http:%s uri:%s rid:%s class:%s method:%s time:%d args:%s}",
        this.http_method, this.uri, this.rid, this.class, this.method, this.time, this.args)
}

// Constructor of AmfWorker
func NewAmfWorker(manager *Manager, name, filename string, seq uint16) IWorker {
    defer T.Un(T.Trace(""))

    this := new(AmfWorker)
    this.self = this // don't forget this
    this.init(manager, name, filename, seq)

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
func (this *AmfWorker) Map(line string, out chan<- interface{}) {
    if x := this.Worker.ExtractLineInfo(line); x != nil {
        if this.manager.option.debug {
            this.Println(line)
        }
    }

    req := new(amfRequest)
    req.parseLine(line)

    d := mr.NewMapData()
    d.Set(0, req.class+"."+req.method, 1.0)
    d.Set(1, req.uri, 10)
    d.Set(2, req.rid, 1.0)

    // emit an intermediate data
    out <- d
}

// TODO
func convert(v []interface{}) []float64 {
    r := make([]float64, 0)
    for i, _ := range v {
        d, ok := v[i].(float64)
        if !ok {
            fmt.Println("shit", d, ok)
        }
        r = append(r, d)
    }
    return r
}

// Reduce
func (this *AmfWorker) Reduce(in mr.ReduceData) (out mr.ReduceResult) {
    defer T.Un(T.Trace(""))

    this.Println(this.name, "start to reduce...")

    out = mr.NewReduceResult(len(in))
    for tagType, d := range in {
        for k, v := range d {
            // sum up []float64 for this key
            x := reflect.TypeOf(v).Elem()
            fmt.Printf("%#v\n", x)
            out[tagType][k] = stats.StatsSum(convert(v))
        }
    }

    return
}
