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
    expected := `Manager{&dlog.Option{files:[]string(nil), debug:false, trace:false, verbose:false, version:false, mapper:"", reducer:"", kind:"", logfile:""}}`
    str := fmt.Sprintf("%s", m)
    if expected != str {
        t.Error("expected:", expected, "real:", str)
    }
    
}

