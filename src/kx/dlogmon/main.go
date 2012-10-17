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

    chScan := make(chan dlog.ScanResult, totalFiles)
    lock := new(sync.Mutex)

    chScanedLines := make(chan dlog.ScanResult)
    go collectTotalLines(totalFiles, chScan, chScanedLines)

    // timing all the jobs up
    start := time.Now()

    // each dlog file is a goroutine
    var executor dlog.IDlogExecutor
    executors := make([]dlog.IDlogExecutor, totalFiles)
    for _, file := range files {
        executor = kindConstructors[options.Kind()](file, chScan, lock, options)
        executors = append(executors, executor)
        go executor.Run(executor)
    }

    // wait for all dlog goroutines done
    scannedLines := <- chScanedLines
    totalLines, validLines := scannedLines.TotalLines, scannedLines.ValidLines

    end := time.Now()
    delta := end.Sub(start)
    fmt.Printf("\nParsed %d/%d lines in %d files within %s [%.1f lines per second]\n",
        validLines,
        totalLines,
        totalFiles,
        delta, float64(totalLines)/delta.Seconds())
}

func collectTotalLines(n int, chScan chan dlog.ScanResult, chTotal chan dlog.ScanResult) (total, valid int) {
    for i:=0; i<n; i++ {
        s := <- chScan
        total += s.TotalLines
        valid += s.ValidLines
    }

    chTotal <- dlog.ScanResult{total, valid}
    return
}

