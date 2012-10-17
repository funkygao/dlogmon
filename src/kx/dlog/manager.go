package dlog

import (
    "sync"
)

// Manager(coordinator) of all the dlog goroutines
type Manager struct {
    TotalLines, ValidLines int
    options *Options
    chFileScanResult chan ScanResult // each dlog goroutine will report to this
    chTotalScanResult chan ScanResult // total scan line collector use this to sync
    lock *sync.Mutex
    executors []IDlogExecutor
}

var (
    constructors = map[string] DlogConstructor {
        "amf": NewAmfDlog}
)

func NewManager(options *Options) *Manager {
    this := new(Manager)
    this.options = options
    this.lock = new(sync.Mutex)
    this.chFileScanResult, this.chTotalScanResult = make(chan ScanResult), make(chan ScanResult)

    return this
}

func (this *Manager) executorsCount() int {
    return this.FilesCount()
}

func (this *Manager) FilesCount() int {
    return len(this.options.files)
}

// Start and manage all the dlog executors
func (this *Manager) StartAll() {
    // wait to collect after all dlog executors done
    go this.collectTotalLines()

    // run each dlog in a goroutine
    var executor IDlogExecutor
    this.executors = make([]IDlogExecutor, 0)
    for _, file := range this.options.files {
        executor = constructors[this.options.Kind()](file, this.chFileScanResult, 
            this.lock, this.options)
        this.executors = append(this.executors, executor)
        go executor.Run(executor)
    }
}

// Wait for all the dlog goroutines finish and collect final result
func (this *Manager) CollectAll() {
    r := <- this.chTotalScanResult
    this.TotalLines, this.ValidLines = r.TotalLines, r.ValidLines
}

func (this *Manager) collectTotalLines() {
    var total, valid int
    for i:=0; i<this.executorsCount(); i++ {
        r := <- this.chFileScanResult
        total += r.TotalLines
        valid += r.ValidLines
    }

    this.chTotalScanResult <- ScanResult{total, valid}
}

