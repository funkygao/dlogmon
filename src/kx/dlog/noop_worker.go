package dlog

import (
    "kx/mr"
    T "kx/trace"
)

// Constructor of NoopWorker
func NewNoopWorker(manager *Manager, name, filename string, seq uint16) IWorker {
    defer T.Un(T.Trace(""))

    this := new(NoopWorker)
    this.name = name
    this.filename = filename
    this.seq = seq
    this.manager = manager
    this.self = this
    this.Logger = this.manager.Logger

    return this
}

func (this *NoopWorker) IsLineValid(line string) bool {
    return true
}

// Extract meta info related to amf from a valid line
func (this *NoopWorker) Map(line string, out chan<- interface{}) {
    d := mr.NewMapData()

    out <- d
}

// Reduce
func (this *NoopWorker) Reduce(in mr.ReduceData) mr.ReduceResult {
    defer T.Un(T.Trace(""))

    this.Println(this.name, "start to reduce...")

    return mr.NewReduceResult(len(in))
}
