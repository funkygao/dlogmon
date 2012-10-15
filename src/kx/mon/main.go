// main entry
package main

import (
    "fmt"
    "kx/dlog"
    "runtime"
    "sync"
    "time"
)

func main() {
    // parallel level
    runtime.GOMAXPROCS((runtime.NumCPU() + 1)/2)

    // cli options
    options := parseFlags()

    chLines := make(chan int, len(options.files))
    lock := new(sync.Mutex)

    // timing all the jobs up
    start := time.Now()

    // each dlog file is a goroutine
    for _, file := range options.files {
        dlog := dlog.NewAmfDlog(file, chLines, lock)
        go dlog.ReadLines()
    }

    // wait for all dlog runner finish
    lines := 0
    for i:=0; i<len(options.files); i++ {
        lines += <- chLines
    }

    end := time.Now()
    delta := end.Sub(start)
    fmt.Printf("\nParsed %d lines in %s [%.1f lines per second]\n", lines, delta, float64(lines)/delta.Seconds())
}

