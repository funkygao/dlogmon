package dlog

import (
    "kx/mr"
    T "kx/trace"
)

// Constructor of NoopWorker
func NewNoopWorker(manager *Manager, name, filename string, seq uint16) IWorker {
    defer T.Un(T.Trace(""))

    this := new(NoopWorker)
    this.self = this
    this.init(manager, name, filename, seq)

    return this
}

func (this *NoopWorker) IsLineValid(line string) bool {
    return false
}

// Extract meta info related to amf from a valid line
func (this *NoopWorker) Map(line string, out chan<- mr.KeyValue) {
    out <- mr.NewKeyValue()
}

// Reduce
func (this *NoopWorker) Reduce(in mr.KeyValues) mr.KeyValue {
    defer T.Un(T.Trace(""))

    this.Println(this.name, "start to reduce...")

    return mr.NewKeyValue()
}
