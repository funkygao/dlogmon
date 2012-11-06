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

func (this KeyValue) exportForGroupped(printer Printer, top int) {
    defer T.Un(T.Trace(""))

    for _, group := range this.Groups() {
        // header for each key type
        this.printGroupHeader(group)

        kvGroup := this.newByGroup(group) // a new kv just for this group
        if p, ok := printer.(Printher); ok {
            p.Printh(kvGroup, top)
        }

        println()
    }
}

func (this KeyValue) printGroupHeader(group string) {
    defer T.Un(T.Trace(""))

    fmt.Println(group)
    fmt.Println(strings.Repeat("=", GROUP_HEADER_LEN))
}

// this with key as mappers' output keys
// and value as reducer output value(KeyValue)
func (this KeyValue) ExportResult(printer Printer, top int) {
    defer T.Un(T.Trace(""))

    println("\n") // seperate from the progress bar

    if !this.Groupped() {
        this.exportForNonGrouped(printer, top)
        return
    } else {
        this.exportForGroupped(printer, top)
    }

}
