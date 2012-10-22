package dlog

import (
    "fmt"
    T "kx/trace"
    . "os"
    "os/signal"
    "runtime"
    "strings"
    "sync"
    "time"
)

// Manager constructor
func NewManager(option *Option) *Manager {
    defer T.Un(T.Trace(""))

    this := new(Manager)
    this.workersStarted = false
    if option.tick > 0 {
        this.ticker = time.NewTicker(time.Millisecond * time.Duration(option.tick))
    }
    this.Logger = newLogger(option)
    this.option = option
    this.lock = new(sync.Mutex)
    this.chFileScanResult, this.chTotalScanResult = make(chan ScanResult), make(chan ScanResult)
    this.chLine = make(chan Any, this.FilesCount())

    return this
}

// Printable Manager
func (this *Manager) String() string {
    return fmt.Sprintf("Manager{%#v}", this.option)
}

func (this *Manager) workersCount() int {
    return this.FilesCount()
}

// Are all dlog workers finished?
func (this *Manager) WorkersAllDone() bool {
    if !this.workersStarted {
        return false
    }

    for _, w := range this.workers {
        if w.Running() {
            return false
        }
    }

    return true
}

// How many dlog files are analyzed
func (this *Manager) FilesCount() int {
    return len(this.option.files)
}

// Altogether how many raw lines parsed
func (this Manager) RawLines() int {
    return this.rawLines
}

// Uock the global mutex
func (this Manager) Lock() {
    this.lock.Lock()
}

// Unlock the global mutex
func (this Manager) Unlock() {
    this.lock.Unlock()
}

// Altogether how many valid lines parsed
func (this Manager) ValidLines() int {
    return this.validLines
}

// Start and manage all the workers
func (this *Manager) StartAll() (err error) {
    // collection the panic's
    defer func() {
        if r := recover(); r != nil {
            var ok bool
            if err, ok = r.(error); !ok {
                err = fmt.Errorf("manager: %v", r)
            }
        }
    }()

    go this.trapSignal()

    if this.ticker != nil {
        go this.runTicker()
    }

    // wait to collect after all dlog workers done
    go this.collectLinesCount()

    this.Println("starting all workers...")

    // run each dlog in a goroutine
    var worker IWorker
    this.workers = make([]IWorker, 0)
    for _, file := range this.option.files {
        worker = workerConstructors[this.option.Kind()](this, file)
        this.workers = append(this.workers, worker)

        // type assertion
        if w, ok := worker.(IWorker); ok {
            if this.option.debug {
                fmt.Fprintf(Stderr, "worker type: %T\n", w)
            }

            go w.SafeRun(w)
        }
    }

    this.Println("all workers started.")
    this.workersStarted = true
    return
}

func (this *Manager) collectWorkerSummary(rawLines, validLines int) {
    this.chFileScanResult <- ScanResult{rawLines, validLines}
}

func (this *Manager) collectLineMeta(meta Any) {
    this.chLine <- meta
}

// Wait for all the dlog goroutines finish and collect final result
func (this *Manager) CollectAll() {
    r := <-this.chTotalScanResult
    this.rawLines, this.validLines = r.RawLines, r.ValidLines

    close(this.chFileScanResult)
    close(this.chLine)
}

func (this *Manager) collectLinesCount() {
    defer T.Un(T.Trace(""))

    this.Println("collector started")

    var rawLines, validLines int
    for {
        if this.WorkersAllDone() {
            break
        }

        select {
        case r, ok := <-this.chFileScanResult:
            if !ok {
                this.Println("chFileScanResult closed")
            }
            rawLines += r.RawLines
            validLines += r.ValidLines

        case r, ok := <-this.chLine:
            if !ok {
                this.Println("chLine closed")
            }
            fmt.Println("reducer: ", r)
        }

        runtime.Gosched()
    }

    this.chTotalScanResult <- ScanResult{rawLines, validLines}
}

func (this Manager) runTicker() {
    for _ = range this.ticker.C {
        this.Println("mem:", T.MemAlloced(), "goroutines:", runtime.NumGoroutine())
    }
}

func (this Manager) Shutdown() {
    Exit(0)
}

func (this Manager) trapSignal() {
    ch := make(chan Signal, 10)

    // register the given channel to receive notifications of the specified signals
    signal.Notify(ch, caredSignals...)

    go func() {
        sig := <-ch
        fmt.Fprintf(Stderr, "%s signal received...\n", strings.ToUpper(sig.String()))
        for _, skip := range skippedSignals {
            if skip == sig {
                this.Printf("%s signal ignored\n", strings.ToUpper(sig.String()))
                return
            }
        }

        // not skipped
        fmt.Fprintf(Stderr, "prepare to shutdown...")
        this.Shutdown()
    }()
}
