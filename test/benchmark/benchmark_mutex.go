package main

import (
    "fmt"
    "sync"
    "testing"
)

func main() {
    fmt.Printf("%s\n", testing.Benchmark(benchmarkMutexLock).String())
}

func benchmarkMutexLock(b *testing.B) {
    var lock sync.Mutex
    for i:=0; i<b.N; i++ {
        lock.Lock()
        lock.Unlock()
    }
}
