package trace

import (
    "fmt"
    "kx/size"
    "runtime"
    "strings"
    "sync/atomic"
    "sync"
    "time"
)

// TODO how about any return values?
type AnyFunc func(args ...interface{})

const (
    DEPTH_CHAR = " "
    DEPTH_STEP = 4
)

var (
    enabled bool
    depth int32
    lock *sync.Mutex
)

// Entering into a func
func Trace(fn string) string {
    if !enabled {
        return fn
    }

    if fn == "" {
        fn = CallerFuncName(2)
    }

    space := strings.Repeat(DEPTH_CHAR, int(atomic.LoadInt32(&depth)))
    if lock != nil {
        lock.Lock()
        defer lock.Unlock()
    }
    fmt.Printf("%s %s %s\n", space, "Entering:", fn)

    atomic.AddInt32(&depth, DEPTH_STEP)

    return fn
}

// Leaving from a func
func Un(fn string) {
    if !enabled {
        return
    }

    if fn == "" {
        fn = CallerFuncName(2)
    }

    atomic.AddInt32(&depth, -DEPTH_STEP)

    space := strings.Repeat(DEPTH_CHAR, int(atomic.LoadInt32(&depth)))
    if lock != nil {
        lock.Lock()
        defer lock.Unlock()
    }
    fmt.Printf("%s %s %s\n", space, "Leaving :", fn)
}

// Enable the trace output
func Enable() {
    enabled = true
}

// Disable the trace output
func Disable() {
    enabled = false
}

// Setup the global mutex
func SetLock(l *sync.Mutex) {
    lock = l
}

// Measure how long it takes to run a func
func Timeit(f AnyFunc, args ...interface{}) time.Duration {
    start := time.Now()

    f(args...) // call the func

    end := time.Now()
    delta := end.Sub(start)
    return delta
}

// How much mem allocated
func MemAlloced() size.ByteSize {
    ms := &runtime.MemStats{}
    runtime.ReadMemStats(ms)
    return size.ByteSize(ms.Alloc)
}

// Caller func name with skip as the call stack level
func CallerFuncName(calldepth int) string {
    pc, _, _, _ := runtime.Caller(calldepth)
    f := runtime.FuncForPC(pc)
    return f.Name() // the caller func name
}
