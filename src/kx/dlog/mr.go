package dlog

import (
    "fmt"
)

func newMapOut() MapOut {
    return make(MapOut)
}

func newReduceIn() ReduceIn {
    return make(ReduceIn)
}

func getKey(t KeyType, key string) string {
    return fmt.Sprintf("%d%s%s", t, KEYTYPE_SEP, key)
}

func (this MapOut) Set(t KeyType, key string, val int) {
    key = getKey(t, key)
    this[key] = val
}

func (this MapOut) Get(key string) (val int, ok bool) {
    val, ok= this[key]
    return
}

func (this ReduceIn) Append(key string, val int) {
    _, ok := this[key]
    if !ok {
        this[key] = make([]int, 1)
        this[key][0] = val
    } else {
        this[key] = append(this[key], val)
    }
}
