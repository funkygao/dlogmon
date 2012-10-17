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

func Trace(s string) string {
    fmt.Println("Entering:", s)
    return s
}

func Un(s string) {
    fmt.Println("Leaving:", s)
}

