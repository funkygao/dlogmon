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

// TODO tag
type KeyType uint8

// Local aggregator
type CombinerFunc func([]float64) float64

// Mapper raw output format
type MapData map[string] float64

// mapper -> TransformData -> reduce
type TransformData map[string] []float64

type ReduceData []TransformData

type ReduceResult []MapData

// dlog parser interface
type DlogParser interface {
    IsLineValid(string) bool
}

// map
type Mapper interface {
    Map(string, chan<- Any)
}

// reduce
type Reducer interface {
    Reduce(ReduceData) ReduceResult
}

type Namer interface {
    Name() string
}

// Worker struct method signatures
type IWorker interface {
    Namer  // each kind of worker has a uniq name
    SafeRun(chan<- Any, chan<- WorkerResult)
    Running() bool
    Combiner() CombinerFunc
    DlogParser
    Mapper
    Reducer
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
    combiner CombinerFunc // can be nil
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
