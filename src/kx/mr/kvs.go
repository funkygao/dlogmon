package mr

func NewKeyValues() KeyValues {
    return make(KeyValues)
}

func (this KeyValues) Empty() bool {
    return len(this) == 0
}

func (this KeyValues) Append(key interface{}, val ...interface{}) {
    if _, ok := this[key]; !ok {
        this[key] = make([]interface{}, 0, 100)
    }
    this[key] = append(this[key], val...)
}
