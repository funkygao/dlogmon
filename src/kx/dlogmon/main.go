package main

import (
    "fmt"
    "kx/dlog"
    T "kx/trace"
    "os"
    "runtime"
    "runtime/pprof"
    "time"
)

const VERSION = "1.0.6r"

const (
    maxprocsenv = "GOMAXPROCS"
)

func main() {
    defer T.Un(T.Trace(""))

    // cli options
    option := dlog.ParseFlags()
    if option.Version() {
        fmt.Fprintf(os.Stderr, "%s %s\n", "dlogmon", VERSION)
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

    // cpu profile
    if option.Cpuprofile() != "" {
        f, err := os.Create(option.Cpuprofile())
        if err != nil {
            panic(err)
        }

        pprof.StartCPUProfile(f)

        defer pprof.StopCPUProfile()
    }

    manager := dlog.NewManager(option)
    go manager.StartAll()

    dumpMemProfile(option.Memprofile())

    manager.CollectAll()
    rawLines, validLines := manager.RawLines(), manager.ValidLines()

    end := time.Now()
    delta := end.Sub(start)
    manager.Printf("Parsed %d/%d lines in %d files within %s [%.1f lines per second]\n",
        validLines,
        rawLines,
        manager.FilesCount(),
        delta, float64(rawLines)/delta.Seconds())
}

// mem profile
func dumpMemProfile(pf string) {
    if pf != "" {
        f, err := os.Create(pf)
        if err != nil {
            panic(err)
        }

        pprof.WriteHeapProfile(f)
        f.Close()
    }
}
