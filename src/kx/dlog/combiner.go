package dlog

func intSum(vals []Any) (sum int) {
    for _, val := range vals {
        sum += val.(int)
    }
    return
}

func intAvg(vals []Any) (avg float64) {
    sum, size := float64(intSum(vals)), float64(len(vals))
    return sum / size
}
