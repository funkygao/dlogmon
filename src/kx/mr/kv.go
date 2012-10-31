package mr

// Factory
func NewKeyValue() KeyValue {
    return make(KeyValue)
}

func (this KeyValue) Empty() bool {
    return len(this) == 0
}

func (this KeyValue) ExportResult(printer Printer, top int) {
    s := newSort(this)
    // sort by value desc
    s.Sort(SORT_BY_VALUE, SORT_ORDER_DESC)
    sortedKeys := s.keys
    if top > 0 {
        sortedKeys = sortedKeys[:top]
    }
    println()
    for _, k := range sortedKeys {
        _ = printer.Printr(k, this[k]) // return sql dml statement, usually 'insert into'
    }
}
