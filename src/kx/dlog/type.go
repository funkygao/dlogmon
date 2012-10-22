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

type KeyType uint8

// Mapper output format
type MapOut map[string] int

// Reducer input format
type ReduceIn map[string] []int

// dlog parser interface
type DlogParser interface {
    IsLineValid(string) bool
    ExtractLineInfo(string) Any
}

// Worker struct method signatures
type IWorker interface {
    SafeRun(IWorker, chan<- Any, chan<- WorkerResult) // IWorker param for dynamic polymorphism
    Running() bool
    DlogParser
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

// Result of a worker
type WorkerResult struct {
    RawLines, ValidLines int
}

// Result of all workers
type TotalResult struct {
    WorkerResult
}

// Manager(coordinator) of all the dlog goroutines
type Manager struct {
    workersStarted       bool // all workers started?
    rawLines, validLines int
    option               *Option
    lock                 *sync.Mutex
    ticker               *time.Ticker
    *log.Logger
    workers []IWorker
    chTotal chan TotalResult
}

// map -> sort -> merge -> reduce
type Sorter interface {
}

// map -> sort -> merge -> reduce
type Merger interface {
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
