package dlog

import "sync"

const (
    LZOP_CMD = "lzop"
    LZOP_OPTION = "-dcf"
    EOL = '\n'
)


// dlog interface
type DlogAware interface {
    ReadLines()
    IsLineValid(string) bool
}

// request interface
type RequestParser interface {
    ParseLine(string)
}

type Request struct {
    RequestParser
    http_method, uri, rid string
}

type Dlog struct {
    DlogAware
    options *Options
    filename string
    chLines chan int
    lock *sync.Mutex
}

func NewAmfDlog(filename string, ch chan int, lock *sync.Mutex, options *Options) *AmfDlog {
    dlog := new(AmfDlog)
    dlog.filename = filename
    dlog.chLines = ch
    dlog.lock = lock
    dlog.options = options

    return dlog
}

