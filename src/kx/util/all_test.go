package util

import (
    "github.com/bmizerany/assert"
    "testing"
    "unsafe"
)

func TestSet(t *testing.T) {
    s := NewSet()
    s.Add(5)
    assert.Equal(t, false, s.Contains(4))
    assert.Equal(t, true, s.Contains(5))
    assert.Equal(t, 1, s.Len())
    s.Add("string")
    assert.Equal(t, 2, s.Len())
    s.Add("string")
    assert.Equal(t, 2, s.Len())
}

func TestEncodeDecodeSlice(t *testing.T) {
    x := []string{"funky", "gao peng", "kaixin"}
    se, e := EncodeStrSlice(x)
    assert.Equal(t, e, nil)

    de, e := DecodeStrToSlice(se)
    assert.Equal(t, e, nil)

    assert.Equal(t, len(x), len(de))

    for i, v := range x {
        assert.Equal(t, de[i], v)
    }
}

func BenchmarkGobEncodeSlice(b *testing.B) {
    x := []string{"123456", "7890", "123456"}
    for i := 0; i < b.N; i++ {
        EncodeStrSlice(x)
    }
    b.SetBytes(int64(unsafe.Sizeof(x)))
}

func BenchmarkGobDecodeSlice(b *testing.B) {
    x := []string{"funky", "gao peng", "kaixin"}
    se, _ := EncodeStrSlice(x)
    for i := 0; i < b.N; i++ {
        DecodeStrToSlice(se)
    }
    b.SetBytes(int64(unsafe.Sizeof(x)))
}
