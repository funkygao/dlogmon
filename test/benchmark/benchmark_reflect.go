package main

import (
    "fmt"
    "testing"
    "reflect"
)

func main() {
    fmt.Println("deepEqual", testing.Benchmark(benchmarkDeepEqual))
    fmt.Println("valueOf", testing.Benchmark(benchmarkValueOf))
    fmt.Println("typeOf", testing.Benchmark(benchmarkTypeOf))
}

var foo = map[string]interface{} {"ff": 3.4}

func benchmarkDeepEqual(b *testing.B) {
    x, y := []string{"a", "b"}, []string{"c", "d", "e"}
    for i:=0; i<b.N; i++ {
        reflect.DeepEqual(x, y)
    }
}

func benchmarkValueOf(b *testing.B) {
    for i:=0; i<b.N; i++ {
        reflect.ValueOf(foo)
    }
}

func benchmarkTypeOf(b *testing.B) {
    for i:=0; i<b.N; i++ {
        reflect.TypeOf(foo)
    }
}
