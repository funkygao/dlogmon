package mr

const KEYTYPE_SEP = ":"

// Local aggregator
type CombinerFunc func([]float64) float64

// TODO tag
type TagType uint8 

// Mapper raw output format
type MapData map[string]float64

// mapper -> TransformData -> reduce
type TransformData map[string][]float64

type ReduceData []TransformData

type ReduceResult []MapData

// map
type Mapper interface {
    Map(string, chan<- interface{})
}

// reduce
type Reducer interface {
    Reduce(ReduceData) ReduceResult
}
