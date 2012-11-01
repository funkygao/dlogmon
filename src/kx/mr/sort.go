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

func (this Sorter) byKey() bool {
    return this.t == SORT_BY_KEY
}

func (this Sorter) asc() bool {
    return this.o == SORT_ORDER_ASC
}

func (this Sorter) lessStringSlice(si, sj []string) bool {
    if len(si) != len(sj) {
        if this.asc() {
            return len(si) < len(sj)
        } else {
            return len(si) > len(sj)
        }
    }

    for i := range si {
        vi, vj := si[i], sj[i]
        if vi < vj {
            if this.asc() {
                return true
            } else {
                return false
            }
        } else if vi > vj {
            if this.asc() {
                return false
            } else {
                return true
            }
        }
    }

    return false
}

func (this *Sorter) Less(i, j int) bool {
    if this.t == SORT_BY_KEY {
        ki := this.keys[i]

        switch ki.(type) {
        case string:
            ki, kj := this.keys[i].(string), this.keys[j].(string)
            if this.asc() {
                return ki < kj
            } else {
                return ki > kj
            }
        case [2]string:
            ki, kj := this.keys[i].([2]string), this.keys[j].([2]string)
            return this.lessStringSlice(ki[:], kj[:])
        case [3]string:
            ki, kj := this.keys[i].([3]string), this.keys[j].([3]string)
            return this.lessStringSlice(ki[:], kj[:])
        case [4]string:
            ki, kj := this.keys[i].([4]string), this.keys[j].([4]string)
            return this.lessStringSlice(ki[:], kj[:])
        case [5]string:
            ki, kj := this.keys[i].([5]string), this.keys[j].([5]string)
            return this.lessStringSlice(ki[:], kj[:])
        case [6]string:
            ki, kj := this.keys[i].([6]string), this.keys[j].([6]string)
            return this.lessStringSlice(ki[:], kj[:])
        case [7]string:
            ki, kj := this.keys[i].([7]string), this.keys[j].([7]string)
            return this.lessStringSlice(ki[:], kj[:])
        case [8]string:
            ki, kj := this.keys[i].([8]string), this.keys[j].([8]string)
            return this.lessStringSlice(ki[:], kj[:])
        }
    } else if this.t == SORT_BY_VALUE {
        // for reducer result, key is mappers' output key
        // and value is reducers' output KeyValue
        valsI, valsJ := this.vals[i].(KeyValue).Values(), this.vals[j].(KeyValue).Values()
        if len(valsI) != len(valsJ) {
            panic("value length not match")
        }
        // 这里的限制是，len(valsI) == 1
        rvi, rvj := valsI[0], valsJ[0]
        switch rvi.(type) {
        case string:
            rvi, rvj := rvi.(string), rvj.(string)
            if this.asc() {
                return rvi < rvj
            } else {
                return rvi > rvj
            }
        case float64:
            rvi, rvj := rvi.(float64), rvj.(float64)
            if this.asc() {
                return rvi < rvj
            } else {
                return rvi > rvj
            }
        case int:
            rvi, rvj := rvi.(int), rvj.(int)
            if this.asc() {
                return rvi < rvj
            } else {
                return rvi > rvj
            }
        case int64:
            rvi, rvj := rvi.(int64), rvj.(int64)
            if this.asc() {
                return rvi < rvj
            } else {
                return rvi > rvj
            }
        }
    } else {
        panic("invalid sort type")
    }
    return false
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

func (this Sorter) Keys() []interface{} {
    return this.keys
}
