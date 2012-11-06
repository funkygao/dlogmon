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

    return this
}

// Does a log line contain 'AMF_SLOW'?
func (this *AmfWorker) IsLineValid(line string) bool {
    if !isSamplerHostLine(line) {
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
    req := new(amfRequest)
    req.parseLine(line)

    kv := mr.NewKeyValue()
    key := mr.NewKey(req.class, req.method, req.uri)
    kv[key] = 1

    kv.Emit(out)
}

func (this *AmfWorker) Reduce(key interface{}, values []interface{}) (kv mr.KeyValue) {
    // here we don't care about the key
    // we only care about values
    aggregate := stats.StatsSum(mr.ConvertAnySliceToFloat(values))
    if aggregate > 0 {
        kv = mr.NewKeyValue()
        kv[NIL_KEY] = aggregate
    }

    return
}

func (this AmfWorker) Printr(key interface{}, value mr.KeyValue) string {
    k := key.(mr.Key)
    keys := k.Keys()
    method := keys[0] + "." + keys[1]
    fmt.Printf("%65s  %-35s %v\n", method, keys[2], value[NIL_KEY])
    return ""
}
