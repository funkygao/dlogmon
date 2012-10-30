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

// Mapper
type Mapper interface {
    Map(line string, out chan<- KeyValue)
}

// Reducer
type Reducer interface {
    Reduce(key interface{}, values []interface{}) (out KeyValue)
}

type SortType uint8

type Sorter struct {
    keys []interface{}
    vals []interface{}
    t SortType
    sort.Interface
}
