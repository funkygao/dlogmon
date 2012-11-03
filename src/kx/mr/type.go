package mr

import (
    "sort"
)

// Local aggregator
// TODO combiner is not just a func, it's a mini reducer
type CombinerFunc func([]float64) float64

// key->value pair with key and value being any type
type KeyValue map[interface{}]interface{}

// key->[]value pair with key and value being any type
type KeyValues map[interface{}][]interface{}

type Mapper interface {
    Map(line string, out chan<- KeyValue)
}

type Reducer interface {
    Reduce(key interface{}, values []interface{}) (out KeyValue)
}

type MapReducer interface {
    Mapper
    Reducer
}

type Grouper interface {
    Groups() []string
}

type Printer interface {
    Printr(key interface{}, value KeyValue) string
}

type Printher interface {
    Printh(KeyValue, int)
}

type (
    SortType     uint8
    SortOrdering uint8
)

type Sorter struct {
    keys []interface{}
    vals []interface{}
    t    SortType
    o    SortOrdering
    sort.Interface
}
