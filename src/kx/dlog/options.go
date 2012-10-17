package dlog

import (
    "flag"
    "fmt"
    T "kx/trace"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

// CLI options object
type Option struct {
    files []string
    debug bool
    trace bool
    verbose bool
    version bool
    mapper string
    reducer string
    kind string
    logfile string
}

// Printable Option
func (this *Option) String() string {
    return fmt.Sprintf("Option{files: %#v debug:%#v mapper:%s reducer:%s}", this.files, this.debug,
        this.mapper, this.reducer)
}

// Names of the dlog files to be analyzed
func (this *Option) Files() []string {
    return this.files
}

// Kind of the current dlog
func (this *Option) Kind() string {
    return this.kind
}

// Does CLI ask for dlogmon version info?
func (this *Option) Version() bool {
    return this.version
}

// Parse CLI options
func ParseFlags() *Option {
    d := flag.String("D", "", "day of dlog[default today] e,g 121005")
    h := flag.String("H", "10", "hour of dlog[default 10] e,g 9-11")
    f := flag.String("f", "", "specify a single dlog file to analyze")
    verbose := flag.Bool("v", false, "verbose")
    version := flag.Bool("V", false, "show version")
    kind := flag.String("k", "amf", "what kind of content to scan in dlog[amf|xxx]")
    debug := flag.Bool("d", false, "debug mode")
    mapper := flag.String("mapper", "", "let a runnable script be the mapper")
    reducer := flag.String("reducer", "", "let a runnable script be the reducer")
    logfile := flag.String("l", "", "log file path, default stderr")
    trace := flag.Bool("t", false, "trace each func call")

    flag.Parse()

    option := new(Option)
    if *f != "" {
        option.files = []string{*f}
    }
    option.debug = *debug
    option.mapper = *mapper
    option.reducer = *reducer
    option.kind = *kind
    option.version = *version
    option.verbose = *verbose
    option.logfile = *logfile
    option.trace = *trace
    if option.trace {
        T.Enable()
    }

    // day
    dir := *d
    if dir == "" {
        // default today
        now := time.Now()
        year, month, day := now.Date()
        dir = fmt.Sprintf("%4d%02d%02d", year, month, day)
        dir = dir[2:] // 20120918 -> 120918
    }

    // hour span
    var h1, h2 int
    var err error
    var hp []string // hour parts
    if strings.Contains(*h, FLAG_TIMESPAN_SEP) {
        hp = strings.SplitN(*h, FLAG_TIMESPAN_SEP, 2)
    } else {
        hp = []string{*h, *h}
    }

    h1, err = strconv.Atoi(strings.TrimSpace(hp[0]))
    if err != nil {
        panic(err)
    }
    h2, err = strconv.Atoi(strings.TrimSpace(hp[1]))
    if err != nil {
        panic(err)
    }
    if h1 > h2 {
        fmt.Println("Invalid hour option:", *h)
        os.Exit(1)
    }
    globs := make([]string, 0)
    for i:=h1; i<=h2; i++ {
        globs = append(globs, fmt.Sprintf("%s%s/*.%s-%02d*", DLOG_BASE_DIR, dir, dir, i))
    }

    for _, glob := range globs {
        files, err := filepath.Glob(glob)
        if err != nil {
            panic(err)
        }

        for _, file := range files {
            option.files = append(option.files, file)
        }
    }

    return option
}

