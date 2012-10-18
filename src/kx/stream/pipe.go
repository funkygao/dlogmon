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

// constructor
func NewStream(name string, arg ...string) *Stream {
    return &Stream{name: name, arg: arg}
}

// open stream
func (this *Stream) Open() {
    this.cmd = exec.Command(this.name, this.arg...)

    // stdout pipe
    out, err := this.cmd.StdoutPipe()
    if err != nil {
        panic(err)
    }

    // stdin pipe
    in, err := this.cmd.StdinPipe()
    if err != nil {
        panic(err)
    }

    // startup
    if err := this.cmd.Start(); err != nil {
        panic(err)
    }

    // prepare the reader/writer
    this.reader = bufio.NewReader(out)
    this.writer = bufio.NewWriter(in)
}

// get reader to read from the pipe output
func (this Stream) Reader() *bufio.Reader {
    return this.reader
}

// get writer to write to the pipe input
func (this Stream) Writer() *bufio.Writer {
    return this.writer
}

// close the stream
func (this *Stream) Close() {
    if err := this.cmd.Wait(); err != nil {
        panic(err)
    }
}
