package dlog

import (
    "fmt"
    "testing"
)

func newOption() *Option {
    return new(Option)
}

func TestEmptyOptionString(t *testing.T) {
    option := newOption()
    s := fmt.Sprintf("%s", option)
    expected := "Option{files: []string(nil) debug:false mapper: reducer:}"
    if s != expected {
        t.Error(s, expected)
    }
}

func TestOptionFiles(t *testing.T) {
    option := newOption()
    var files = []string{"a", "c", "cd"}
    option.files = files
    if len(option.Files()) != len(option.files) {
        t.Error("Option.files wrong!")
    }

    for i, file := range option.files {
        if file != files[i] {
            t.Errorf("Option file: %d is wrong\n", i)
        }
    }
}
