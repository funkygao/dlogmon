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
    filename string
    chEof chan bool
    lock *sync.Mutex
}

func NewAmfDlog(filename string, ch chan bool, lock *sync.Mutex) *AmfDlog {
    dlog := new(AmfDlog)
    dlog.filename = filename
    dlog.chEof  = ch
    dlog.lock = lock

    return dlog
}

