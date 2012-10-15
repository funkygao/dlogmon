// main entry
package main

import (
    "kx/dlog"
)

func main() {
    options := parseFlags()

    for _, file := range options.files {
        dlog := new(dlog.AmfDlog)
        dlog.SetFile(file)
        go dlog.ReadLines()
    }

}

