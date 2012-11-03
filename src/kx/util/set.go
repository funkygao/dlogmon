package util

type Set struct {
    data map[interface{}]bool
}

func NewSet() *Set {
    s := new(Set)
    s.data = make(map[interface{}]bool)
    return s
}

func (this *Set) Add(v interface{}) {
    this.data[v] = true
}

func (this *Set) Remove(v interface{}) (found bool) {
    if _, found = this.data[v]; !found {
        return
    }

    delete(this.data, v)
    found = true
    return
}

func (this Set) Len() int {
    return len(this.data)
}

func (this Set) Contains(v interface{}) bool {
    if _, found := this.data[v]; found {
        return true
    }
    return false
}

func (this Set) Values() []interface{} {
    r := make([]interface{}, len(this.data))
    var i int
    for k := range this.data {
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
