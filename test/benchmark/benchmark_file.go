package main

import (
    "fmt"
    "testing"
    "os"
    "io/ioutil"
)

func main() {
    fmt.Printf("%s\n", testing.Benchmark(benchmarkReadFile).String())
}

func benchmarkReadFile(b *testing.B) {
    const F = "/Users/gaopeng/github/dlogmon/bin/dlogmon"

    for i:=0; i<b.N; i++ {
        ioutil.ReadFile(F)
    }

    f, _ := os.Open(F)
    st, _ := f.Stat()
    b.SetBytes(st.Size())
}
