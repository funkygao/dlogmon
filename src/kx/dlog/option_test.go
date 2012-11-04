package dlog

import (
    "fmt"
    "github.com/bmizerany/assert"
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
    assert.Equal(t, len(option.files), len(option.Files()))

    for i, file := range option.files {
        if file != files[i] {
            t.Errorf("Option file: %d is wrong\n", i)
        }
    }
}
