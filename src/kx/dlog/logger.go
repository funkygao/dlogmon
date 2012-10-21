package dlog

import (
    "io"
    "log"
    "os"
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

    if option.debug {
        return log.New(logWriter, LOG_PREFIX_DEBUG, LOG_OPTIONS_DEBUG)
    }

    return log.New(logWriter, LOG_PREFIX, LOG_OPTIONS)
}
