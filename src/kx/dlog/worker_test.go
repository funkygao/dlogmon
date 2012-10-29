package dlog

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestWorkerConst(t *testing.T) {
	assert.Equal(t, LZOP_CMD, "lzop")
	assert.Equal(t, LZOP_OPTION, "-dcf")
}
