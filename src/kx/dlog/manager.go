package dlog

import (
    "fmt"
    T "kx/trace"
    "log"
    . "os"
    "os/signal"
    "runtime"
    "syscall"
    "sync"
    "time"
)

// Manager(coordinator) of all the dlog goroutines
type Manager struct {
    workersStarted     bool // all workers started?
    rawLines, validLines int
    option               *Option
    chFileScanResult     chan ScanResult // each dlog goroutine will report to this
    chTotalScanResult    chan ScanResult // total scan line collector use this to sync
    chLine               chan Any
    lock                 *sync.Mutex
    ticker *time.Ticker
    *log.Logger
    workers []IWorker
}

var (
    constructors = map[string]WorkerConstructor{
        "amf": NewAmfWorker}
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
        worker = constructors[this.option.Kind()](this, file)
        this.workers = append(this.workers, worker)

        // type assertion
        if e, ok := worker.(IWorker); ok {
            if this.option.debug {
                fmt.Fprintf(Stderr, "worker type: %T\n", e)
            }

            go e.SafeRun(e)
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
        this.Println("mem:", T.MemAlloced())
    }
}

func (this Manager) Shutdown() {
    Exit(0)
}

func (this Manager) trapSignal() {
    sch := make(chan Signal, 10)
    signal.Notify(sch, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT,
        syscall.SIGHUP, syscall.SIGSTOP, syscall.SIGQUIT)
    go func(ch <- chan Signal) {
        sig := <- ch
        fmt.Fprintln(Stderr, "signal received...", sig)
        this.Shutdown()
    }(sch)
}
