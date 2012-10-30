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
    assert.Equal(t, d,
        KeyValues{"a": []float64{1.0, 2.0}, "b": []float64{3.0}})
}

func TestAppendSlice(t *testing.T) {
    d := NewKeyValues()
    d.Append("a", 2.0)
    d.AppendSlice("a", []float64{3.1, 4.1})
    d.Append("b", 3.0)
    assert.Equal(t, d, KeyValues{"a": []float64{2.0, 3.1, 4.1}, "b": []float64{3.0}})
}
