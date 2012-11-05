package util

import (
    "github.com/bmizerany/assert"
    "testing"
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
