package main

import (
    "fmt"
    "testing"
)

var buf int

func main() {
    bs := []int{0, 1, 2, 5, 10, 50, 100, 200, 500, 1000, 5000}
    for _, s := range bs {
        buf = s
        fmt.Printf("%5d %s\n", buf, testing.Benchmark(benchmarkChan).String())
    }
}

func benchmarkChan(b *testing.B) {
    ch := make(chan int, buf)
    go func() {
        for i:=0; i<b.N; i++ {
            ch <- i
        }

        close(ch)
    }()

    for _ = range ch {
    }
}
