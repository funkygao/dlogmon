package main

import (
    "fmt"
    "runtime"
    "testing"
)

var ops uint64

func main() {
    fmt.Printf("%5s %s\n", "loop", testing.Benchmark(benchmarkEmptyLoop).String())
    fmt.Printf("%5s %s\n", "+", testing.Benchmark(benchmarkAdd).String())
    fmt.Printf("%5s %s\n", "-", testing.Benchmark(benchmarkMinus).String())
    fmt.Printf("%5s %s\n", "*", testing.Benchmark(benchmarkProduct).String())
    fmt.Printf("%5s %s\n", "/", testing.Benchmark(benchmarkDivide).String())
    fmt.Printf("%5s %s\n", "if", testing.Benchmark(benchmarkIf).String())
    fmt.Printf("%5s %s\n", "print", testing.Benchmark(benchmarkPrint).String())
    fmt.Printf("%5s %s\n", "go", testing.Benchmark(benchmarkGo).String())
    fmt.Printf("%5s %s\n", "sched", testing.Benchmark(benchmarkGosched).String())
}

func benchmarkPrint(b *testing.B) {
    for i:=0; i<b.N; i++ {
        print("")
    }
}

func benchmarkGosched(b *testing.B) {
    for i:=0; i<b.N; i++ {
        runtime.Gosched()
    }
}

func benchmarkGo(b *testing.B) {
    for i:=0; i<b.N; i++ {
        go func() {
        }()
    }
}

func benchmarkEmptyLoop(b *testing.B) {
    for i:=0; i<b.N; i++ {
    }
}

func benchmarkAdd(b *testing.B) {
    var x uint64
    for i:=0; i<b.N; i++ {
        x += 1
    }
}

func benchmarkProduct(b *testing.B) {
    var x uint64 = 2
    for i:=0; i<b.N; i++ {
        x *= 2
    }
}

func benchmarkDivide(b *testing.B) {
    var x uint64 = 2 << 60
    for i:=0; i<b.N; i++ {
        x /= 2
    }
}

func benchmarkMinus(b *testing.B) {
    var x uint64 = 2 << 60
    for i:=0; i<b.N; i++ {
        x -= 2
    }
}

func benchmarkIf(b *testing.B) {
    var x uint64 = 2 << 60
    for i:=0; i<b.N; i++ {
        if x < 554454 {
        }
    }
}
