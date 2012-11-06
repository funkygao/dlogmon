package dlog

import (
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
    if !isSamplerHostLine(line) {
        return false
    }

    // mapper stream will decide valididation of the line
    return true
}

func (this *KxiWorker) Map(line string, out chan<- mr.KeyValue) {
    var streamResult StreamResult
    if streamResult = this.Worker.streamedResult(line); streamResult.Empty() {
        return
    }

    type record struct {
        Url     string  `json:"u"`
        Rid     string  `json:"i"`
        Service string  `json:"s"`
        Time    float64 `json:"t"`
        Sql     string  `json:"q"`
    }
    rec := new(record)
    if err := streamResult.Decode(rec); err != nil {
        panic(err)
    }

    if this.manager.option.debug {
        fmt.Fprintf(os.Stderr, "DEBUG<= %s %s %s %f %s\n",
            rec.Url, rec.Rid, rec.Service, rec.Time, rec.Sql)
    }

    // TODO refactor from here

    kv := mr.NewKeyValue()
    kg1 := mr.NewGroupKey("url call kxi service", rec.Url, rec.Service)
    kg2 := mr.NewGroupKey("url within a request", rec.Url, rec.Rid)
    kg3 := mr.NewGroupKey("url query db sql", rec.Url, rec.Sql)
    kg4 := mr.NewGroupKey("kxi servants", rec.Service)
    kv[kg1] = rec.Time
    kv[kg2] = rec.Time
    kv[kg3] = rec.Time
    kv[kg4] = rec.Time

    kv.Emit(out)
}

// The key is already sorted
func (this *KxiWorker) Reduce(key interface{}, values []interface{}) (kv mr.KeyValue) {
    kv = mr.NewKeyValue()
    switch key.(mr.GroupKey).Group() {
    case "url call kxi service":
        kv[TIME_ALL] = stats.StatsSum(mr.ConvertAnySliceToFloat(values))
        kv[CALL_ALL] = float64(len(values))
    case "url within a request":
    case "kxi servants":
    case "url query db sql":
    }

    return
}

func (this KxiWorker) SortCol(group string) string {
    switch group {
    case "url call kxi service":
        return TIME_ALL
    }
    return ""
}

func (this KxiWorker) KeyLengths(group string) []int {
    switch group {
    case "url call kxi service":
        return []int{60, 20}
    }
    return nil
}

