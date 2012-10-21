package main

import (
    "fmt"
    "sync/atomic"
    "time"
    "testing"
)

var ops uint64

func main() {
    fmt.Printf("%12d %s\n", ops, testing.Benchmark(benchmarkCAS).String())
    final := atomic.LoadUint64(&ops)
    time.Sleep(time.Second)
    fmt.Printf("%12d\n", final)
}

func benchmarkCAS(b *testing.B) {
    for i:=0; i<b.N; i++ {
        atomic.AddUint64(&ops, 1)
    }
}
