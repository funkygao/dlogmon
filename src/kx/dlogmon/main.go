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

    chTotalLines := make(chan int)
    go collectTotalLines(totalFiles, chLines, chTotalLines)

    // timing all the jobs up
    start := time.Now()

    // each dlog file is a goroutine
    var executor dlog.IDlogExecutor
    executors := make([]dlog.IDlogExecutor, totalFiles)
    for _, file := range files {
        executor = kindConstructors[options.Kind()](file, chLines, lock, options)
        executors = append(executors, executor)
        go executor.Run(executor)
    }

    // wait for all dlog goroutines done
    totalLines := <- chTotalLines

    end := time.Now()
    delta := end.Sub(start)
    fmt.Printf("\nParsed %d lines in %d files within %s [%.1f lines per second]\n",
        totalLines,
        totalFiles,
        delta, float64(totalLines)/delta.Seconds())
}

func collectTotalLines(n int, ch chan int, chT chan int) (total int) {
    for i:=0; i<n; i++ {
        total += <- ch
    }

    chT <- total
    return
}

