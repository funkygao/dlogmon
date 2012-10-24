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
    for i:=0; i<size; i++ {
        r[i] = NewTransformData()
    }
    return r
}

// Factory
func NewReduceResult(size int) ReduceResult {
    r := make(ReduceResult, size)
    for i:=0; i<size; i++ {
        r[i] = NewMapData()
    }
    return r
}

func getKeyByType(t KeyType, key string) string {
    return fmt.Sprintf("%d%s%s", t, KEYTYPE_SEP, key)
}

func GetKeyType(key string) (r KeyType, k string) {
    format := "%d" + KEYTYPE_SEP + "%s"
    fmt.Sscanf(key, format, &r, &k)
    return
}

// Self printable
func (this MapData) Println() {
    for k, v := range this {
        fmt.Println("mr", k, v)
    }
}

func (this MapData) Set(t KeyType, key string, val float64) {
    key = getKeyByType(t, key)
    this[key] = val
}

func (this MapData) Get(key string) (val float64, ok bool) {
    val, ok= this[key]
    return
}

func (this TransformData) Append(key string, val float64) {
    _, ok := this[key]
    if !ok {
        this[key] = make([]float64, 1)
        this[key][0] = val
    } else {
        this[key] = append(this[key], val)
    }
}

func (this TransformData) AppendSlice(key string, val []float64) {
    _, ok := this[key]
    if !ok {
        this[key] = make([]float64, 1)
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

// Get key types into slice of KeyType
func (this TransformData) KeyTypes() (r []KeyType) {
    var m = make(map[KeyType] bool)
    for k, _ := range this {
        key, _ := GetKeyType(k)
        m[key] = true
    }

    r = make([]KeyType, len(m))
    var i int
    for k, _ := range m {
        r[i] = k
        i ++
    }
    return
}

// Self printable
func (this ReduceData) Println() {
    for keyType, d := range this {
        println("\nmr KeyType:", keyType)
        d.Println()
    }
}

// Self printable
func (this ReduceResult) Println() {
    for keyType, d := range this {
        println("\nmr KeyType:", keyType)
        d.Println()
    }
}
