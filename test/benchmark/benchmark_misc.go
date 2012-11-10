package main

import (
    "fmt"
    "testing"
)

func main() {
    fmt.Printf("method by ref %s\n", testing.Benchmark(benchmarkMethodByRef).String())
    fmt.Printf("method by val %s\n", testing.Benchmark(benchmarkMethodByVal).String())
    fmt.Printf("empty func %s\n", testing.Benchmark(benchmarkEmptyFunc).String())
}

func e() {
}

type Foo struct {}

func (this Foo) X() {
}

func (this *Foo) Y() {
}

var f = Foo{}

func benchmarkEmptyFunc(b *testing.B) {
    for i:=0; i<b.N; i++ {
        e()
    }
}

func benchmarkMethodByRef(b *testing.B) {
    for i:=0; i<b.N; i++ {
        f.Y()
    }
}

func benchmarkMethodByVal(b *testing.B) {
    for i:=0; i<b.N; i++ {
        f.X()
    }
}
