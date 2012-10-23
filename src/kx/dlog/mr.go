package dlog

import (
    "fmt"
)

func newMapData() MapData {
    return make(MapData)
}

func newTransformData() TransformData {
    return make(TransformData)
}

func getKey(t KeyType, key string) string {
    return fmt.Sprintf("%d%s%s", t, KEYTYPE_SEP, key)
}

func (this MapData) Set(t KeyType, key string, val float64) {
    key = getKey(t, key)
    this[key] = val
}

func (this MapData) Get(key string) (val float64, ok bool) {
    val, ok= this[key]
    return
}

func (this TransformData) Append(key string, val float64) {
    _, ok := this[key]
    if !ok {
        this[key] = make([]float64, 1)
        this[key][0] = val
    } else {
        this[key] = append(this[key], val)
    }
}

func (this TransformData) AppendSlice(key string, val []float64) {
    _, ok := this[key]
    if !ok {
        this[key] = make([]float64, 1)
        this[key] = val
    } else {
        this[key] = append(this[key], val...)
    }
}
