package mr

import (
    "fmt"
    "github.com/pmylund/sortutil"
    T "kx/trace"
)

const (
    PRINT_FMT = "%s %60v %v\n"
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
        fmt.Printf(PRINT_FMT, T.CallerFuncName(1), k, v)
    }
}

func (this KeyValue) DumpToSql() {
    println()
    this.Println()
}

// Self printable
func (this KeyValues) Println() {
    for k, v := range this {
        fmt.Printf(PRINT_FMT, T.CallerFuncName(1), k, v)
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

func (this KeyValues) Keys() []interface{} {
    keys := make([]interface{}, len(this))
    var i int
    for k, _ := range this {
        keys[i] = k
        i ++
    }
    return keys
}

func (this KeyValues) SortedKeys() []string {
    keys := make([]string, len(this))
    for i, v := range this.Keys() {
        keys[i] = v.(string)
    }
    fmt.Println(keys)
    sortutil.Desc(keys)
    fmt.Println(keys)

    return keys
}
