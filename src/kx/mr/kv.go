package mr

import (
    "encoding/gob"
    "kx/util"
    "os"
)

func NewKeyValue() KeyValue {
    return make(KeyValue)
}

func (this KeyValue) Empty() bool {
    return len(this) == 0
}

func (this KeyValue) getOneKey() (key interface{}) {
    for k := range this {
        key = k
        return
    }
    return
}

func (this KeyValue) Keys() (keys []interface{}) {
    keys = make([]interface{}, len(this))
    var i int
    for k := range this {
        keys[i] = k
        i++
    }
    return
}

func (this KeyValue) Values() (values []interface{}) {
    values = make([]interface{}, len(this))
    var i int
    for _, v := range this {
        values[i] = v
        i++
    }
    return
}

// Is my keys groupped?
func (this KeyValue) Groupped() bool {
    k := this.getOneKey()
    if _, ok := k.(Grouper); ok {
        return true
    }
    return false
}

// 如果要分组，那么mapper输出的key，必须是[]string，而且第一个key值为KEY_GROUP
// 第二个key值为group value
func (this KeyValue) Groups() []string {
    r := util.NewSet()
    for k := range this {
        if g, ok := k.(GroupKey); ok {
            r.Add(g.Group())
        }
    }

    return r.StrValues()
}

func (this KeyValue) Emit(ch chan<- KeyValue) {
    ch <- this
}

func (this KeyValue) newByGroup(group string) KeyValue {
    if !this.Groupped() {
        return this
    }

    r := NewKeyValue()
    for k, v := range this {
        key := k.(GroupKey)
        if key.group == group {
            r[key] = v
        }
    }
    return r
}

func (this KeyValue) Save(filename string) {
    file, e := os.OpenFile(filename, GOB_FILE_FLAG, GOB_FILE_PERM)
    if e != nil {
        panic(e)
    }
    defer file.Close()

    enc := gob.NewEncoder(file)
    if e := enc.Encode(this); e != nil {
        // TODO
        // type not registered for interface: [3]string
        panic(e)
    }
}

// TODO decode gob
func (this *KeyValue) Load(filename string) {
}
