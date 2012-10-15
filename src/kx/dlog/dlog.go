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
    debug bool
    filename string
    chLines chan int
    lock *sync.Mutex
}

func NewAmfDlog(filename string, ch chan int, lock *sync.Mutex, debug bool) *AmfDlog {
    dlog := new(AmfDlog)
    dlog.filename = filename
    dlog.chLines = ch
    dlog.lock = lock
    dlog.debug = debug

    return dlog
}

