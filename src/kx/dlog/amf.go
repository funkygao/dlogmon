// line reader
package dlog

import (
    "bufio"
    "io"
    "log"
    "os/exec"
    "strings"
)

// AMF_SLOW dlog analyzer
type AmfDlog struct {
    Dlog
}

func (dlog AmfDlog) ReadLines() {
    run := exec.Command(LZOP_CMD, LZOP_OPTION, dlog.filename)
    out, err := run.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }

    inputReader := bufio.NewReader(out)
    for {
        line, err := inputReader.ReadString(EOL)
        if err != nil {
            if err != io.EOF {
                log.Fatal(err)
            }

            break
        }

        if !dlog.IsLineValid(line) {
            continue
        }

        // extract info from this line
        dlog.OperateLine(line)
    }

    dlog.chEof <- true
}

func (dlog AmfDlog) IsLineValid(line string) bool {
    regexes := []string{"AMF_SLOW", "100.123", "PHP.CDlog"}
    for _, regex := range regexes {
        if !strings.Contains(line, regex) {
            return false
        }
    }

    return true
}

func (dlog *AmfDlog) OperateLine(line string) {
    print(line)
}

