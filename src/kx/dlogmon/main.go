// main entry
package main

import (
    "fmt"
    "kx/dlog"
    "runtime"
    "sync"
    "time"
)

var kindMapping = map[string] func(string, chan int, *sync.Mutex, *dlog.Options) dlog.IDlogExecutor {
    "amf": dlog.NewAmfDlog}

func main() {
    // parallel level
    parallel := runtime.NumCPU()/2 + 1
    runtime.GOMAXPROCS(parallel)
    fmt.Printf("Parallel CPU: %d / %d\n", parallel, runtime.NumCPU())

    // cli options
    options := dlog.ParseFlags()
    files := options.GetFiles()

    chLines := make(chan int, len(files))
    lock := new(sync.Mutex)

    // timing all the jobs up
    start := time.Now()

    // each dlog file is a goroutine
    for _, file := range files {
        var executor dlog.IDlogExecutor
        executor = kindMapping[options.GetKind()](file, chLines, lock, options)
        go executor.ScanLines()
    }

    // wait for all dlog runner finish
    lines := 0
    for i:=0; i<len(files); i++ {
        lines += <- chLines
    }

    end := time.Now()
    delta := end.Sub(start)
    fmt.Printf("\nParsed %d lines in %d files within %s [%.1f lines per second]\n",
        lines,
        len(files),
        delta, float64(lines)/delta.Seconds())
}

