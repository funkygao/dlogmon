package mr

import (
    "sort"
)

func newSort(p interface{}) *Sorter {
    if kv, ok := p.(KeyValue); ok {
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
    } else if kv, ok := p.(KeyValues); ok {
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
    } else {
        panic("invalid sorter type")
    }

    return nil
}

func (this *Sorter) Len() int {
    return len(this.keys)
}

func (this *Sorter) Less(i, j int) bool {
    if this.t == SORT_BY_KEY {
        ki, kj := this.keys[i].(string), this.keys[j].(string)
        return ki < kj
        
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

func (this *Sorter) Sort(t SortType, o SortOrdering) {
    this.t = t
    this.o = o

    sort.Sort(this)
}
