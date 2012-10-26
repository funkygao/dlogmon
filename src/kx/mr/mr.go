package mr

import (
    "fmt"
)

// Factory
func NewMapData() MapData {
    return make(MapData)
}

// Factory
func NewTransformData() TransformData {
    return make(TransformData)
}

// Factory
func NewReduceData(size int) ReduceData {
    r := make(ReduceData, size)
    for i := 0; i < size; i++ {
        r[i] = NewTransformData()
    }
    return r
}

// Factory
func NewReduceResult(size int) ReduceResult {
    r := make(ReduceResult, size)
    for i := 0; i < size; i++ {
        r[i] = NewMapData()
    }
    return r
}

func getTagByType(t TagType, key interface{}) string {
    return fmt.Sprintf("%d%s%s", t, KEYTYPE_SEP, key)
}

func GetTagType(key string) (t TagType, k string) {
    format := "%d" + KEYTYPE_SEP + "%s"
    fmt.Sscanf(key, format, &t, &k)
    return
}

// Self printable
func (this MapData) Println() {
    for k, v := range this {
        fmt.Println("mr", k, v)
    }
}

func (this MapData) Set(t TagType, key interface{}, val interface{}) {
    key = getTagByType(t, key)
    this[key] = val
}

func (this MapData) Get(key interface{}) (val interface{}, ok bool) {
    val, ok = this[key]
    return
}

func (this TransformData) Append(key interface{}, val interface{}) {
    _, ok := this[key]
    if !ok {
        this[key] = make([]interface{}, 1)
        this[key][0] = val
    } else {
        this[key] = append(this[key], val)
    }
}

func (this TransformData) AppendSlice(key interface{}, val []interface{}) {
    _, ok := this[key]
    if !ok {
        this[key] = make([]interface{}, 1)
        this[key] = val
    } else {
        this[key] = append(this[key], val...)
    }
}

// Self printable
func (this TransformData) Println() {
    for k, v := range this {
        fmt.Println("mr", k, v)
    }
}

// Get key types into slice of TagType
func (this TransformData) TagTypes() (r []TagType) {
    var m = make(map[TagType]bool)
    for k := range this {
        key, _ := GetTagType(k.(string))
        m[key] = true
    }

    r = make([]TagType, len(m))
    var i int
    for k := range m {
        r[i] = k
        i++
    }
    return
}

// Self printable
func (this ReduceData) Println() {
    for _, d := range this {
        d.Println()
    }
}

// Self printable
func (this ReduceResult) Println() {
    for _, d := range this {
        d.Println()
    }
}

// Dump to plain sql statements
// TODO
func (this ReduceResult) DumpToSql() {
    this.Println()
}
