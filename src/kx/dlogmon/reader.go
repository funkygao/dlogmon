// line reader
package main

import (
    "bytes"
    "fmt"
    "io"
    "log"
    "os/exec"
)

const (
    LZOP_CMD = "lzop"
    LZOP_OPTION = "-dcf"
    EOL = '\n'
)

func read(file string, line chan string) {
    run := exec.Command(LZOP_CMD, LZOP_OPTION, file)
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

        print(line)
    }
}

