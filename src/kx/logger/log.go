// TODO vitess/relog
package logger

import (
    "flag"
    stdlog "log"
    "sync"
)

var verbosity = flag.Int("verbosity", 0, "logging verbosity (-1==QUIET to 1==DEBUG)")

const (
    QUIET = iota - 1
    NORMAL
    DEBUG
)

type Level int

var mutex sync.RWMutex

func test(level Level) bool {
    mutex.RLock()
    defer mutex.RUnlock()
    return Level(*verbosity) >= level
}

func SetVerbosity(level Level) {
    mutex.Lock()
    defer mutex.Unlock()
    *verbosity = int(level)
}

// Not silenceable, terminates
var (
    Fatal   = stdlog.Fatal
    Fatalf  = stdlog.Fatalf
    Fatalln = stdlog.Fatalln
)

// Not silenceable, does not terminate
func Always(args ...interface{})            { stdlog.Print(args...) }
func Alwaysf(s string, args ...interface{}) { stdlog.Printf(s, args...) }
func Alwaysln(args ...interface{})          { stdlog.Println(args...) }

// Not silenceable, does not terminate (might revise this in the future)
var (
    Error   = Always
    Errorf  = Alwaysf
    Errorln = Alwaysln
)

// Debug level
func Debug(args ...interface{}) {
    if test(DEBUG) {
        stdlog.Print(args...)
    }
}

func Debugf(s string, args ...interface{}) {
    if test(DEBUG) {
        stdlog.Printf(s, args...)
    }
}

func Debugln(args ...interface{}) {
    if test(DEBUG) {
        stdlog.Println(args...)
    }
}

// Normal level
func Print(args ...interface{}) {
    if test(NORMAL) {
        stdlog.Print(args...)
    }
}

func Printf(s string, args ...interface{}) {
    if test(NORMAL) {
        stdlog.Printf(s, args...)
    }
}

func Println(args ...interface{}) {
    if test(NORMAL) {
        stdlog.Println(args...)
    }
}

// Specified level
func Log(level Level, args ...interface{}) {
    if test(level) {
        stdlog.Print(args...)
    }
}
func Logf(level Level, s string, args ...interface{}) {
    if test(level) {
        stdlog.Printf(s, args...)
    }
}
func Logln(level Level, args ...interface{}) {
    if test(level) {
        stdlog.Println(args...)
    }
}
