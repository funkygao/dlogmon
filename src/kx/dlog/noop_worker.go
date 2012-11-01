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

func (this *NoopWorker) Map(line string, out chan<- mr.KeyValue) {
    out <- mr.NewKeyValue()
}

func (this *NoopWorker) Reduce(key interface{}, values []interface{}) (kv mr.KeyValue) {
    return
}

func (this NoopWorker) Printr(key interface{}, value mr.KeyValue) string {
    return ""
}
