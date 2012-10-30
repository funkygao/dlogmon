package dlog

import (
    "fmt"
    "io"
    "log"
    "os"
    "syscall"
)

func newLogger(option *Option) *log.Logger {
    var logWriter io.Writer = os.Stderr
    if option.conf != nil {
        logfile, err := option.conf.String("default", "logfile")
        if err == nil {
            logWriter, err = os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
            if err != nil {
                panic(err)
            }
        }
    }

    prefix := fmt.Sprintf("[%d] ", syscall.Getpid())
    if option.debug {
        return log.New(logWriter, prefix, LOG_OPTIONS_DEBUG)
    }

    return log.New(logWriter, prefix, LOG_OPTIONS)
}
