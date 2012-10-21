package main

import (
    "fmt"
    "sync/atomic"
    "testing"
)

var ops uint64

func main() {
    fmt.Printf("%12d %s\n", ops, testing.Benchmark(benchmarkChan).String())
}

func benchmarkChan(b *testing.B) {
    for i:=0; i<b.N; i++ {
        atomic.AddUint64(&ops, 1)
    }
}
