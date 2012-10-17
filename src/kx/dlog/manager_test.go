package dlog

import (
    "testing"
    "fmt"
)

func mockOption() *Option {
    return new(Option)
}

func mockManager() *Manager {
    return NewManager(mockOption())
}

func TestNewManager(t *testing.T) {
    m := mockManager()
    if m.rawLines != 0 {
        t.Error("rawLines")
    }
}

func TestString(t *testing.T) {
    m := mockManager()
    expected := "shit"
    str := fmt.Sprintf("%s", m)
    if expected != str {
        t.Error("expected:", expected, "real:", str)
    }
    
}

