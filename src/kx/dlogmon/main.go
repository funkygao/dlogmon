package main

import (
    "fmt"
    "kx/dlog"
    T "kx/trace"
    "kx/size"
    "kx/util"
    "log"
    "os"
    "runtime"
    "runtime/pprof"
    "time"
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
    kvResult := manager.WaitForCompletion()

    displaySummary(manager.Logger, start,
        manager.FilesCount(), manager.RawLines, manager.ValidLines)

    if option.Shell {
        cliCmdloop(manager.GetOneWorker(), kvResult)
    }
}

func displaySummary(logger *log.Logger, start time.Time, files, rawLines, validLines int) {
    defer T.Un(T.Trace(""))

    delta := time.Since(start)
    summary := fmt.Sprintf("Parsed %s/%s(%.4f%s) lines in %d files within %s [%.1f lines per second]\n",
        size.Comma(int64(validLines)),
        size.Comma(int64(rawLines)),
        100*float64(validLines)/float64(rawLines),
        "%%",
        files,
        delta,
        float64(rawLines)/delta.Seconds())
    // render to both log and stderr
    logger.Print(summary)
    fmt.Fprintf(os.Stderr, summary)
}

func initialize(option *dlog.Option, err error) {
    defer T.Un(T.Trace(""))

    if option.Version() {
        fmt.Fprintf(os.Stderr, "%s %s %s %s\n", "dlogmon", VERSION,
            runtime.GOOS, runtime.GOARCH)
        os.Exit(0)
    }

    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    // enable gc trace
    // this will not work, the only way is to setenv before invoke me
    os.Setenv("GOGCTRACE", "1")

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
