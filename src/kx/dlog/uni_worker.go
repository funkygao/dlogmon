package dlog

import (
    "fmt"
    "kx/mr"
    "regexp"
    "kx/stats"
    T "kx/trace"
    "strings"
)

var re *regexp.Regexp

func init() {
    re = regexp.MustCompile(`\d*`)
}

func NewUniWorker(manager *Manager, name, filename string, seq uint16) IWorker {
    defer T.Un(T.Trace(""))

    this := new(UniWorker)
    this.self = this // don't forget this
    this.init(manager, name, filename, seq)

    return this
}

func (this *UniWorker) IsLineValid(line string) bool {
    return true
}

func (this *UniWorker) Map(line string, out chan<- mr.KeyValue) {
    line = line[:len(line)-1]
    line = strings.Replace(line, "+", " ", 2)
    kv := mr.NewKeyValue()
    rs := re.ReplaceAll([]byte(line), []byte{})
    line = string(rs)
    line = trimAllRune(line, []rune{'>', '~', ';', '.', '-', '*'})
    words := strings.Split(line, " ")
    for _, w := range words {
        line = strings.Trim(line, "  ")
        k := mr.NewKey(w)
        kv[k] = 1
    }

    kv.Emit(out)
}

func (this *UniWorker) Reduce(key interface{}, values []interface{}) (kv mr.KeyValue) {
    // here we don't care about the key
    // we only care about values
    aggregate := stats.StatsSum(mr.ConvertAnySliceToFloat(values))
    if aggregate > 100 {
        kv = mr.NewKeyValue()
        kv[NIL_KEY] = aggregate
    }

    return
}

func (this UniWorker) Printr(key interface{}, value mr.KeyValue) string {
    k := key.(mr.Key)
    fmt.Printf("%85s %10.0f\n", k.Keys()[0], value[NIL_KEY])
    return ""
}
