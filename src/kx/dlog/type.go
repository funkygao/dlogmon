package dlog

import (
    "bufio"
    "github.com/kless/goconfig/config"
    "log"
    "sync"
    "time"
    "kx/mr"
)

// dlog parser interface
type DlogParser interface {
    IsLineValid(string) bool
}

type Namer interface {
    Name() string
}

// Worker struct method signatures
type IWorker interface {
    Namer  // each kind of worker has a uniq name
    SafeRun(chan<- interface{}, chan<- WorkerResult)
    Running() bool
    Combiner() mr.CombinerFunc
    DlogParser
    mr.Mapper
    mr.Reducer
}

// For 1 dlog file worker
type Worker struct {
    name      string
    running   bool
    filename  string // dlog filename
    mapReader *bufio.Reader
    mapWriter *bufio.Writer
    *log.Logger
    manager *Manager
    combiner mr.CombinerFunc // can be nil
    executor IWorker // runtime dispatch
}

// AMF_SLOW tag analyzer
type AmfWorker struct {
    Worker
}

// Worker constructor signature
type WorkerConstructor func(*Manager, string, string) IWorker

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
