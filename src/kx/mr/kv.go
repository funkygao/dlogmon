package mr

import (
    "encoding/gob"
    "os"
    "reflect"
)

// Factory
func NewKeyValue() KeyValue {
    return make(KeyValue)
}

func (this KeyValue) Empty() bool {
    return len(this) == 0
}

func (this KeyValue) sortType() SortType {
    var key interface{}
    for k, _ := range this {
        key = k
        break
    }

    switch v := reflect.ValueOf(key); v.Kind() {
    case reflect.Array:
        var possibleSortType string
        if t, ok := v.Interface().([2]string); ok {
            possibleSortType = t[0]
        } else if t, ok := v.Interface().([3]string); ok {
            possibleSortType = t[0]
        } else if t, ok := v.Interface().([4]string); ok {
            possibleSortType = t[0]
        } else if t, ok := v.Interface().([5]string); ok {
            possibleSortType = t[0]
        } else if t, ok := v.Interface().([6]string); ok {
            possibleSortType = t[0]
        } else if t, ok := v.Interface().([7]string); ok {
            possibleSortType = t[0]
        } else if t, ok := v.Interface().([8]string); ok {
            possibleSortType = t[0]
        }

        if possibleSortType == KEY_SECONDARY_KV {
            return SORT_SECONDARY_KV
        } else if possibleSortType == KEY_SECONDARY_VK {
            return SORT_SECONDARY_VK
        }
    }

    // default sort type
    return SORT_BY_VALUE
}

func (this KeyValue) ExportResult(printer Printer, top int) {
    s := newSort(this)
    s.Sort(this.sortType(), SORT_ORDER_DESC)
    sortedKeys := s.keys
    if top > 0 && top < len(sortedKeys) {
        sortedKeys = sortedKeys[:top]
    }

    println()
    for _, k := range sortedKeys {
        _ = printer.Printr(k, this[k]) // return sql dml statement, usually 'insert into'
    }
}

func (this KeyValue) serialize(filename string) {
    file, e := os.OpenFile(filename, GOB_FILE_FLAG, GOB_FILE_PERM)
    if e != nil {
        panic(e)
    }
    defer file.Close()

    enc := gob.NewEncoder(file)
    if e := enc.Encode(this); e != nil {
        // TODO
        // type not registered for interface: [3]string
        panic(e)
    }
}
