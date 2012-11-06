// TODO we need only save key or value in Sorter
package mr

import (
    "sort"
)

func NewSort(p interface{}) *Sorter {
    if kv, ok := p.(KeyValue); ok {
        s := &Sorter{
            keys: make([]interface{}, len(kv)),
            vals: make([]interface{}, len(kv)),
        }
        var i int
        for k, v := range kv {
            s.keys[i], s.vals[i] = k, v

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
            s.keys[i], s.vals[i] = k, v

            i++
        }

        return s
    } else {
        panic("invalid sorter type")
    }

    return nil
}

func (this Sorter) asc() bool {
    return this.o == SORT_ORDER_ASC
}

func (this *Sorter) Len() int {
    return len(this.keys)
}

func (this *Sorter) Less(i, j int) bool {
    if this.t == SORT_BY_KEY {
        ki := this.keys[i]

        switch ki.(type) {
        case GroupKey:
            ki, kj := this.keys[i].(GroupKey), this.keys[j].(GroupKey)
            return ki.Less(kj, this.asc())
        case Key:
            ki, kj := this.keys[i].(Key), this.keys[j].(Key)
            return ki.Less(kj, this.asc())
        case string:
            ki, kj := this.keys[i].(string), this.keys[j].(string)
            return less(ki, kj, this.asc())
        default:
            panic("Invalid key type")
        }
    } else if this.t == SORT_BY_VALUE {
        // for reducer result, key is mappers' output key
        // and value is reducers' output KeyValue
        valsI, valsJ := this.vals[i].(KeyValue).Values(), this.vals[j].(KeyValue).Values()
        if len(valsI) != 1 || len(valsJ) != 1 {
            panic("len must be 1")
        }
        return less(valsI[0], valsJ[0], this.asc())
    } else if this.t == SORT_BY_COL {
        if this.col == nil {
            panic("Must call SortCol(col) before Sort()")
        }
        vi, vj := this.vals[i].(KeyValue), this.vals[j].(KeyValue)
        return less(vi[this.col], vj[this.col], this.asc())
    } else {
        panic("invalid sort type")
    }
    return false
}

func (this *Sorter) Swap(i, j int) {
    this.keys[i], this.keys[j] = this.keys[j], this.keys[i]
    this.vals[i], this.vals[j] = this.vals[j], this.vals[i]
}

func (this *Sorter) SortCol(col interface{}) {
    this.col = col
}

func (this *Sorter) Sort(t SortType, o SortOrdering) {
    this.t = t
    this.o = o

    sort.Sort(this)
}

func (this Sorter) Keys() []interface{} {
    return this.keys
}
