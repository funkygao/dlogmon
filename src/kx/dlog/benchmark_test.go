package dlog

import (
    "sync"
    "testing"
)

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

func BenchmarkMutex(b *testing.B) {
    var lock sync.Mutex
    for i := 0; i < b.N; i++ {
        lock.Lock()
        lock.Unlock()
    }
}
