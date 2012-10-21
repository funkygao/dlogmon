package dlog

import (
    "fmt"
    "sync"
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

var Native_HashVal uint64 = 14695981039346656037

func native_fnv1() {
    Native_HashVal *= 1099511628211
    Native_HashVal ^= 0xff
}

func BenchmarkNative(b *testing.B) {
    b.StopTimer()
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        native_fnv1()
    }
}

func BenchmarkMutext(b *testing.B) {
    var lock sync.Mutex
    for i := 0; i < b.N; i++ {
        lock.Lock()
        lock.Unlock()
    }
}
