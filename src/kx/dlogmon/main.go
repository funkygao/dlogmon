package main

import (
    "fmt"
    "kx/dlog"
    T "kx/trace"
    "os"
    "runtime"
    "time"
)

const VERSION = "1.0.6r"

const (
    maxprocsenv = "GOMAXPROCS"
)

func main() {
    defer T.Un(T.Trace("main"))

    // cli options
    option := dlog.ParseFlags()
    if option.Version() {
        fmt.Println("dlogmon", VERSION)
        os.Exit(0)
    }

    // parallel level
    if os.Getenv(maxprocsenv) == "" {
        parallel := runtime.NumCPU()/2 + 1
        runtime.GOMAXPROCS(parallel)
        fmt.Printf("Parallel CPU: %d / %d\n", parallel, runtime.NumCPU())
    }

    // timing all the jobs up
    start := time.Now()

    manager := dlog.NewManager(option)
    go manager.StartAll()
    manager.CollectAll()

    rawLines, validLines := manager.RawLines(), manager.ValidLines()

    end := time.Now()
    delta := end.Sub(start)
    fmt.Printf("\nParsed %d/%d lines in %d files within %s [%.1f lines per second]\n",
        validLines,
        rawLines,
        manager.FilesCount(),
        delta, float64(rawLines)/delta.Seconds())
}
