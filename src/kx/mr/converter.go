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

func InterfaceArrayToStringSlice(in []interface{}) []string {
    r := make([]string, len(in))
    for i, v := range in {
        r[i] = v.(string)
    }
    return r
}

func ArrayToSlice(in interface{}) []string {
    if v, ok := in.([2]string); ok {
        return arrayToSlice2(v)
    } else if v, ok := in.([3]string); ok {
        return arrayToSlice3(v)
    } else if v, ok := in.([4]string); ok {
        return arrayToSlice4(v)
    } else if v, ok := in.([5]string); ok {
        return arrayToSlice5(v)
    } else if v, ok := in.([6]string); ok {
        return arrayToSlice6(v)
    } else if v, ok := in.([7]string); ok {
        return arrayToSlice7(v)
    } else if v, ok := in.([8]string); ok {
        return arrayToSlice8(v)
    }

    return nil
}

func arrayToSlice2(in [2]string) []string {
    out := make([]string, len(in))
    for i, s := range in {
        out[i] = s
    }
    return out
}

func arrayToSlice3(in [3]string) []string {
    out := make([]string, len(in))
    for i, s := range in {
        out[i] = s
    }
    return out
}

func arrayToSlice4(in [4]string) []string {
    out := make([]string, len(in))
    for i, s := range in {
        out[i] = s
    }
    return out
}

func arrayToSlice5(in [5]string) []string {
    out := make([]string, len(in))
    for i, s := range in {
        out[i] = s
    }
    return out
}

func arrayToSlice6(in [6]string) []string {
    out := make([]string, len(in))
    for i, s := range in {
        out[i] = s
    }
    return out
}

func arrayToSlice7(in [7]string) []string {
    out := make([]string, len(in))
    for i, s := range in {
        out[i] = s
    }
    return out
}
func arrayToSlice8(in [8]string) []string {
    out := make([]string, len(in))
    for i, s := range in {
        out[i] = s
    }
    return out
}
