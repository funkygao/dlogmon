package dlog

// TODO
func convertAnySliceToFloat(v []interface{}) []float64 {
    r := make([]float64, 0)
    for i, _ := range v {
        d, ok := v[i].(float64)
        if !ok {
            panic("invalid type")
        }
        r = append(r, d)
    }

    return r
}
