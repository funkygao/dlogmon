// main entry
package main

import (
    "kx/dlog"
)

func main() {
    options := parseFlags()

    chDlogDone := make(chan bool, 10)

    for _, file := range options.files {
        dlog := new(dlog.AmfDlog)
        dlog.SetFile(file)
        go dlog.ReadLines(chDlogDone)
    }

    // wait for all dlog runner finish
    for i:=0; i<len(options.files); i++ {
        <- chDlogDone
    }

}

