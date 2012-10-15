package dlog

const (
    LZOP_CMD = "lzop"
    LZOP_OPTION = "-dcf"
    EOL = '\n'
)

type DlogAware interface {
    ReadLines()
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

