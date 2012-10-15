package dlog

import (
    "testing"
    "fmt"
)

func newOptions() *Options {
    return new(Options)
}

func TestEmptyOptionsString(t *testing.T) {
    options := newOptions()
    s := fmt.Sprintf("%s", options)
    expected := "{files: []string(nil) debug:false mapper: reducer:}"
    if s != expected {
        t.Error(s, expected)
    }
}

func TestOptionsGetFiles(t *testing.T) {
    options := newOptions()
    var files = []string{"a", "c", "cd"}
    options.files = files
    if len(options.GetFiles()) != len(options.files) {
        t.Error("Options.files wrong!")
    }

    for i, file := range options.files {
        if file != files[i] {
            t.Errorf("Options file: %d is wrong\n", i)
        }
    }
}

