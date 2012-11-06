package mr

// TODO group reducer into goroutines
func (this KeyValues) LaunchReducer(r Reducer) (out KeyValue) {
    out = NewKeyValue()

    // sort by key asc
    s := NewSort(this)
    s.Sort(SORT_BY_KEY, SORT_ORDER_ASC)
    for _, k := range s.keys {
        // k is keys of mappers' output key
        if kv := r.Reduce(k, this[k]); kv != nil && !kv.Empty() {
            out[k] = kv
        }
    }

    return
}
