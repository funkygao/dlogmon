// codable stream or pipe
package stream

import (
    "bufio"
    "errors"
    "fmt"
    "os"
    "os/exec"
)

// Constructor factory
func NewStream(mode StreamMode, name string, arg ...string) *Stream {
    return &Stream{name: name, arg: arg, mode: mode}
}

func (this *Stream) openPlainFile() error {
    file, err := os.Open(this.name)
    if err != nil {
        return err
    }

    this.reader = bufio.NewReader(file)
    return nil
}

func (this *Stream) openExecPipe() error {
    this.cmd = exec.Command(this.name, this.arg...)

    var err error
    this.pr, err = this.cmd.StdoutPipe()
    if err != nil {
        return err
    }

    this.pw, err = this.cmd.StdinPipe()
    if err != nil {
        return err
    }

    // startup
    if err := this.cmd.Start(); err != nil {
        return err
    }

    // prepare the reader/writer
    this.reader = bufio.NewReader(this.pr)
    this.writer = bufio.NewWriter(this.pw)
    return nil
}

// open stream
func (this *Stream) Open() error {
    switch this.mode {
    case EXEC_PIPE:
        return this.openExecPipe()
    case PLAIN_FILE:
        return this.openPlainFile()
    }

    return errors.New(fmt.Sprintf("non supported mode, valid modes: %v", this.mode.Valids()))
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
func (this *Stream) Close() error {
    switch this.mode {
    case EXEC_PIPE:
        // close my writer stream so that client can get EOF
        this.pw.Close()
        // wait for client to exit
        if err := this.cmd.Wait(); err != nil {
            return err
        }
    }

    return errors.New(fmt.Sprintf("current mode:%s can't be closed", this.mode.String()))
}
