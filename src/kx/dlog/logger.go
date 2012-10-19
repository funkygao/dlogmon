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
        var err error
        logWriter, err = os.OpenFile(option.logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
        if err != nil {
            panic(err)
        }
    }

    return log.New(logWriter, "", log.Ldate|log.Lshortfile|log.Ltime|log.Lmicroseconds)
}
