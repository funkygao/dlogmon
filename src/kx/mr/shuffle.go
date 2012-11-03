package mr

// shuffle k=>v to k=>[v1, ...vn]
func Shuffle(in <-chan KeyValue, out chan<- KeyValues) {
    kvs := NewKeyValues()
    for x := range in {
        for k, v := range x {
            kvs.Append(k, v)
        }
    }

    out <- kvs
}
