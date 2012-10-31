package main

import (
    "fmt"
    "kx/dlog"
    T "kx/trace"
    "kx/util"
    "log"
    "os"
    "runtime"
    "runtime/pprof"
    "time"
)

const VERSION = "1.0.6r"

const (
    maxprocsenv = "GOMAXPROCS"
)

// assert cwd is right
func init() {
    if !util.FileExists(dlog.VarDir) {
        panic("must run on top dir")
    }
}

func main() {
    // cli options
    option, err := dlog.ParseFlags()
    initialize(option, err)

    // construct the manager
    manager := dlog.NewManager(option)
    // mutex pass through
    T.SetLock(manager.GetLock())

    defer T.Un(T.Trace(""))

    // cpu profile
    if option.Cpuprofile() != "" {
        defer pprof.StopCPUProfile()
    }

    // timing all the jobs up
    start := time.Now()

    manager.Println("about to submit jobs")
    go manager.Submit()

    // mem profile
    dumpMemProfile(option.Memprofile())

    manager.Println("waiting for completion...")
    manager.WaitForCompletion()

    displaySummary(manager.Logger, start,
        manager.FilesCount(), manager.RawLines(), manager.ValidLines())
}

func displaySummary(logger *log.Logger, start time.Time, files, rawLines, validLines int) {
    defer T.Un(T.Trace(""))

    end := time.Now()
    delta := end.Sub(start)
    summary := fmt.Sprintf("Parsed %d/%d lines in %d files within %s [%.1f lines per second]\n",
        validLines,
        rawLines,
        files,
        delta, float64(rawLines)/delta.Seconds())
    logger.Print(summary)
    fmt.Fprintf(os.Stderr, summary)
}

func initialize(option *dlog.Option, err error) {
    defer T.Un(T.Trace(""))

    if option.Version() {
        fmt.Fprintf(os.Stderr, "%s %s\n", "dlogmon", VERSION)
        os.Exit(0)
    }

    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    // parallel level
    if os.Getenv(maxprocsenv) == "" {
        parallel := runtime.NumCPU()/2 + 1
        runtime.GOMAXPROCS(parallel)
        fmt.Fprintf(os.Stderr, "Parallel CPU(core): %d / %d, Concurrent workers: %d\n", parallel,
            runtime.NumCPU(), option.Nworkers)
    }
    fmt.Fprintln(os.Stderr, option.Timespan)

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
