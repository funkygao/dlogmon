package dlog

import (
    "fmt"
    T "kx/trace"
    "log"
    "runtime"
    "sync"
    "time"
)

// Manager(coordinator) of all the dlog goroutines
type Manager struct {
    executorsStarted     bool
    rawLines, validLines int
    option               *Option
    chFileScanResult     chan ScanResult // each dlog goroutine will report to this
    chTotalScanResult    chan ScanResult // total scan line collector use this to sync
    chLine               chan Any
    lock                 *sync.Mutex
    ticker *time.Ticker
    *log.Logger
    executors []IDlogExecutor
}

var (
    constructors = map[string]DlogConstructor{
        "amf": NewAmfDlog}
)

// Manager constructor
func NewManager(option *Option) *Manager {
    defer T.Un(T.Trace("NewManager"))

    this := new(Manager)
    this.executorsStarted = false
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

func (this *Manager) executorsCount() int {
    return this.FilesCount()
}

// Are all dlog executors finished?
func (this *Manager) ExecutorsAllDone() bool {
    if !this.executorsStarted {
        return false
    }

    for _, dlog := range this.executors {
        if dlog.Running() {
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

// Global mutex
func (this Manager) Lock() *sync.Mutex {
    return this.lock
}

// Altogether how many valid lines parsed
func (this Manager) ValidLines() int {
    return this.validLines
}

// Start and manage all the dlog executors
func (this *Manager) StartAll() {
    if this.ticker != nil {
        go this.runTicker()
    }

    // wait to collect after all dlog executors done
    go this.collectLinesCount()

    this.Println("starting all executors...")

    // run each dlog in a goroutine
    var executor IDlogExecutor
    this.executors = make([]IDlogExecutor, 0)
    for _, file := range this.option.files {
        executor = constructors[this.option.Kind()](this, file)
        this.executors = append(this.executors, executor)
        go executor.Run(executor)
    }

    this.Println("all executors started.")
    this.executorsStarted = true
}

func (this *Manager) collectExecutorSummary(rawLines, validLines int) {
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
    defer T.Un(T.Trace("collectExecutors"))

    this.Println("collector started")

    var rawLines, validLines int
    for {
        if this.ExecutorsAllDone() {
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
