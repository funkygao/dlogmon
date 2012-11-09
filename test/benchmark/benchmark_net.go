package main

import (
    "fmt"
    "testing"
    "net"
)

func main() {
    fmt.Printf("%s\n", testing.Benchmark(benchmarkDialAndClose).String())
}

func benchmarkDialAndClose(b *testing.B) {
    for i:=0; i<b.N; i++ {
        conn, e := net.Dial("tcp", "localhost:80")
        if e != nil {
            panic(e)
        }
        conn.Close()
    }
}
