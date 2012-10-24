package main

import (
    "fmt"
    "testing"
)

func main() {
    fmt.Println(testing.Benchmark(benchmarkSliceAppend))
}

func benchmarkSliceAppend(b *testing.B) {
    slice := make([]int, 0)
    for i:=0; i<b.N; i++ {
        slice = append(slice, i)
    }
}
