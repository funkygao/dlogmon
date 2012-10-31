// codable stream or pipe
package stream

import (
    "bufio"
    "os"
    "os/exec"
)

// Constructor factory
func NewStream(mode StreamMode, name string, arg ...string) *Stream {
    return &Stream{name: name, arg: arg, mode: mode}
}

func (this *Stream) openPlainFile() {
    file, err := os.Open(this.name)
    if err != nil {
        panic(this.name + " not exist")
    }

    this.reader = bufio.NewReader(file)
}

func (this *Stream) openExecPipe() {
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

// open stream
func (this *Stream) Open() {
    switch this.mode {
    case EXEC_PIPE:
        this.openExecPipe()
    case PLAIN_FILE:
        this.openPlainFile()
    }
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
    switch this.mode {
    case EXEC_PIPE:
        if err := this.cmd.Wait(); err != nil {
            panic(err)
        }
    }
}
