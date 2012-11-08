package dlog

import (
    "bufio"
    "github.com/kless/goconfig/config"
    "kx/mr"
    "log"
    "sync"
    "time"
)

type StreamResult string

type Kinder interface {
    Kind() string
}

type LineCounter interface {
    TotalLines() int
}

type TopNer interface {
    TopN() int
}

// Worker struct method signatures
type IWorker interface {
    Kinder // each kind of worker has a uniq name
    SafeRun(chan<- int, chan<- mr.KeyValue, chan<- Worker)
    Combiner() mr.CombinerFunc
    LineCounter
    mr.Mapper
    mr.Reducer
    mr.Filter
    TopNer
    mr.Printer
}

// For 1 dlog file worker
// Abstract
type Worker struct {
    kind                      string
    seq                       uint16 // sequence number
    filename                  string // dlog filename
    CreatedAt, StartAt, EndAt time.Time
    RawLines, ValidLines      int
    mapReader                 *bufio.Reader
    mapWriter                 *bufio.Writer
    *log.Logger
    manager  *Manager
    combiner mr.CombinerFunc // can be nil
    self     IWorker         // runtime dispatch
}

// Workers
type (
    KxiWorker struct {
        Worker
    }

    AmfWorker struct {
        Worker
    }

    NoopWorker struct {
        Worker
    }

    FileWorker struct {
        Worker
    }
)

// Worker constructor signature
type WorkerConstructor func(*Manager, string, string, uint16) IWorker

// Manager(coordinator) of all the dlog goroutines
type Manager struct {
    RawLines, ValidLines int
    option               *Option
    lock                 *sync.Mutex
    ticker               *time.Ticker
    *log.Logger
    workers       []IWorker
    chWorkersDone chan mr.KeyValue
    chProgress    chan int // default <nil>
}

// CLI options object
type Option struct {
    files                  []string
    Timespan               string
    debug                  bool
    trace                  bool
    verbose                bool
    version                bool
    filemode               bool
    Shell                  bool
    rpc                    bool
    progress               bool
    Nworkers               uint8 // how many concurrent workers(goroutines) permitted
    tick                   int   // in ms
    cpuprofile, memprofile string
    mapper                 string // mapper stream exe filename
    reducer                string // reducer stream exe filename
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
