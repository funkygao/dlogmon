package main

import (
    "fmt"
    "testing"
    "time"
)

func main() {
    fmt.Printf("%s\n", testing.Benchmark(benchmarkNow).String())
    fmt.Printf("%s\n", testing.Benchmark(benchmarkSleep).String())
}

func benchmarkNow(b *testing.B) {
    for i:=0; i<b.N; i++ {
        time.Now()
    }
}

func benchmarkSleep(b *testing.B) {
    for i:=0; i<b.N; i++ {
        time.Sleep(0)
    }
}
