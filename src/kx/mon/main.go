// main entry
package main

import (
    "kx/dlog"
    "runtime"
)

func main() {
    // parallel level
    runtime.GOMAXPROCS(runtime.NumCPU() + 1)

    // cli options
    options := parseFlags()

    chDlogDone := make(chan bool, len(options.files))

    for _, file := range options.files {
        dlog := dlog.NewAmfDlog(file, chDlogDone)
        go dlog.ReadLines()
    }

    // wait for all dlog runner finish
    for i:=0; i<len(options.files); i++ {
        <- chDlogDone
    }

}

