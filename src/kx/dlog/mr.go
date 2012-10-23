package dlog

import (
    "fmt"
)

func newMapData() MapData {
    return make(MapData)
}

func newTransformData() TransformData {
    return make(TransformData)
}

func newReduceData(size int) ReduceData {
    r := make(ReduceData, size)
    for i:=0; i<size; i++ {
        r[i] = newTransformData()
    }
    return r
}

func newReduceResult(size int) ReduceResult {
    r := make(ReduceResult, size)
    for i:=0; i<size; i++ {
        r[i] = newMapData()
    }
    return r
}

func getKeyByType(t KeyType, key string) string {
    return fmt.Sprintf("%d%s%s", t, KEYTYPE_SEP, key)
}

func getKeyType(key string) (r KeyType, k string) {
    format := "%d" + KEYTYPE_SEP + "%s"
    fmt.Sscanf(key, format, &r, &k)
    return
}

func (this MapData) Println() {
    for k, v := range this {
        fmt.Println(k, v)
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

func (this TransformData) Println() {
    for k, v := range this {
        fmt.Println(k, v)
    }
}

// Get key types into slice of KeyType
func (this TransformData) KeyTypes() (r []KeyType) {
    var m = make(map[KeyType] bool)
    for k, _ := range this {
        key, _ := getKeyType(k)
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

func (this ReduceData) Println() {
    for keyType, d := range this {
        println("\nKeyType:", keyType)
        d.Println()
    }
}

func (this ReduceResult) Println() {
    for keyType, d := range this {
        println("\nKeyType:", keyType)
        d.Println()
    }
}
