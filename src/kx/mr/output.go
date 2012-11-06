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
        if group != "" && grp !=group {
            continue
        }
        kvGroup := this.newByGroup(grp) // a new kv just for this group
        kvGroup.OutputGroup(printer, grp, "", top)
        println()
    }
}

// this with key as mappers' output keys
// and value as reducer output value(KeyValue)
func (this KeyValue) ExportResult(printer Printer, group, sortCol string, top int) {
    defer T.Un(T.Trace(""))

    println("\n") // seperate from the progress bar

    if !this.Groupped() {
        this.exportForNonGrouped(printer, top)
        return
    } else {
        this.exportForGroupped(printer, group, sortCol, top)
    }

}

func (kv KeyValue) OutputGroup(printer Printer, group, sortCol string, top int) {
    defer T.Un(T.Trace(""))

    // print group header
    fmt.Println(group)
    fmt.Println(strings.Repeat("=", OUTPUT_GROUP_HEADER_LEN))

    // output the aggregate columns header
    oneVal := kv.OneValue().(KeyValue)
    valKeys := oneVal.Keys()
    keyLengths := printer.(KeyLengther).KeyLengths(group)
    for _, l := range keyLengths {
        fmt.Printf("%*s", l, "")
    }
    for _, x := range valKeys {
        fmt.Printf("%*s", OUTPUT_VAL_WIDTH, x)
    }
    println()

    // sort by column
    s := NewSort(kv)
    if sortCol != "" {
        s.SortCol(sortCol)
        println("ha", sortCol)
    } else {
        // default sort by 1st colomn
        s.SortCol(valKeys[0])
    }
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
            if len(k) > keyLengths[i] {
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
