// Measure the performance of chan in go
package main

import (
    "fmt"
    "runtime"
    "strings"
    "time"
)

const (
    LOOPS = 1000000
)

type param struct {
    chbuf, datasize int
    yield bool
}

func timeit(f func(param), p param) time.Duration {
    start := time.Now()

    f(p) // call the func

    end := time.Now()
    delta := end.Sub(start)
    return delta
}

func bench(p param) {
    c := make(chan []byte, p.chbuf)
    data := make([]byte, p.datasize)

    go func() {
        for i:=0; i<LOOPS; i++ {
            c <- data

            if p.yield {
                runtime.Gosched()
            }
        }

        close(c)
    }()

    for {
        if _, ok := <- c; !ok {
            break
        }
    }

}

func main() {
    println("Loops:", LOOPS)
    println(strings.Repeat("=", 50))

    runBenches()
}

func runBenches() {
    for _, d := range []int{1, 128, 1024, 1024*1024} {
        fmt.Println("datasize:", d)
        for _, b := range []int{0, 1, 10, 1000, 10000} {
            for _, y := range []bool{true, false} {
                t := timeit(bench, param{b, d, y})
                fmt.Printf("\t%10s  %5s: buf:%5d, yield:%v\n", t, t/LOOPS, b, y)
            }
        }
    }
}
