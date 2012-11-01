package dlog

import (
    "fmt"
    "kx/mr"
    "kx/stats"
    T "kx/trace"
    "os"
    "strings"
)

const AMF_KEY_LEN = 2

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

func (this *AmfWorker) Map(line string, out chan<- mr.KeyValue) {
    if x := this.Worker.ExtractLineInfo(line); x != nil {
        if this.manager.option.debug {
            fmt.Fprintf(os.Stderr, "%s", x)
        }
    }

    req := new(amfRequest)
    req.parseLine(line)

    kv := mr.NewKeyValue()
    kv[[AMF_KEY_LEN]string{req.class + "." + req.method, req.uri}] = 1

    // emit an intermediate data
    out <- kv
}

func (this *AmfWorker) Reduce(key interface{}, values []interface{}) (kv mr.KeyValue) {
    // here we don't care about the key
    // we only care about values
    aggregate := stats.StatsSum(mr.ConvertAnySliceToFloat(values))
    if aggregate > 0 {
        kv = mr.NewKeyValue()
        kv[key] = aggregate
    }

    return
}

func (this AmfWorker) Printr(key interface{}, value interface{}) string {
    v := value.(mr.KeyValue)
    k := key.([AMF_KEY_LEN]string)
    if count := v[k].(float64); count >= this.printThreshold() {
        fmt.Printf("%65s  %-35s %5.0f\n", k[0], k[1], v[k])
    }

    return fmt.Sprintf("insert into %s(method, uri, c) values('%s', '%s', %d)",
        TABLE_AMF, k[0], k[1], v[k])
}

func (this AmfWorker) printThreshold() float64 {
    const default_threshold = 5
    t, e := this.manager.Conf().Float("amf", "export.threshold")
    if e != nil {
        return default_threshold
    }
    return t
}
