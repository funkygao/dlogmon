package main

import (
    "fmt"
    "testing"
)

func main() {
    fmt.Println("mem", testing.Benchmark(benchmarkMemAccess))
}

func benchmarkMemAccess(b *testing.B) {
    slice := make([]int, 10)
    slice[4] = 10
    b.StartTimer()
    for i:=0; i<b.N; i++ {
        _ = slice[4]
    }
}
