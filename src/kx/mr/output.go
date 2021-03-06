package mr

import (
    "fmt"
    T "kx/trace"
    "strings"
)

func (this KeyValue) exportForNonGrouped(printer Printer, top int) {
    defer T.Un(T.Trace(""))

    s := NewSort(this)
    s.Sort(SORT_BY_VALUE, SORT_ORDER_DESC) // sort by value desc
    sortedKeys := s.keys
    if top > 0 && top < len(sortedKeys) {
        sortedKeys = sortedKeys[:top]
    }

    for _, k := range sortedKeys {
        _ = printer.(Printrer).Printr(k, this[k].(KeyValue)) // return sql dml statement, usually 'insert into'
    }
}

func (this KeyValue) exportForGroupped(printer Printer, group, sortCol string, top int) {
    defer T.Un(T.Trace(""))

    for _, grp := range this.Groups() {
        if group != "" && grp != group {
            continue
        }

        kvGroup := this.newByGroup(grp) // a new kv just for this group
        kvGroup.OutputGroup(printer, grp, sortCol, top)
        println()
    }
}

// this with key as mappers' output keys
// and value as reducer output value(KeyValue)
func (this KeyValue) ExportResult(printer Printer, group, sortCol string, top int) {
    defer T.Un(T.Trace(""))

    if !this.Groupped() {
        this.exportForNonGrouped(printer, top)
        return
    } else {
        this.exportForGroupped(printer, group, sortCol, top)
    }

}

func (kv KeyValue) OutputGroup(printer Printer, group, sortCol string, top int) {
    defer T.Un(T.Trace(""))

    // print group title
    fmt.Println(group)
    fmt.Println(strings.Repeat("-", OUTPUT_GROUP_HEADER_LEN))

    // output the aggregate columns title
    oneVal := kv.OneValue().(KeyValue)
    valKeys := oneVal.Keys()
    keyLengths := printer.(KeyLengther).KeyLengths(group)
    var keyLen int // key placeholder len total
    for _, l := range keyLengths {
        keyLen += l
    }
    fmt.Printf("%*d#", keyLen-1, len(kv))
    // default sort column
    if sortCol == "" {
        sortCol = valKeys[0].(string)
    }
    for _, x := range valKeys {
        if x == sortCol {
            x = x.(string) + "*"
        }
        fmt.Printf("%*s", OUTPUT_VAL_WIDTH, x)
    }

    // title done
    println()

    // sort by column
    s := NewSort(kv)
    s.SortCol(sortCol)
    s.Sort(SORT_BY_COL, SORT_ORDER_DESC)
    sortedKeys := s.Keys()
    if top > 0 && top < len(sortedKeys) {
        sortedKeys = sortedKeys[:top]
    }

    // output each key's values per line
    for _, sk := range sortedKeys {
        mapKey := sk.(GroupKey)
        // the keys
        for i, k := range mapKey.Keys() {
            if len(k) >= keyLengths[i] {
                k = k[:keyLengths[i]-1]
            }
            fmt.Printf("%*s", keyLengths[i], k)
        }

        // the values
        val := kv[sk].(KeyValue)
        for _, k := range valKeys {
            fmt.Printf("%*.1f", OUTPUT_VAL_WIDTH, val[k])
        }

        println()
    }
}
