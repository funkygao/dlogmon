// codable stream or pipe
package stream

import (
	"bufio"
	"os"
	"os/exec"
)

// Constructor factory
func NewStream(name string, arg ...string) *Stream {
	return &Stream{name: name, arg: arg}
}

func (this *Stream) plainFileMode() bool {
	return len(this.arg) == 0
}

// open stream
func (this *Stream) Open() {
	if this.plainFileMode() {
		// direct file open instead of exec pipe stream
		file, err := os.Open(this.name)
		if err != nil {
			panic(this.name + " not exist")
		}

		this.reader = bufio.NewReader(file)
		return
	}

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
	if this.plainFileMode() {
		return
	}

	if err := this.cmd.Wait(); err != nil {
		panic(err)
	}
}
