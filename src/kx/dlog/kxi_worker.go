package dlog

import (
    "encoding/json"
    "fmt"
    "kx/mr"
    T "kx/trace"
    "os"
)

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

    // mapper stream will decide valididation of the line
    return true
}

func (this *KxiWorker) Map(line string, out chan<- mr.KeyValue) {
    const (
        KEY_URL = "u"
        KEY_RID = "i"
        KEY_SERVICE = "s"
        KEY_TIME = "t"
        KEY_SQL = "q"
    )
    var streamResult interface{}
    if streamResult = this.Worker.ExtractLineInfo(line); streamResult == nil {
        // this line is invalid
        return
    }

    streamKv := make(map[string]interface{})
    if err := json.Unmarshal([]byte(streamResult.(string)), &streamKv); err != nil {
        panic(err)
    }

    if this.manager.option.debug {
        fmt.Fprintf(os.Stderr, "stream=> %#v\n", streamKv)
    }

    kv := mr.NewKeyValue()
    out <- kv
}

func (this *KxiWorker) Reduce(key interface{}, values []interface{}) (kv mr.KeyValue) {
    return
}

func (this KxiWorker) Printr(key interface{}, value interface{}) string {
    return ""
}
