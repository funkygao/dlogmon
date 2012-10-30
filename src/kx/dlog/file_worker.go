// Find co-occurence of 2 terms in a file
// Just a simplel demo
package dlog

import (
    "fmt"
    "kx/mr"
    "kx/stats"
    T "kx/trace"
    "strings"
)

type tuple [2]string

func (this tuple) String() string {
    return strings.TrimSpace(this[0]) + "," + strings.TrimSpace(this[1])
}

// Constructor of NoopWorker
func NewFileWorker(manager *Manager, name, filename string, seq uint16) IWorker {
    defer T.Un(T.Trace(""))

    this := new(FileWorker)
    this.self = this
    this.init(manager, name, filename, seq)

    return this
}

func (this *FileWorker) IsLineValid(line string) bool {
    return true
}

// Extract meta info related to amf from a valid line
func (this *FileWorker) Map(line string, out chan<- mr.KeyValue) {
    kv := mr.NewKeyValue()
    line = trimAllRune(line, []rune{'=', ':', '+', '.', '-'})
    line = strings.Trim(line, "  ")
    if len(line) == 0 {
        return
    }

    terms := strings.Split(line, " ")
    for i, term := range terms {
        for j := i + 1; j < len(terms); j++ {
            var coOccurence [2]string
            coOccurence[0] = strings.TrimSpace(term)
            coOccurence[1] = strings.TrimSpace(terms[j])
            kv[coOccurence] = 1
        }
    }

    out <- kv
}

// Reduce
func (this *FileWorker) Reduce(key interface{}, values []interface{}) (kv mr.KeyValue) {
    const threhold = 0
    var occurence = stats.StatsSum(mr.ConvertAnySliceToFloat(values))
    if occurence > threhold {
        kv = mr.NewKeyValue()
        kv[key] = occurence
    }

    return
}

func (this FileWorker) Printr(key interface{}, value interface{}) string {
    v := value.(mr.KeyValue)
    k := key.([2]string)
    if v[k].(float64) > 1 {
        fmt.Printf("%25s%25s %4.0f\n", k[0], k[1], v[k])
    }
    return ""
}
