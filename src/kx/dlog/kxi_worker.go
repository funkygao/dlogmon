package dlog

import (
    "encoding/json"
    "fmt"
    "kx/mr"
    "kx/stats"
    T "kx/trace"
    "os"
)

const KXI_KEY_LEN = 3

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

    var (
        url, rid, service, sql string
        time float64
    )
    url, rid, service, time, sql = streamKv[KEY_URL].(string), streamKv[KEY_RID].(string),
        streamKv[KEY_SERVICE].(string), streamKv[KEY_TIME].(float64),
        streamKv[KEY_SQL].(string)
    if this.manager.option.debug {
        fmt.Fprintf(os.Stderr, "DEBUG<= %s %s %s %f %s\n", url, rid, service, time, sql)
    }

    kv := mr.NewKeyValue()
    kv[[KXI_KEY_LEN]string{"url_service", url, service}] = time
    kv[[KXI_KEY_LEN]string{"url", url, ""}] = time
    kv[[KXI_KEY_LEN]string{"service", service, ""}] = time
    kv[[KXI_KEY_LEN]string{"url_rid", url, rid}] = time
    kv[[KXI_KEY_LEN]string{"url_sql", url, sql}] = time
    out <- kv
}

// The key is already sorted
func (this *KxiWorker) Reduce(key interface{}, values []interface{}) (kv mr.KeyValue) {
    parts := key.([KXI_KEY_LEN]string)
    kind, k1, k2 := parts[0], parts[1], parts[2]
    if this.manager.option.debug {
        fmt.Fprintf(os.Stderr, "DEBUG=> %s %s %s %v\n", kind, k1, k2, values)
    }

    kv = mr.NewKeyValue()
    if kind == "service" {
        kv[k1] = stats.StatsSum(mr.ConvertAnySliceToFloat(values))
    }
    return
}

func (this KxiWorker) Printr(key interface{}, value mr.KeyValue) string {
    for k, v := range value {
        fmt.Printf("%50s %.0f\n", k, v)
    }
    return ""
}
