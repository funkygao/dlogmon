// Altanative name: controller
package dlog

import (
    "fmt"
    "github.com/kless/goconfig/config"
    "kx/db"
    "kx/mr"
    "kx/progress"
    T "kx/trace"
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

func (this Manager) Conf() *config.Config {
    return this.option.conf
}

func (this *Manager) String() string {
    return fmt.Sprintf("Manager{%+v}", this.option)
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

// How many dlog files are being and to be analyzed
func (this *Manager) FilesCount() int {
    return len(this.option.files)
}

// How many lines added up
// For progress bar purpose
func (this Manager) totalLines() (total int) {
    for _, w := range this.workers {
        total += w.TotalLines()
    }
    return
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

// Create all workers instances
func (this *Manager) newWorkers() {
    var worker IWorker
    this.workers = make([]IWorker, 0)
    for seq, file := range this.option.files {
        worker = workerConstructors[this.option.Kind()](this, this.option.Kind(), file, uint16(seq+1))
        this.workers = append(this.workers, worker)
    }

    this.Println("all worker instances created")
}

func (this *Manager) initRateLimit() chan bool {
    if this.option.Nworkers == 0 || int(this.option.Nworkers) > this.FilesCount() {
        this.option.Nworkers = uint8(this.FilesCount())
    }

    chRateLimit := make(chan bool, this.option.Nworkers)

    // first let it start Nworkers without being blocked
    for i := uint8(0); i < this.option.Nworkers; i++ {
        chRateLimit <- true
    }

    return chRateLimit
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

    chMap := make(chan mr.KeyValue, this.workersCount()*LINE_CHANBUF_PER_WORKER)
    chWorker := make(chan Worker, this.workersCount())
    this.chWorkersDone = make(chan bool)

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
    chRateLimit := this.initRateLimit()
    go this.collectWorkers(chRateLimit, chMap, chWorker)

    // launch workers in chunk
    go this.launchWorkers(chRateLimit, chMap, chWorker)

    return
}

func (this Manager) launchWorkers(chRateLimit <-chan bool, chMap chan<- mr.KeyValue, chWorker chan<- Worker) {
    this.Println("starting workers...")

    for seq := 0; seq < this.workersCount(); seq++ {
        <-chRateLimit // 放闸

        worker := this.workers[seq]
        go worker.SafeRun(this.chProgress, chMap, chWorker)
    }

    this.Println("all workers started")
}

// Wait for all the dlog goroutines finish and collect final result
// Must run after collectWorkers() finished
func (this *Manager) WaitForCompletion() {
    defer T.Un(T.Trace(""))

    // 也可能我走的太快，得等他们先创建好再开始
    for this.chWorkersDone == nil {
        runtime.Gosched()
    }

    select {
    case _, ok := <-this.chWorkersDone:
        if !ok {
            panic("unkown error")
        }
    case <-time.After(time.Hour):
        // timeout 1 hour? just demo useage of timeout
        break
    }

    close(this.chWorkersDone)
    if this.chProgress != nil {
        close(this.chProgress)
    }

    // stop the ticker
    if this.ticker != nil {
        this.ticker.Stop()
    }

    this.Println("got workers summary, ready to finish")
}

// Collect worker's output
// including map data and worker summary
func (this *Manager) collectWorkers(chRateLimit chan bool, chInMap chan mr.KeyValue, chInWorker chan Worker) {
    defer T.Un(T.Trace(""))

    this.Println("collectWorkers started")

    shuffledKvs := make(chan mr.KeyValues)
    go mr.Shuffle(chInMap, shuffledKvs)

    const MSG_PER_WORKER = 2

    var doneWorkers int
    for {
        if doneWorkers == this.workersCount() {
            break
        }

        select {
        case w, ok := <-chInWorker: // each worker send 1 msg to this chan
            if !ok {
                // this can never happens, worker can't close this chan
                this.Fatal("worker chan closed")
                break
            }

            doneWorkers++
            this.Printf("workers done: %d/%d %.1f%%\n", doneWorkers,
                this.workersCount(), float64(100*doneWorkers/this.workersCount()))

            this.RawLines += w.RawLines
            this.ValidLines += w.ValidLines

            chRateLimit <- true // 让贤
        }

        //runtime.Gosched()
    }

    this.Println("all workers Map done, start to Shuffle...")

    // all workers done, so close the channels
    close(chInMap)
    close(chInWorker)
    close(chRateLimit)

    this.invokeGc()

    // mappers must complete before reducers can begin
    worker := this.getOneWorker()
    this.Println(worker.Kind(), "worker start to reduce with ordered keys...")
    reduceResult := (<-shuffledKvs).LaunchReducer(worker)
    this.Println(worker.Kind(), "worker reduce done")

    this.invokeGc()

    // enter into output phase
    // export final result, possibly export to db
    this.Println(worker.Kind(), "worker start to export result...")
    reduceResult.ExportResult(worker, worker.TopN())

    // WaitForCompletion will wait for this
    this.chWorkersDone <- true
}

func (this Manager) invokeGc() {
    start := time.Now()
    alloced := T.MemAlloced()
    runtime.GC()
    this.Println("GC", time.Now().Sub(start), alloced, "->", T.MemAlloced())
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
