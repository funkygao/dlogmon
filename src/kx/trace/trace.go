package trace

import (
    "fmt"
    "runtime"
    "kx/size"
    "time"
)

// TODO how about any return values?
type AnyFunc func(args ...interface{})

var enabled bool

// Caller func name with skip as the call stack level
func CallerFuncName(skip int) string {
    pc, _, _, _ := runtime.Caller(skip)
    f := runtime.FuncForPC(pc)
    return f.Name() // the caller func name
}

// Entering into a func
func Trace(fn string) string {
    if fn == "" {
        fn = CallerFuncName(2)
    }
    if enabled {
        fmt.Println("Entering:", fn)
    }
    return fn
}

// Leaving from a func
func Un(fn string) {
    if fn == "" {
        fn = CallerFuncName(2)
    }
    if enabled {
        fmt.Println("Leaving:", fn)
    }
}

// Enable the trace output
func Enable() {
    enabled = true
}

// Disable the trace output
func Disable() {
    enabled = false
}

// Measure how long it takes to run a func
func Timeit(f AnyFunc, args ...interface{}) time.Duration {
    start := time.Now()

    f(args...) // call the func

    end := time.Now()
    delta := end.Sub(start)
    return delta
}

func MemAlloced() size.ByteSize {
    ms := &runtime.MemStats{}
    runtime.ReadMemStats(ms)
    return size.ByteSize(ms.TotalAlloc)
}
