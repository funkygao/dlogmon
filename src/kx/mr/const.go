package mr

import (
    "os"
)

const (
    SORT_BY_KEY SortType = iota + 1
    SORT_BY_VALUE
    SORT_BY_COL

    SORT_ORDER_ASC  = 1
    SORT_ORDER_DESC = 2

    KEY_SECONDARY_KV = "skv"
    KEY_SECONDARY_VK = "svk"
)

const (
    GOB_FILE_FLAG = os.O_CREATE | os.O_WRONLY
    GOB_FILE_PERM = 0600
)

const (
    OUTPUT_GROUP_HEADER_LEN = 100
    OUTPUT_VAL_WIDTH = 7
)

var KEY_SEP string
func init() {
    sep := []byte{0x0}
    KEY_SEP = string(sep)
}
