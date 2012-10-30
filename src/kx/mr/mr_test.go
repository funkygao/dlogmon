package mr

import (
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

func TestAppendSlice(t *testing.T) {
    d := NewKeyValues()
    d.Append("a", 2.0)

    s := make([]interface{}, 2)
    s[0], s[1] = 3.1, 4.1
    d.AppendSlice("a", s)

    d.Append("b", 3.0)
    f := ConvertAnySliceToFloat(d["a"])
    assert.Equal(t, f, []float64{2, 3.1, 4.1})
}
