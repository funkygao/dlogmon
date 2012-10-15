// line reader
package dlog

import (
    "bytes"
    "io"
    "log"
    "os/exec"
    "strings"
)

// AMF_SLOW dlog analyzer
type AmfDlog struct {
    Dlog
}

func (dlog *AmfDlog) SetFile(file string) {
    dlog.filename = file
}

func (dlog AmfDlog) ReadLines() {
    run := exec.Command(LZOP_CMD, LZOP_OPTION, dlog.filename)
    var stdout bytes.Buffer
    run.Stdout = &stdout
    if err := run.Run(); err != nil {
        log.Fatal(err)
    }

    for {
        line, err := stdout.ReadString(EOL)
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

