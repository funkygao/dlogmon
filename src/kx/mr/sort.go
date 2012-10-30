package mr

import (
    "sort"
)

func newKvSort(kv KeyValue) *Sorter {
    s := &Sorter{
        keys: make([]interface{}, len(kv)),
        vals: make([]interface{}, len(kv)),
    }
    var i int
    for k, v := range kv {
        s.keys[i] = k
        s.vals[i] = v

        i++
    }

    return s
}

func (this *Sorter) Len() int {
    return len(this.keys)
}

func (this *Sorter) Less(i, j int) bool {
    if this.t == SORT_BY_KEY {
    } else if this.t == SORT_BY_VALUE {
    } else {
        panic("invalid sort type")
    }
    return true
}

func (this *Sorter) Swap(i, j int) {
    this.keys[i], this.keys[j] = this.keys[j], this.keys[i]
    this.vals[i], this.vals[j] = this.vals[j], this.vals[i]
}

func (this *Sorter) Sort(t SortType) {
    this.t = t
    sort.Sort(this)
}
