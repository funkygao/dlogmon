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

func (this KeyValue) Keys() (keys []interface{}) {
    keys = make([]interface{}, len(this))
    var i int
    for k, _ := range this {
        keys[i] = k
        i++
    }
    return
}

func (this KeyValue) Values() (values []interface{}) {
    values = make([]interface{}, len(this))
    var i int
    for _, v := range this {
        values[i] = v
        i++
    }
    return
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

// Return []string
func (this KeyValue) KeyTypes() []string {
    t := make(KeyValue)
    for k,_ := range this {
        switch k.(type) {
        case [2]string:
            t[k.([2]string)[0]] = true
        case [3]string:
            t[k.([3]string)[0]] = true
        case [4]string:
            t[k.([4]string)[0]] = true
        case [5]string:
            t[k.([5]string)[0]] = true
        case [6]string:
            t[k.([6]string)[0]] = true
        case [7]string:
            t[k.([7]string)[0]] = true
        case [8]string:
            t[k.([8]string)[0]] = true
        }
    }

    return InterfaceArrayToStringSlice(t.Keys())
}

// this with key as mappers' output keys
// and value as reducer output value(KeyValue)
func (this KeyValue) ExportResult(printer Printer, top int) {
    for _, kt := range this.KeyTypes() {
        println(kt)
    }

    s := newSort(this)
    s.Sort(this.sortType(), SORT_ORDER_DESC)
    sortedKeys := s.keys
    if top > 0 && top < len(sortedKeys) {
        sortedKeys = sortedKeys[:top]
    }

    println() // seperate from the progress bar
    for _, k := range sortedKeys {
        _ = printer.Printr(k, this[k].(KeyValue)) // return sql dml statement, usually 'insert into'
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
