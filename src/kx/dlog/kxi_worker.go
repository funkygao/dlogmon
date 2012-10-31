package dlog

import (
    "fmt"
    "kx/mr"
    T "kx/trace"
    "os"
)

// Constructor of KxiWorker
func NewKxiWorker(manager *Manager, name, filename string, seq uint16) IWorker {
    defer T.Un(T.Trace(""))

    this := new(KxiWorker)
    this.self = this
    this.init(manager, name, filename, seq)

    return this
}

func (this *KxiWorker) IsLineValid(line string) bool {
    if !this.Worker.IsLineValid(line) {
        return false
    }

    return true
}

func (this *KxiWorker) Map(line string, out chan<- mr.KeyValue) {
    if x := this.Worker.ExtractLineInfo(line); x != nil {
        if this.manager.option.debug {
            fmt.Fprintf(os.Stderr, "%s", x)
        }
    }
    out <- mr.NewKeyValue()
}

func (this *KxiWorker) Reduce(key interface{}, values []interface{}) (kv mr.KeyValue) {
    return
}

func (this KxiWorker) Printr(key interface{}, value interface{}) string {
    return ""
}
