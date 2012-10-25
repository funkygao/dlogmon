package dlog

import (
    "fmt"
    "kx/db"
    "kx/mr"
    T "kx/trace"
    "kx/progress"
    . "os"
    "os/signal"
    "runtime"
    "strings"
    "sync"
    "time"
)

func init() {
    db.Initialize(DbEngine, DbFile)
}

// Construct a TotalResult instance
func newTotalResult(rawLines, validLines int) TotalResult {
    return TotalResult{WorkerResult{rawLines, validLines}}
}

// Manager constructor
func NewManager(option *Option) *Manager {
    defer T.Un(T.Trace(""))

    this := new(Manager)
    if option.tick > 0 {
        this.ticker = time.NewTicker(time.Millisecond * time.Duration(option.tick))
    }
    this.Logger = newLogger(option)
    this.option = option
    this.lock = new(sync.Mutex)

    this.Println("manager created")

    return this
}

// Printable Manager
func (this *Manager) String() string {
    return fmt.Sprintf("Manager{%#v}", this.option)
}

// Get any worker of the same type TODO
func (this *Manager) getOneWorker() IWorker {
    defer T.Un(T.Trace(""))

    return this.workers[0]
}

// How many workers are running
func (this *Manager) workersCount() int {
    return this.FilesCount()
}

// How many lines added up
// For progress bar purpose
func (this Manager) totalLines() (total int) {
    for _, w := range this.workers {
        total += w.TotalLines()
    }
    return
}

// How many dlog files are being and to be analyzed
func (this *Manager) FilesCount() int {
    return len(this.option.files)
}

// Altogether how many raw lines parsed
func (this Manager) RawLines() int {
    return this.rawLines
}

// Get the global mutex object
func (this Manager) GetLock() *sync.Mutex {
    return this.lock
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

// Create all workers instances
func (this *Manager) newWorkers() {
    var worker IWorker
    this.workers = make([]IWorker, 0)
    for seq, file := range this.option.files {
        worker = workerConstructors[this.option.Kind()](this, this.option.Kind(), file, uint16(seq+1))
        this.workers = append(this.workers, worker)

        // type assertion
        if w, ok := worker.(IWorker); ok {
            if this.option.debug {
                fmt.Fprintf(Stderr, "worker type: %T\n", w)
            }
        }
    }

    this.Println("all worker instances created")
}

// Submit the job and start the job
func (this *Manager) Submit() (err error) {
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

    this.Println("submitted job accepted")

    // create channels first
    chMap, chWorker := make(chan interface{}, this.workersCount()), make(chan WorkerResult, this.workersCount())
    // the barrier
    this.chTotal = make(chan TotalResult)

    // create workers first
    this.newWorkers()

    // TODO
    go this.trapSignal()

    if this.ticker != nil {
        go this.runTicker()
    }

    if this.option.progress {
        this.chProgress = make(chan int, PROGRESS_CHAN_BUF)
        go this.showProgress()
    }

    // collect all workers output
    go this.collectWorkers(chMap, chWorker)

    // launch workers in chunk
    go this.launchWorkers(chMap, chWorker)

    return
}

func (this Manager) launchWorkers(chMap chan<- interface{}, chWorker chan<- WorkerResult) {
    this.Println("starting workers...")
    for i, worker := range this.workers {
        if this.option.Nworkers > 0 {
            running := i - this.doneWorkers // running workers count
            if running > this.option.Nworkers {
                runtime.Gosched()
            }
        }
        go worker.SafeRun(this.chProgress, chMap, chWorker)
    }
    this.Println("all workers started")
}

// Wait for all the dlog goroutines finish and collect final result
// Must run after collectWorkers() finished
func (this *Manager) WaitForCompletion() {
    defer T.Un(T.Trace(""))

    // 也可能我走的太快，得等他们先把chTotal创建好之后再开始
    for this.chTotal == nil {
        runtime.Gosched()
    }

    select {
    case r, ok := <-this.chTotal:
        if !ok {
            panic("chTotal unkown error")
        }

        this.rawLines, this.validLines = r.RawLines, r.ValidLines
    case <-time.After(time.Hour * 10):
        // timeout 10 hours? just demo useage of timeout
        break
    }

    close(this.chTotal)
    if this.chProgress != nil {
        close(this.chProgress)
    }

    this.Println("got workers summary, ready to finish")
}

// Collect worker's output
// including map data and worker summary
func (this *Manager) collectWorkers(chInMap chan interface{}, chInWorker chan WorkerResult) {
    defer T.Un(T.Trace(""))

    this.Println("collectWorkers started")

    transFromMapper := mr.NewTransformData()

    var (
        rawLines, validLines int
        doneWorkers int
        allDone int = 2 * this.workersCount()
    )
    for {
        if doneWorkers == allDone {
            break
        }

        select {
        case w, ok := <-chInWorker: // each worker send 1 msg to this chan
            if !ok {
                // this can never happens, worker can't close this chan
                this.Fatal("worker chan closed")
                break
            }
            doneWorkers ++
            this.doneWorkers ++
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
            doneWorkers ++
        }

        runtime.Gosched()
    }

    this.Println("all workers collected, next to merge and reduce...")

    // reduce the merged result
    // reduce cannot start until all the mappers have finished
    worker := this.getOneWorker()
    var r mr.ReduceResult = worker.Reduce(this.merge(worker.Name(), transFromMapper))
    this.exportToDb(worker.Name(), r)

    // all workers done, so close the channels
    this.Println("closing channels")
    close(chInMap)
    close(chInWorker)

    // WaitForCompletion will wait for this
    this.chTotal <- newTotalResult(rawLines, validLines)
}

func (this Manager) merge(name string, t mr.TransformData) (r mr.ReduceData) {
    defer T.Un(T.Trace(""))

    this.Printf("start to merge %s...\n", name)

    // init the ReduceData
    tagTypes := t.TagTypes()
    r = mr.NewReduceData(len(tagTypes))

    // trans -> reduce
    for k, v := range t {
        tagType, key := mr.GetTagType(k)
        r[tagType].AppendSlice(key, v)
    }

    return
}

func (this Manager) exportToDb(name string, r mr.ReduceResult) {
    defer T.Un(T.Trace(""))

    this.Printf("export %s reduce result to db\n", name)
    db.ImportResult(name, r)
}

func (this Manager) runTicker() {
    defer T.Un(T.Trace(""))

    for _ = range this.ticker.C {
        this.Println("mem:", T.MemAlloced(), "goroutines:", runtime.NumGoroutine())
    }
}

// Show progress bar
func (this Manager) showProgress() {
    total := this.totalLines()
    p := progress.New(total)

    var lines int
    for n := range this.chProgress {
        lines += n
        p.ShowProgress(lines)
    }
    println()
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
