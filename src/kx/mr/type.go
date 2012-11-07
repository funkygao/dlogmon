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

type Filter interface {
    IsLineValid(string) bool
}

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

type GroupKey struct {
    group string
    Key
}

type Key struct {
    key string // gob'ed slice of strings
}

type Grouper interface {
    Group() string
}

type Printer interface {
}

type Printrer interface {
    Printr(key interface{}, value KeyValue) string
}

type KeyLengther interface {
    KeyLengths(group string) []int
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
    col  interface{}
    sort.Interface
}
