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

// Construct a TotalResult instance
func newTotalResult(rawLines, validLines int) TotalResult {
    return TotalResult{WorkerResult{rawLines, validLines}}
}

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

    return this
}

// Printable Manager
func (this *Manager) String() string {
    return fmt.Sprintf("Manager{%#v}", this.option)
}

// How many workers are running
func (this *Manager) workersCount() int {
    return this.FilesCount()
}

// Are all dlog workers finished?
func (this *Manager) workersAllDone() bool {
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

// How many dlog files are being and to be analyzed
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

// Altogether how many valid lines were parsed
func (this Manager) ValidLines() int {
    return this.validLines
}

// Start and manage all the workers safely
func (this *Manager) SafeRun() (err error) {
    // safely: collection the panic's
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

    chLine, chWorker := make(chan Any, CH_LINES_BUFSIZE), make(chan WorkerResult, this.workersCount())
    this.chTotal = make(chan TotalResult)

    // collect all workers output
    go this.collectWorkers(chLine, chWorker, this.chTotal)

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

            go w.SafeRun(w, chLine, chWorker)
        }
    }

    this.Println("all workers started.")
    this.workersStarted = true
    return
}

// Wait for all the dlog goroutines finish and collect final result
func (this *Manager) Wait() {
    r := <-this.chTotal
    this.rawLines, this.validLines = r.RawLines, r.ValidLines

    //close(this.chFileScanResult)
    //close(this.chLine)
}

// Collect worker's output
// including line meta and worker summary
func (this *Manager) collectWorkers(chInLine <-chan Any, chInWorker <-chan WorkerResult, chOutTotal chan<- TotalResult) {
    defer T.Un(T.Trace(""))

    this.Println("collectWorkers started")

    var rawLines, validLines int
    for {
        if this.workersAllDone() {
            break
        }

        select {
        case w, ok := <-chInWorker:
            if !ok {
                this.Println("worker chan closed")
            }
            rawLines += w.RawLines
            validLines += w.ValidLines

        case l, ok := <-chInLine:
            if !ok {
                this.Println("line chan closed")
            }
            println(l.(string))
        }

        //runtime.Gosched()
    }

    chOutTotal <- newTotalResult(rawLines, validLines)
}

func (this Manager) runTicker() {
    for _ = range this.ticker.C {
        this.Println("mem:", T.MemAlloced(), "goroutines:", runtime.NumGoroutine())
    }
}

func (this Manager) Shutdown() {
    this.Println("shutdown now")
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
