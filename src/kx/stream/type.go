// codable stream or pipe
package stream

import (
	"bufio"
	"os/exec"
)

// Stream data
type Stream struct {
	name   string   // command name
	arg    []string // command arguments
	cmd    *exec.Cmd
	reader *bufio.Reader
	writer *bufio.Writer
}
