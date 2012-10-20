package size

import (
    "fmt"
    "testing"
)

func TestByteSize(t *testing.T) {
    var b ByteSize = 12121212212
    s := fmt.Sprintf("%s", b)
    expected := "11.29GB"
    if s != expected {
        t.Error(expected, b)
    }
}

func TestConsts(t *testing.T) {
    if MB / KB != 1024 {
        t.Error("MB/KB != 1024")
    }
}
