package main

import (
    "fmt"
    "testing"
    "time"
)

func main() {
    fmt.Printf("%s\n", testing.Benchmark(benchmarkNow).String())
}

func benchmarkNow(b *testing.B) {
    for i:=0; i<b.N; i++ {
        time.Now()
    }
}
