package mr

import (
    "fmt"
    T "kx/trace"
)

// Factory
func NewKeyValue() KeyValue {
    return make(KeyValue)
}

// Factory
func NewKeyValues() KeyValues {
    return make(KeyValues)
}

// Self printable
func (this KeyValue) Println() {
    for k, v := range this {
        fmt.Printf("%s %v %v\n", T.CallerFuncName(1), k, v)
    }
}

func (this KeyValue) DumpToSql() {
    this.Println()
}

// Self printable
func (this KeyValues) Println() {
    for k, v := range this {
        fmt.Printf("%s %v %v\n", T.CallerFuncName(1), k, v)
    }
}

func (this KeyValues) Append(key interface{}, val interface{}) {
    if _, ok := this[key]; !ok {
        this[key] = make([]interface{}, 1)
        this[key][0] = val
    } else {
        this[key] = append(this[key], val)
    }
}

func (this KeyValues) AppendSlice(key interface{}, val []interface{}) {
    if _, ok := this[key]; !ok {
        this[key] = make([]interface{}, 1)
        this[key] = val
    } else {
        this[key] = append(this[key], val...)
    }
}
