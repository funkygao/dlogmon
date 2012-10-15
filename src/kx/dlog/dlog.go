package dlog

type DlogAware interface {
    SetFile(string)
    ReadLines()
    IsLineValid(string) bool
    OperateLine(string)
}

type Dlog struct {
    DlogAware
    filename string
}

const (
    LZOP_CMD = "lzop"
    LZOP_OPTION = "-dcf"
    EOL = '\n'
)
