package mr

// Factory
func NewKeyValues() KeyValues {
    return make(KeyValues)
}

func (this KeyValues) getOneKey() (key interface{}) {
    for k := range this {
        key = k
        return
    }

    return
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

func (this KeyValues) AppendSlice(key interface{}, val []interface{}) {
    if _, ok := this[key]; !ok {
        this[key] = make([]interface{}, 1)
        this[key] = val
    } else {
        this[key] = append(this[key], val...)
    }
}

func (this KeyValues) LaunchReducer(r Reducer) (out KeyValue) {
    out = NewKeyValue()

    // sort by key asc
    // the shuffling process
    s := newSort(this)
    s.Sort(SORT_BY_KEY, SORT_ORDER_ASC)
    for _, k := range s.keys {
        if v := r.Reduce(k, this[k]); v != nil && !v.Empty() {
            out[k] = v
        }
    }

    return
}
