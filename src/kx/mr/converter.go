package mr

// []interface{}.([]float64) is not supported in golang, see FAQ
// so we need do it ourself
func ConvertAnySliceToFloat(v []interface{}) []float64 {
	r := make([]float64, 0)
	for i := range v {
		if d, ok := v[i].(int); ok {
			r = append(r, float64(d))
		} else if d, ok := v[i].(int32); ok {
			r = append(r, float64(d))
		} else if d, ok := v[i].(int64); ok {
			r = append(r, float64(d))
        } else if d, ok := v[i].(float64); ok {
			r = append(r, d)
		} else {
			panic("Unkown type")
		}
	}

	return r
}

