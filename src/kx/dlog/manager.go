package dlog

import (
    "fmt"
    "log"
    "sync"
)

// Manager(coordinator) of all the dlog goroutines
type Manager struct {
    totalLines, validLines int
    options *Options
    ChFileScanResult chan ScanResult // each dlog goroutine will report to this
    chTotalScanResult chan ScanResult // total scan line collector use this to sync
    lock *sync.Mutex
    *log.Logger
    executors []IDlogExecutor
}

var (
    constructors = map[string] DlogConstructor {
        "amf": NewAmfDlog}
)

// Manager constructor
func NewManager(options *Options) *Manager {
    this := new(Manager)
    this.options = options
    this.lock = new(sync.Mutex)
    this.ChFileScanResult, this.chTotalScanResult = make(chan ScanResult), make(chan ScanResult)

    return this
}

// Printable Manager
func (this *Manager) String() string {
    return fmt.Sprintf("Manager{%#v}", this.options)
}

func (this *Manager) executorsCount() int {
    return this.FilesCount()
}

// Are all dlog executors finished?
func (this *Manager) DlogsDone() bool {
    for _, dlog := range this.executors {
        if dlog.Running() {
            return false
        }
    }

    return true
}

// How many dlog files are analyzed
func (this *Manager) FilesCount() int {
    return len(this.options.files)
}

// Altogether how many raw lines parsed
func (this Manager) TotalLines() int {
    return this.totalLines
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
    go this.collectTotalLines()

    // run each dlog in a goroutine
    var executor IDlogExecutor
    this.executors = make([]IDlogExecutor, 0)
    for _, file := range this.options.files {
        executor = constructors[this.options.Kind()](this, file)
        this.executors = append(this.executors, executor)
        go executor.Run(executor)
    }
}

// Wait for all the dlog goroutines finish and collect final result
func (this *Manager) CollectAll() {
    r := <- this.chTotalScanResult
    this.totalLines, this.validLines = r.TotalLines, r.ValidLines
}

func (this *Manager) collectTotalLines() {
    var total, valid int
    for i:=0; i<this.executorsCount(); i++ {
        r := <- this.ChFileScanResult

        total += r.TotalLines
        valid += r.ValidLines
    }

    this.chTotalScanResult <- ScanResult{total, valid}
}

