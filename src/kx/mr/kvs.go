package mr

func NewKeyValues() KeyValues {
    return make(KeyValues)
}

func (this KeyValues) Empty() bool {
    return len(this) == 0
}

func (this KeyValues) Append(key interface{}, val interface{}) {
    if _, ok := this[key]; !ok {
        this[key] = make([]interface{}, 1)
        this[key][0] = val
    } else {
        this[key] = append(this[key], val)
    }
}

func (this KeyValues) LaunchReducer(r Reducer) (out KeyValue) {
    out = NewKeyValue()

    // sort by key asc
    // the shuffling process
    s := NewSort(this)
    s.Sort(SORT_BY_KEY, SORT_ORDER_ASC)
    for _, k := range s.keys {
        // k is keys of mappers' output
        if v := r.Reduce(k, this[k]); v != nil && !v.Empty() {
            // v is output of reducer: KeyValue
            out[k] = v
        }
    }

    return
}
