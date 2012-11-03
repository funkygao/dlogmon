package dlog

import (
    "strings"
)

func isSamplerHostLine(line string) bool {
    if !strings.Contains(line, SAMPLER_HOST) {
        return false
    }
    return true
}
