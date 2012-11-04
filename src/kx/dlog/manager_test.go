package dlog

import (
    "testing"
)

func mockOption() *Option {
    return new(Option)
}

func mockManager() *Manager {
    return NewManager(mockOption())
}

func TestNewManager(t *testing.T) {
    m := mockManager()
    if m.RawLines != 0 {
        t.Error("rawLines")
    }
}
