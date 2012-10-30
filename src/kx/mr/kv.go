package mr

// Factory
func NewKeyValue() KeyValue {
    return make(KeyValue)
}

func (this KeyValue) Empty() bool {
    return len(this) == 0
}

func (this KeyValue) ExportResult(printer Printer) {
    s := newSort(this)
    s.Sort(SORT_BY_KEY, SORT_ORDER_DESC)
    println()
    for _, k := range s.keys {
        _ = printer.Printr(k, this[k]) // return sql dml statement, usually 'insert into'
    }
}
