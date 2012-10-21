package dlog

import (
    "bufio"
    "github.com/kless/goconfig/config"
    "log"
    "sync"
    "time"
)

// Any kind of things
type Any interface{}

// Worker struct method signatures
type IWorker interface {
    SafeRun(IWorker) // IWorker param for dynamic polymorphism
    IsLineValid(string) bool
    ExtractLineInfo(string) Any
    Running() bool
}

// For 1 dlog file worker
type Worker struct {
    running   bool
    filename  string // dlog filename
    mapReader *bufio.Reader
    mapWriter *bufio.Writer
    *log.Logger
    manager *Manager
}

// AMF_SLOW tag analyzer
type AmfWorker struct {
    Worker
}

// Worker constructor signature
type WorkerConstructor func(*Manager, string) IWorker

// Request object for a line
type Request struct {
    http_method, uri, rid string
}

// a single line meta info
type amfRequest struct {
    Request
    class, method, args string
    time                int16
}

// Scan result of raw lines and valid lines
type ScanResult struct {
    RawLines, ValidLines int
}

// Manager(coordinator) of all the dlog goroutines
type Manager struct {
    workersStarted       bool // all workers started?
    rawLines, validLines int
    option               *Option
    chFileScanResult     chan ScanResult // each dlog goroutine will report to this
    chTotalScanResult    chan ScanResult // total scan line collector use this to sync
    chLine               chan Any
    lock                 *sync.Mutex
    ticker               *time.Ticker
    *log.Logger
    workers []IWorker
}

// CLI options object
type Option struct {
    files                  []string
    debug                  bool
    trace                  bool
    verbose                bool
    version                bool
    Nworkers               int // how many concurrent workers(goroutines) permitted
    tick                   int // in ms
    cpuprofile, memprofile string
    mapper                 string
    reducer                string
    kind                   string
    conf                   *config.Config
}
