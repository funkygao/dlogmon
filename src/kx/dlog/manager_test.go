package dlog

import (
    "fmt"
    "github.com/bmizerany/assert"
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
    if m.rawLines != 0 {
        t.Error("rawLines")
    }
}

func TestString(t *testing.T) {
    m := mockManager()
    expected := `Manager{&dlog.Option{files:[]string(nil), debug:false, trace:false, verbose:false, version:false, Nworkers:0, tick:0, cpuprofile:"", memprofile:"", mapper:"", reducer:"", kind:"", conf:(*config.Config)(nil)}}`
    got := fmt.Sprintf("%s", m)
    assert.Equal(t, expected, got)
}
