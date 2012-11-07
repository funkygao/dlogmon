// codable stream or pipe
package stream

import (
    "bufio"
    "io"
    "os/exec"
)

// Stream data
type Stream struct {
    name   string   // command name
    arg    []string // command arguments
    cmd    *exec.Cmd
    reader *bufio.Reader
    writer *bufio.Writer
    mode   StreamMode
    pw     io.WriteCloser
    pr     io.ReadCloser
}

type StreamMode uint8
