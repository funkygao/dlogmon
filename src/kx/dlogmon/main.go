// main entry
package main

import (
    "fmt"
    "kx/dlog"
    "os"
    "runtime"
    "sync"
    "time"
)

const version = "1.0.5r"

var kindConstructors = map[string] dlog.DlogConstructor {
    "amf": dlog.NewAmfDlog}

func main() {
    // cli options
    options := dlog.ParseFlags()
    if options.Version() {
        fmt.Println("dlogmon", version)
        os.Exit(0)
    }
    files := options.Files()
    totalFiles := len(files)

    // parallel level
    parallel := runtime.NumCPU()/2 + 1
    runtime.GOMAXPROCS(parallel)
    fmt.Printf("Parallel CPU: %d / %d\n", parallel, runtime.NumCPU())

    chLines := make(chan int, totalFiles)
    lock := new(sync.Mutex)

    // timing all the jobs up
    start := time.Now()

    // each dlog file is a goroutine
    var executor dlog.IDlogExecutor
    executors := make([]dlog.IDlogExecutor, totalFiles)
    for _, file := range files {
        executor = kindConstructors[options.Kind()](file, chLines, lock, options)
        executors = append(executors, executor)
        go executor.ScanLines(executor)
    }

    // wait for all dlog runner finish
    lines := 0
    for i:=0; i<totalFiles; i++ {
        lines += <- chLines
        executor.Progress(i)
    }

    end := time.Now()
    delta := end.Sub(start)
    fmt.Printf("\nParsed %d lines in %d files within %s [%.1f lines per second]\n",
        lines,
        totalFiles,
        delta, float64(lines)/delta.Seconds())
}

