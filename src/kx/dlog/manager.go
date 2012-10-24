package dlog

import (
    "fmt"
    T "kx/trace"
    "kx/mr"
    "kx/db"
    . "os"
    "os/signal"
    "runtime"
    "strings"
    "sync"
    "time"
)

func init() {
    db.Initialize(DbEngine, DbFile, SQL_CREATE_TABLE)
}

// Construct a TotalResult instance
func newTotalResult(rawLines, validLines int) TotalResult {
    return TotalResult{WorkerResult{rawLines, validLines}}
}

// Manager constructor
func NewManager(option *Option) *Manager {
    defer T.Un(T.Trace(""))

    this := new(Manager)
    this.workersStarted = false
    if option.tick > 0 {
        this.ticker = time.NewTicker(time.Millisecond * time.Duration(option.tick))
    }
    this.Logger = newLogger(option)
    this.option = option
    this.lock = new(sync.Mutex)

    return this
}

// Printable Manager
func (this *Manager) String() string {
    return fmt.Sprintf("Manager{%#v}", this.option)
}

func (this *Manager) getOneWorker() IWorker {
    defer T.Un(T.Trace(""))

    return this.workers[0]
}

// How many workers are running
func (this *Manager) workersCount() int {
    return this.FilesCount()
}

// Are all dlog workers finished?
func (this *Manager) workersAllDone() bool {
    if !this.workersStarted {
        return false
    }

    for _, w := range this.workers {
        if w.Running() {
            return false
        }
    }

    return true
}

// How many dlog files are being and to be analyzed
func (this *Manager) FilesCount() int {
    return len(this.option.files)
}

// Altogether how many raw lines parsed
func (this Manager) RawLines() int {
    return this.rawLines
}

// Uock the global mutex
func (this Manager) Lock() {
    this.lock.Lock()
}

// Unlock the global mutex
func (this Manager) Unlock() {
    this.lock.Unlock()
}

// Altogether how many valid lines were parsed
func (this Manager) ValidLines() int {
    return this.validLines
}

// Start and manage all the workers safely
func (this *Manager) SafeRun() (err error) {
    defer T.Un(T.Trace(""))

    // safely: collection the panic's
    defer func() {
        if r := recover(); r != nil {
            var ok bool
            if err, ok = r.(error); !ok {
                err = fmt.Errorf("manager: %v", r)
            }
        }
    }()

    go this.trapSignal()

    if this.ticker != nil {
        go this.runTicker()
    }

    chMap, chWorker := make(chan interface{}, this.workersCount()), make(chan WorkerResult, this.workersCount())
    this.chTotal = make(chan TotalResult)

    // collect all workers output
    go this.collectWorkers(chMap, chWorker)

    this.Println("starting workers...")

    // run each dlog in a goroutine
    var worker IWorker
    this.workers = make([]IWorker, 0)
    for _, file := range this.option.files {
        worker = workerConstructors[this.option.Kind()](this, this.option.Kind(), file)
        this.workers = append(this.workers, worker)

        // type assertion
        if w, ok := worker.(IWorker); ok {
            if this.option.debug {
                fmt.Fprintf(Stderr, "worker type: %T\n", w)
            }

            go w.SafeRun(chMap, chWorker)
        }
    }

    this.Println("all workers started")
    this.workersStarted = true
    return
}

// Wait for all the dlog goroutines finish and collect final result
func (this *Manager) WaitForCompletion() {
    defer T.Un(T.Trace(""))

    select {
    case r := <-this.chTotal:
        this.Println("got the summary")
        this.rawLines, this.validLines = r.RawLines, r.ValidLines
    case <- time.After(time.Hour * 10):
        // timeout 10 hours? just demo useage of timeout
        break
    }

    close(this.chTotal)

    this.Println("manager ready to finish")
}

// Collect worker's output
// including map data and worker summary
func (this *Manager) collectWorkers(chInMap chan interface{}, chInWorker chan WorkerResult) {
    defer T.Un(T.Trace(""))

    this.Println(T.CallerFuncName(1), "started")

    transFromMapper := mr.NewTransformData()

    var rawLines, validLines int
    for {
        if this.workersAllDone() {
            break
        }

        select {
        case w, ok := <-chInWorker:
            if !ok {
                // this can never happens, worker can't close this chan
                this.Fatal("worker chan closed")
                break
            }
            rawLines += w.RawLines
            validLines += w.ValidLines

        case m, ok := <-chInMap:
            if !ok {
                // this can never happens, worker can't close this chan
                this.Fatal("line chan closed")
                break
            }
            for k, v := range m.(mr.TransformData) {
                transFromMapper.AppendSlice(k, v)
            }
        }

        runtime.Gosched()
    }

    // reduce the merged result
    worker := this.getOneWorker()
    var r mr.ReduceResult = worker.Reduce(this.merge(worker.Name(), transFromMapper))
    this.exportToDb(worker.Name(), r)

    // all workers done, so close the channels
    close(chInMap)
    close(chInWorker)

    this.chTotal <- newTotalResult(rawLines, validLines)

    this.Println(T.CallerFuncName(1), "all workers collected")
}

func (this Manager) merge(name string, t mr.TransformData) (r mr.ReduceData) {
    defer T.Un(T.Trace(""))

    this.Println(name, "merge")

    // init the ReduceData
    keyTypes := t.KeyTypes()
    r = mr.NewReduceData(len(keyTypes))

    // trans -> reduce
    for k, v := range t {
        keyType, key := mr.GetKeyType(k)
        r[keyType].AppendSlice(key, v)
    }

    return
}

func (this Manager) exportToDb(name string, r mr.ReduceResult) {
    defer T.Un(T.Trace(""))

    this.Println(name, "export reduce result to db")

    db.ImportResult(name, r)
}

func (this Manager) runTicker() {
    defer T.Un(T.Trace(""))

    for _ = range this.ticker.C {
        this.Println("mem:", T.MemAlloced(), "goroutines:", runtime.NumGoroutine())
    }
}

func (this Manager) Shutdown() {
    this.Println("shutdown now")
    Exit(0)
}

func (this Manager) trapSignal() {
    defer T.Un(T.Trace(""))

    ch := make(chan Signal, 10)

    // register the given channel to receive notifications of the specified signals
    signal.Notify(ch, caredSignals...)

    go func() {
        sig := <-ch
        fmt.Fprintf(Stderr, "%s signal received...\n", strings.ToUpper(sig.String()))
        for _, skip := range skippedSignals {
            if skip == sig {
                this.Printf("%s signal ignored\n", strings.ToUpper(sig.String()))
                return
            }
        }

        // not skipped
        fmt.Fprintf(Stderr, "prepare to shutdown...")
        this.Shutdown()
    }()
}
