package mr

import (
    "fmt"
    "reflect"
)

func GenericSlice(slice interface{}) []interface{} {
    v := reflect.ValueOf(slice)
    if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
        panic(fmt.Errorf("Expected slice or array, got %T", slice))
    }
    l := v.Len()
    r := make([]interface{}, l)
    for i:=0; i<l; i++ {
        r[i] = v.Index(i).Interface()
    }
    return r
}

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

func ConvertAnySliceToString(v []interface{}) []string {
    r := make([]string, 0)
    for i := range v {
        if d, ok := v[i].(string); ok {
            r = append(r, d)
        } else {
            panic("Unkown type")
        }
    }

    return r
}
