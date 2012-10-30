package mr

import (
    "fmt"
    T "kx/trace"
)

// Factory
func NewKeyValue() KeyValue {
    return make(KeyValue)
}

// Self printable
func (this KeyValue) Println() {
    for k, v := range this {
        fmt.Printf(PRINT_FMT, T.CallerFuncName(1), k, v)
    }
}

// Self printable
func (this KeyValue) PrintByOrderedKey(sortedKeys interface{}) {
    for _, k := range sortedKeys.([]string) {
        fmt.Printf(PRINT_FMT, T.CallerFuncName(1), k, this[k])
    }
}

func (this KeyValue) Empty() bool {
    return len(this) == 0
}

func (this KeyValue) DumpToSql(sortedKeys interface{}) {
}
