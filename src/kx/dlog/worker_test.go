package dlog

import "testing"

func TestWorkerConst(t *testing.T) {
    if LZOP_CMD != "lzop" {
        t.Error(LZOP_CMD)
    }

    if LZOP_OPTION != "-dcf" {
        t.Error(LZOP_OPTION)
    }
}
