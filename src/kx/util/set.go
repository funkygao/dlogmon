package util

type Set struct {
    data map[interface{}] bool
}

func (this *Set) Add(v interface{}) {
    this.data[v] = true
}

func (this *Set) Remove(v interface{}) (ok bool) {
    if _, ok = this.data[v]; !ok {
        return
    }

    delete(this.data, v)
    ok = true
    return
}

func (this Set) Values() []interface{} {
    r := make([]interface{}, len(this.data))
    var i int
    for k, _ := range this.data {
        r[i] = k
        i++
    }
    return r
}

func (this Set) StrValues() []string {
    r := make([]string, len(this.data))
    for i, v := range this.Values() {
        r[i] = v.(string)
    }

    return r
}
