package util

import (
    "os"
)

func FileExists(fileOrDir string) bool {
    var err error
    var f *os.File
    if f, err = os.Open(fileOrDir); err != nil {
        return false
    }

    f.Close()
    return true
}
