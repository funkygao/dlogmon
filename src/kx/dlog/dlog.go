package dlog

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
}

func NewAmfDlog(filename string, ch chan bool) *AmfDlog {
    dlog := new(AmfDlog)
    dlog.filename = filename
    dlog.chEof  = ch

    return dlog
}

