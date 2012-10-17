/*
defer is LIFO

Usage:

    import t "kx/trace"

    func foo() {
        defer t.Un(t.Trace("foo"))

        //
    }
*/
package trace

import "fmt"

var enabled bool = true

// Entering into a func
func Trace(s string) string {
    if enabled {
        fmt.Println("Entering:", s)
    }
    return s
}

// Leaving from a func
func Un(s string) {
    if enabled {
        fmt.Println("Leaving:", s)
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

