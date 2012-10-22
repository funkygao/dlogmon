package dlog

import (
    "github.com/bmizerany/assert"
    "testing"
)

func TestIntSum(t *testing.T) {
    vals := []Any{1, 5, 9, 12}
    assert.Equal(t, intSum(vals), 27)
}

func TestIntAvg(t *testing.T) {
    vals := []Any{1, 5, 9, 12}
    assert.Equal(t, intAvg(vals), 6.75)
}

