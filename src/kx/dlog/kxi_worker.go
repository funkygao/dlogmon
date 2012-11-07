package dlog

import (
    "fmt"
    "kx/mr"
    "kx/stats"
    T "kx/trace"
    "os"
)

const (
    GROUP_URL_SERV = "how kxi are called by web"
    GROUP_URL_RID = "which request hit the most kxi call"
    GROUP_URL_SQL = "what kind of sql are most frequent"
    GROUP_KXI = "running stats of kxi servants"
)

const (
    TIME_ALL = "T"
    TIME_AVG = "Tm"
    TIME_MAX = "Tmax"
    TIME_MIN = "Tmin"
    TIME_STD = "Tstd"

    CALL_ALL = "C"
    CALL_AVG = "Cm"
    CALL_MAX = "Cmax"
    CALL_MIN = "Cmin"
    CALL_STD = "Cstd"
)

var KEY_LENS = map[string][]int{
    GROUP_URL_SERV: []int{50, 24},
    GROUP_KXI: []int{50},
    GROUP_URL_SQL:[]int{40, 57},
    GROUP_URL_RID: []int{60, 20}}

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

    kv := mr.NewKeyValue()
    kg1 := mr.NewGroupKey(GROUP_URL_SERV, rec.Url, rec.Service)
    kg2 := mr.NewGroupKey(GROUP_URL_RID, rec.Url, rec.Rid)
    kg3 := mr.NewGroupKey(GROUP_KXI, rec.Service)
    kv[kg1] = rec.Time
    kv[kg2] = rec.Time
    kv[kg3] = rec.Time

    if rec.Sql != "" {
        kg4 := mr.NewGroupKey(GROUP_URL_SQL, rec.Url, rec.Sql)
        kv[kg4] = rec.Time
    }

    kv.Emit(out)
}

// The key is already sorted
func (this *KxiWorker) Reduce(key interface{}, values []interface{}) (kv mr.KeyValue) {
    kv = mr.NewKeyValue()
    vals := mr.ConvertAnySliceToFloat(values)
    switch key.(mr.GroupKey).Group() {
    case GROUP_URL_SERV:
        kv[TIME_ALL] = stats.StatsSum(vals)
        kv[TIME_MIN] = stats.StatsMin(vals)
        kv[TIME_MAX] = stats.StatsMax(vals)
        kv[TIME_AVG] = stats.StatsMean(vals)
        kv[TIME_STD] = stats.StatsSampleStandardDeviation(vals)
        kv[CALL_ALL] = float64(stats.StatsCount(vals))
    case GROUP_KXI:
        kv[TIME_ALL] = stats.StatsSum(vals)
        kv[TIME_MIN] = stats.StatsMin(vals)
        kv[TIME_MAX] = stats.StatsMax(vals)
        kv[TIME_AVG] = stats.StatsMean(vals)
        kv[TIME_STD] = stats.StatsSampleStandardDeviation(vals)
        kv[CALL_ALL] = float64(stats.StatsCount(vals))
    case GROUP_URL_RID:
        kv[CALL_ALL] = float64(stats.StatsCount(vals))
        kv[TIME_ALL] = stats.StatsSum(vals)
    case GROUP_URL_SQL:
        kv[CALL_ALL] = float64(stats.StatsCount(vals))
        kv[TIME_MAX] = stats.StatsMax(vals)
        kv[TIME_AVG] = stats.StatsMean(vals)
    }

    return
}

func (this KxiWorker) KeyLengths(group string) []int {
    if r, e := this.manager.ConfInts(W_KXI, group + " keylen"); e == nil {
        return r
    }

    return KEY_LENS[group]
}
