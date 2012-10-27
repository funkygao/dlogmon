// Find co-occurence of 2 terms in a file
// Just a simplel demo
package dlog

import (
    "kx/mr"
    "kx/stats"
    "strings"
    T "kx/trace"
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
    line = strings.TrimSpace(strings.Replace(line, "  ", " ", 0))
    terms := strings.Split(line, " ")
    for i, term := range terms {
        for j:=i+1; j<len(terms); j++ {
            var coOccurence tuple
            coOccurence[0] = strings.TrimSpace(term)
            coOccurence[1] = strings.TrimSpace(terms[j])
            kv[coOccurence] = 1
        }
    }

    out <- kv
}

// Reduce
func (this *FileWorker) Reduce(in mr.KeyValues) (out mr.KeyValue) {
    defer T.Un(T.Trace(""))

    this.Println(this.name, "start to reduce...")

    out = mr.NewKeyValue()
    for k, v := range in {
        var occurence = stats.StatsSum(convertAnySliceToFloat(v))
        if occurence > 1 {
            out[k.(tuple)] = occurence
        }
    }

    return
}
