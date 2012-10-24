package main

import (
    "fmt"
    "kx/dlog"
    T "kx/trace"
    "kx/util"
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
    // cli options
    option := dlog.ParseFlags()

    defer T.Un(T.Trace(""))

    setup(option)

    // construct the manager
    manager := dlog.NewManager(option)

    // cpu profile
    if option.Cpuprofile() != "" {
        defer pprof.StopCPUProfile()
    }

    // timing all the jobs up
    start := time.Now()

    go manager.SafeRun()

    // mem profile
    dumpMemProfile(option.Memprofile())
    manager.WaitForCompletion()

    displaySummary(start, manager)
}

func displaySummary(start time.Time, manager *dlog.Manager) {
    defer T.Un(T.Trace(""))

    end := time.Now()
    delta := end.Sub(start)
    manager.Printf("Parsed %d/%d lines in %d files within %s [%.1f lines per second]\n",
        manager.ValidLines(),
        manager.RawLines(),
        manager.FilesCount(),
        delta, float64(manager.RawLines())/delta.Seconds())
}

func setup(option *dlog.Option) {
    defer T.Un(T.Trace(""))

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

    // cpu profile
    if option.Cpuprofile() != "" {
        f, err := os.Create(option.Cpuprofile())
        if err != nil {
            panic(err)
        }

        pprof.StartCPUProfile(f)
    }
}

// mem profile
func dumpMemProfile(pf string) {
    defer T.Un(T.Trace(""))

    if pf != "" {
        f, err := os.Create(pf)
        if err != nil {
            panic(err)
        }

        pprof.WriteHeapProfile(f)
        f.Close()
    }
}

// assert cwd is right
func init() {
    if !util.FileExists(dlog.VarDir) {
        panic("must run on top dir")
    }
}
