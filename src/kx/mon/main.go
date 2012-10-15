// main entry
package main

import (
    "kx/dlog"
)

func main() {
    options := parseFlags()

    chDlogDone := make(chan bool, 10)

    for _, file := range options.files {
        dlog := dlog.NewAmfDlog(file, chDlogDone)
        go dlog.ReadLines()
    }

    // wait for all dlog runner finish
    for i:=0; i<len(options.files); i++ {
        <- chDlogDone
    }

}

