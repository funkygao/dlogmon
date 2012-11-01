package dlog

import (
    "encoding/json"
    "fmt"
    "kx/mr"
    "kx/stats"
    T "kx/trace"
    "os"
)

const (
    KXI_M_KEYLEN = 4
    KXI_R_KEYLEN = 3
)

const (
    TIME_ALL = "T"
    TIME_AVG = "Tm"
    TIME_MAX = "Ta"
    TIME_MIN = "Ti"
    TIME_STD = "Td"

    CALL_ALL = "C"
    CALL_AVG = "Cm"
    CALL_MAX = "Ca"
    CALL_MIN = "Ci"
    CALL_STD = "Cd"
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
    kv[[KXI_M_KEYLEN]string{mr.KEY_GROUP, "url call kxi service", url, service}] = time
    kv[[KXI_M_KEYLEN]string{mr.KEY_GROUP, "url within a request", url, rid}] = time
    kv[[KXI_M_KEYLEN]string{mr.KEY_GROUP, "url query db sql", url, sql}] = time
    kv[[KXI_M_KEYLEN]string{mr.KEY_GROUP, "kxi servants", service, ""}] = time
    out <- kv
}

// The key is already sorted
func (this *KxiWorker) Reduce(key interface{}, values []interface{}) (kv mr.KeyValue) {
    parts := key.([KXI_M_KEYLEN]string)
    grp, kind, k1, k2 := parts[0], parts[1], parts[2], parts[3]
    if this.manager.option.debug {
        fmt.Fprintf(os.Stderr, "DEBUG=> %s %s %s %s %v\n", grp, kind, k1, k2, values)
    }

    kv = mr.NewKeyValue()
    switch kind {
    case "url call kxi service":
        url, service := k1, k2
        kv[[KXI_R_KEYLEN]string{url, service, TIME_ALL}] = stats.StatsSum(mr.ConvertAnySliceToFloat(values))
        kv[[KXI_R_KEYLEN]string{url, service, CALL_ALL}] = float64(len(values))
    case "url within a request":
        url, rid := k1, k2
        kv[[KXI_R_KEYLEN]string{url, rid, TIME_ALL}] = stats.StatsSum(mr.ConvertAnySliceToFloat(values))
        kv[[KXI_R_KEYLEN]string{url, rid, CALL_ALL}] = float64(len(values))
    case "kxi servants":
    case "url query db sql":
    }

    return
}

// kv are in the same group
func (this KxiWorker) Printh(kv mr.KeyValue, top int) {
    s := mr.NewSort(kv)
    s.Sort(mr.SORT_BY_KEY, mr.SORT_ORDER_DESC)
    sortedKeys := s.Keys()
    if top > 0 && top < len(sortedKeys) {
        sortedKeys = sortedKeys[:top]
    }

    metrics := mr.NewKeyValue()
    for _, sk := range sortedKeys {
        value := kv[sk].(mr.KeyValue)
        for k, _ := range value {
            metric := k.([KXI_R_KEYLEN]string)[2]
            metrics[metric] = true
        }
    }
    fmt.Printf("%70s %20s", "", "")
    for _, x := range metrics.Keys() {
        fmt.Printf("%8s", x.(string))
    }
    println()

    var lastK12 string
    for _, sk := range sortedKeys {
    //    fmt.Println(sk)
     //   continue
        value := kv[sk].(mr.KeyValue)
        tt := make(map[string]float64)
        for k, v := range value {
            key := k.([KXI_R_KEYLEN]string)
            tt[key[2]] = v.(float64)
            //col := key[2]
            if lastK12 == "" {
                lastK12 = key[0]+key[1]
            }
            if lastK12 != key[0]+key[1] {
                fmt.Printf("%70s %20s", key[0], key[1])
                for _, bb := range metrics.Keys() {
                    fmt.Printf("%8.0f", tt[bb.(string)])
                }
                println()
                //fmt.Printf("%70s %20s %5s %10.0f\n", key[0], key[1], key[2], v)
                lastK12 = key[0]+key[1]
            }
        }
    }
}

// key is key of mapper output
// value is output of reducer
// Should the 2 key have relationship? TODO
func (this KxiWorker) Printr(key interface{}, value mr.KeyValue) string {
    for k, v := range value {
        key := k.([KXI_R_KEYLEN]string)
        fmt.Printf("%70s %20s %5s %10.0f\n", key[0], key[1], key[2], v)
    }
    return ""
}
