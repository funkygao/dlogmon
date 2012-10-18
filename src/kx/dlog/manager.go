package dlog

import (
    "fmt"
    "log"
    "runtime"
    "sync"
    T "kx/trace"
)

// Manager(coordinator) of all the dlog goroutines
type Manager struct {
    executorsStarted bool
    rawLines, validLines int
    option *Option
    chFileScanResult chan ScanResult // each dlog goroutine will report to this
    chTotalScanResult chan ScanResult // total scan line collector use this to sync
    chLine chan Any
    lock *sync.Mutex
    *log.Logger
    executors []IDlogExecutor
}

var (
    constructors = map[string] DlogConstructor {
        "amf": NewAmfDlog}
)

// Manager constructor
func NewManager(option *Option) *Manager {
    defer T.Un(T.Trace("NewManager"))

    this := new(Manager)
    this.executorsStarted = false
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
    // wait to collect after all dlog executors done
    go this.collectLinesCount()

    // run each dlog in a goroutine
    var executor IDlogExecutor
    this.executors = make([]IDlogExecutor, 0)
    for _, file := range this.option.files {
        executor = constructors[this.option.Kind()](this, file)
        this.executors = append(this.executors, executor)
        go executor.Run(executor)
    }

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
    r := <- this.chTotalScanResult
    this.rawLines, this.validLines = r.RawLines, r.ValidLines
}

func (this *Manager) collectLinesCount() {
    var rawLines, validLines int
    //i := 0
    for {
        /*
        i ++
        if i > this.executorsCount() {
            break
        }*/
        select {
        case r := <- this.chFileScanResult:
            rawLines += r.RawLines
            validLines += r.ValidLines
        case r := <- this.chLine:
            fmt.Println("shit", r)
            /*
        default:

            if this.ExecutorsAllDone() {
                goto endloop
            }*/
        }

        runtime.Gosched()
    }

    //endloop:

    this.chTotalScanResult <- ScanResult{rawLines, validLines}
}

