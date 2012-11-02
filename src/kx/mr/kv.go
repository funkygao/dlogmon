package mr

import (
    "encoding/gob"
    "fmt"
    "os"
    "reflect"
    "strings"
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

// Get a key
func (this KeyValue) getOneKey() (key interface{}) {
    for k, _ := range this {
        key = k
        return
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

// Is my keys groupped?
func (this KeyValue) Groupped() bool {
    return len(this.Groups()) > 0
}

// 如果要分组，那么mapper输出的key，必须是[]string，而且第一个key值为KEY_GROUP
// 第二个key值为group value
func (this KeyValue) Groups() []string {
    t := make(KeyValue)
    for k,_ := range this {
        switch k.(type) {
        case [2]string:
            if k.([2]string)[0] == KEY_GROUP {
                t[k.([2]string)[1]] = true
            }
        case [3]string:
            if k.([3]string)[0] == KEY_GROUP {
                t[k.([3]string)[1]] = true
            }
        case [4]string:
            if k.([4]string)[0] == KEY_GROUP {
                t[k.([4]string)[1]] = true
            }
        case [5]string:
            if k.([5]string)[0] == KEY_GROUP {
                t[k.([5]string)[1]] = true
            }
        case [6]string:
            if k.([6]string)[0] == KEY_GROUP {
                t[k.([6]string)[1]] = true
            }
        case [7]string:
            if k.([7]string)[0] == KEY_GROUP {
                t[k.([7]string)[1]] = true
            }
        case [8]string:
            if k.([8]string)[0] == KEY_GROUP {
                t[k.([8]string)[1]] = true
            }
        }
    }

    return InterfaceArrayToStringSlice(t.Keys())
}

func (this KeyValue) exportForNonGrouped(printer Printer, top int) {
    s := NewSort(this)
    s.Sort(this.sortType(), SORT_ORDER_DESC)
    sortedKeys := s.keys
    if top > 0 && top < len(sortedKeys) {
        sortedKeys = sortedKeys[:top]
    }

    for _, k := range sortedKeys {
        _ = printer.Printr(k, this[k].(KeyValue)) // return sql dml statement, usually 'insert into'
    }
}

func (this KeyValue) printGroupHeader(group string) {
    fmt.Println(group)
    fmt.Println(strings.Repeat("=", GROUP_HEADER_LEN))
}

func (this KeyValue) newByGroup(group string) KeyValue {
    if !this.Groupped() {
        return this
    }

    r := NewKeyValue()
    for k, v := range this {
        // discard the first [2]string with 0:KEY_GROUP 1:group name
        if key, ok := k.([3]string); ok && key[1] == group {
            r[[...]string{key[2]}] = v
        } else if key, ok := k.([4]string); ok && key[1] == group {
            r[[...]string{key[2], key[3]}] = v
        } else if key, ok := k.([5]string); ok && key[1] == group {
            r[[...]string{key[2], key[3], key[4]}] = v
        } else if key, ok := k.([6]string); ok && key[1] == group {
            r[[...]string{key[2], key[3], key[4], key[5]}] = v
        } else if key, ok := k.([7]string); ok && key[1] == group {
            r[[...]string{key[2], key[3], key[4], key[5], key[6]}] = v
        } else if key, ok := k.([8]string); ok && key[1] == group {
            r[[...]string{key[2], key[3], key[4], key[5], key[6], key[7]}] = v
        }
    }
    return r
}

// this with key as mappers' output keys
// and value as reducer output value(KeyValue)
func (this KeyValue) ExportResult(printer Printer, top int) {
    println("\n") // seperate from the progress bar

    if !this.Groupped() {
        this.exportForNonGrouped(printer, top)
        return
    }

    for _, group := range this.Groups() {
        // header for each key type
        this.printGroupHeader(group)

        kvGroup := this.newByGroup(group) // a new kv just for this group
        if p, ok := printer.(Printher); ok {
            // not only keys, values are also groupped
            p.Printh(kvGroup, top)
            continue
        }

        // only keys are groupped
        s := NewSort(kvGroup)
        s.Sort(SORT_BY_KEY, SORT_ORDER_DESC)
        sortedKeys := s.keys
        if top > 0 && top < len(sortedKeys) {
            sortedKeys = sortedKeys[:top]
        }
        for _, k := range sortedKeys {
            _ = printer.Printr(k, kvGroup[k].(KeyValue))
        }

        println()
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
