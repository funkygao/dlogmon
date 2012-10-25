package dlog

import (
    "bufio"
    "github.com/kless/goconfig/config"
    "kx/mr"
    "log"
    "sync"
    "time"
)

// dlog parser interface
type DlogParser interface {
    IsLineValid(string) bool
}

type Namer interface {
    Name() string
}

type LineCounter interface {
    TotalLines() int
}

// Worker struct method signatures
type IWorker interface {
    Namer // each kind of worker has a uniq name
    SafeRun(chan<- int, chan<- interface{}, chan<- WorkerResult)
    Combiner() mr.CombinerFunc
    LineCounter
    DlogParser
    mr.Mapper
    mr.Reducer
}

// For 1 dlog file worker
type Worker struct {
    name      string
    filename  string // dlog filename
    mapReader *bufio.Reader
    mapWriter *bufio.Writer
    *log.Logger
    manager  *Manager
    combiner mr.CombinerFunc // can be nil
    self IWorker         // runtime dispatch
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
    rawLines, validLines int
    option               *Option
    lock                 *sync.Mutex
    ticker               *time.Ticker
    *log.Logger
    workers []IWorker
    chTotal chan TotalResult
    chProgress chan int // default <nil>
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
    progress               bool
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
