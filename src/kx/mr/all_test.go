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

func TestGenericSlice(t *testing.T) {
    ints := []int{3, 2, 8}
    s := GenericSlice(ints)
    assert.Equal(t, 3, len(s))
}
