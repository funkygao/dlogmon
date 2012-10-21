package main

import (
    "fmt"
    "testing"
)

var ops uint64

func main() {
    fmt.Printf("%5s %s\n", "+", testing.Benchmark(benchmarkAdd).String())
    fmt.Printf("%5s %s\n", "-", testing.Benchmark(benchmarkMinus).String())
    fmt.Printf("%5s %s\n", "*", testing.Benchmark(benchmarkProduct).String())
    fmt.Printf("%5s %s\n", "/", testing.Benchmark(benchmarkDivide).String())
    fmt.Printf("%5s %s\n", "if", testing.Benchmark(benchmarkIf).String())
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
