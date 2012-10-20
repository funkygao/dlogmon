package dlog

import (
    "io"
    "log"
    "os"
)

const (
   LOG_OPTIONS_DEBUG = log.Ldate | log.Lshortfile | log.Ltime | log.Lmicroseconds
   LOG_OPTIONS = log.LstdFlags
   LOG_PREFIX_DEBUG = "debug "
   LOG_PREFIX = ""
)

func newLogger(option *Option) *log.Logger {
    var logWriter io.Writer = os.Stderr
    if option.conf != nil {
        logfile, err := option.conf.String("default", "logfile")
        if err == nil {
            logWriter, err = os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
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
