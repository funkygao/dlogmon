package mapreduce

import (
    T "testing"
)

func TestFoo(t *T.T) {
    if 1 < 0 {
        t.Errorf("1<0")
    }
}

