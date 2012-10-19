// Measure the performance of chan in go
package main

import (
    "fmt"
    "strings"
    "time"
)

const (
    LOOPS = 1000000
)

type param struct {
    chbuf, datasize int
}

func timeit(f func(param), p param) time.Duration {
    start := time.Now()

    f(p) // call the func

    end := time.Now()
    delta := end.Sub(start)
    return delta
}

func bench(p param) {
    c, data := make(chan []byte, p.chbuf), make([]byte, p.datasize)

    go func() {
        for i:=0; i<LOOPS; i++ {
            c <- data
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
            t := timeit(bench, param{b, d})
            fmt.Printf("\t%10s  %5s: buf:%5d\n", t, t/LOOPS, b)
        }
    }
}
