package mr

import (
    "kx/progress"
)

var (
    KEY_SEP string
    OUTPUT_GROUP_HEADER_LEN int
)

func init() {
    sep := []byte{0x0}
    KEY_SEP = string(sep)

    OUTPUT_GROUP_HEADER_LEN = progress.TerminalWidth()
}
