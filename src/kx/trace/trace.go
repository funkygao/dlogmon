package trace

import (
    "fmt"
    "runtime"
    "time"
)

// TODO how about any return values?
type AnyFunc func(args ...interface{})

var enabled bool

// Entering into a func
func Trace(fn string) string {
    if fn == "" {
        pc, _, _, _ := runtime.Caller(1)
        f := runtime.FuncForPC(pc)
        fn = f.Name() // the caller func name
    }
    if enabled {
        fmt.Println("Entering:", fn)
    }
    return fn
}

// Leaving from a func
func Un(fn string) {
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
