package mr

const KEYTYPE_SEP = ":"

// Local aggregator
type CombinerFunc func([]float64) float64

// TODO tag
type TagType uint8 

// Mapper raw output format
type MapData map[interface{}]interface{}

// mapper -> TransformData -> merge -> reduce
type TransformData map[interface{}][]interface{}

// Input of Reducer
type ReduceData []TransformData

// Output of Reducer
type ReduceResult []MapData

// map
type Mapper interface {
    Map(string, chan<- interface{})
}

// reduce
type Reducer interface {
    Reduce(ReduceData) ReduceResult
}
