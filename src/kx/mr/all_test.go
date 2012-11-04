package mr

import (
    "fmt"
    "github.com/bmizerany/assert"
    "testing"
)

func TestAppend(t *testing.T) {
    d := NewKeyValues()
    d.Append("a", 1.0)
    d.Append("a", 2.0)
    d.Append("b", 3.0)

    f := ConvertAnySliceToFloat(d["a"])
    assert.Equal(t, f, []float64{1.0, 2.0})
}

func TestAppendVariadics(t *testing.T) {
    d := NewKeyValues()
    d.Append("a", 1.0, 2., 3.)
    d.Append("b", 3.0)
    fmt.Println(d)

    f := ConvertAnySliceToFloat(d["a"])
    assert.Equal(t, f, []float64{1., 2.0, 3})
    f = ConvertAnySliceToFloat(d["b"])
    assert.Equal(t, f, []float64{3})
}

func TestKeyValueSortType(t *testing.T) {
    kv := NewKeyValue()
    kv["a"] = 1
    kv["b"] = 2
    assert.Equal(t, SORT_BY_VALUE, kv.sortType())

    kv = NewKeyValue()
    kv[[...]string{"avg", "www"}] = 1
    kv[[...]string{"avg", "game"}] = 2
    assert.Equal(t, SORT_BY_VALUE, kv.sortType())
}
