package main

import (
    "fmt"
    "testing"
    "reflect"
)

func main() {
    fmt.Println(testing.Benchmark(benchmarkDeepEqual))
}

func benchmarkDeepEqual(b *testing.B) {
    x, y := []string{"a", "b"}, []string{"c", "d", "e"}
    for i:=0; i<b.N; i++ {
        reflect.DeepEqual(x, y)
    }
}
