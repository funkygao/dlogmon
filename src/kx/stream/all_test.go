package stream

import (
    "github.com/bmizerany/assert"
    "testing"
)

func TestStreamModeEnum(t *testing.T) {
    m := EXEC_PIPE
    assert.Equal(t, m, FIRST_MODE)
    assert.Equal(t, LAST_MODE, PLAIN_FILE)
}

func TestStreamModeString(t *testing.T) {
    assert.Equal(t, EXEC_PIPE.String(), "ExecPipe")
    assert.Equal(t, LZOP_FILE.String(), "LzopFile")
    assert.Equal(t, PLAIN_FILE.String(), "PlainText")
}

func TestStreamModeValids(t *testing.T) {
    var mode StreamMode
    assert.Equal(t, mode.Valids(), []string{"ExecPipe", "LzopFile", "PlainText"})
}

