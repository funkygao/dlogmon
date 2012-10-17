/*
Scan each specified dlog files concurrently and merge parsed meta
info from each dlog file and render the final calculated report
*/
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

    rawLines, validLines := manager.RawLines(), manager.ValidLines()

    end := time.Now()
    delta := end.Sub(start)
    fmt.Printf("\nParsed %d/%d lines in %d files within %s [%.1f lines per second]\n",
        validLines,
        rawLines,
        manager.FilesCount(),
        delta, float64(rawLines)/delta.Seconds())
}

