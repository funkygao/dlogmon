package dlog

func intSum(vals []int) (sum int) {
    for _, val := range vals {
        sum += val
    }
    return
}

func intAvg(vals []int) (avg int) {
    avg = intSum(vals) / len(vals)
    return
}
