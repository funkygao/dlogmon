package dlog

import (
    "kx/mr"
    "kx/stats"
    T "kx/trace"
)

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
    kv[line] = 1.0
    out <- kv
}

// Reduce
func (this *FileWorker) Reduce(in mr.KeyValues) (out mr.KeyValue) {
    defer T.Un(T.Trace(""))

    this.Println(this.name, "start to reduce...")

    out = mr.NewKeyValue()
    for k, v := range in {
        out[k.(string)] = stats.StatsSum(convertAnySliceToFloat(v))
    }

    return
}
