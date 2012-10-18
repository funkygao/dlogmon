package dlog

import (
    "io"
    "log"
    "os"
)

// FIXME, log file has no output
func newLogger(option *Option) *log.Logger {
    var logWriter io.Writer
    if option.logfile == "" {
        logWriter = os.Stderr
    } else {
        f, e := os.OpenFile(option.logfile, os.O_APPEND | os.O_CREATE, 0666)
        if e != nil {
            panic(e)
        }
        logWriter = f
    }

    return log.New(logWriter, "",  log.Ldate | log.Lshortfile | log.Ltime | log.Lmicroseconds)
}

