package dlog

import (
    "os"
    "syscall"
)

var (
    workerConstructors = map[string]WorkerConstructor{
        "amf": NewAmfWorker,
        "noop": NewNoopWorker}

    amfLineValidatorRegexes = [...][]string{
        {"AMF_SLOW", "PHP.CDlog"}, // must exists
        {"Q=DLog.log"}}            // must not exists

    caredSignals = []os.Signal{
        syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT,
        syscall.SIGHUP, syscall.SIGSTOP, syscall.SIGQUIT}

    skippedSignals = [...]syscall.Signal{
        syscall.SIGHUP,
        syscall.SIGSTOP,
        syscall.SIGQUIT}
)
