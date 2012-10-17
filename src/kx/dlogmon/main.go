// main entry
package main

import (
    "fmt"
    "kx/dlog"
    "os"
    "runtime"
    "time"
)

const version = "1.0.5r"

func main() {
    // cli options
    options := dlog.ParseFlags()
    if options.Version() {
        fmt.Println("dlogmon", version)
        os.Exit(0)
    }

    // parallel level
    parallel := runtime.NumCPU()/2 + 1
    runtime.GOMAXPROCS(parallel)
    fmt.Printf("Parallel CPU: %d / %d\n", parallel, runtime.NumCPU())

    // timing all the jobs up
    start := time.Now()

    manager := dlog.NewManager(options)
    manager.StartAll()

    manager.CollectAll()

    totalLines, validLines := manager.TotalLines, manager.ValidLines

    end := time.Now()
    delta := end.Sub(start)
    fmt.Printf("\nParsed %d/%d lines in %d files within %s [%.1f lines per second]\n",
        validLines,
        totalLines,
        manager.FilesCount(),
        delta, float64(totalLines)/delta.Seconds())
}

