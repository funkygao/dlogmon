package dlog

import (
    "errors"
    "flag"
    "fmt"
    T "kx/trace"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

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

// Cpuprofile name
func (this *Option) Cpuprofile() string {
    return this.cpuprofile
}

// Memprofile name
func (this *Option) Memprofile() string {
    return this.memprofile
}

// Does CLI ask for dlogmon version info?
func (this *Option) Version() bool {
    return this.version
}

// Parse CLI options
func ParseFlags() (*Option, error) {
    var (
        d          = flag.String("D", "", "day of dlog[default today] e,g 121005")
        h          = flag.String("H", "10", "hour of dlog[default 10] e,g 9-11")
        f          = flag.String("f", "", "specify a single dlog file to analyze")
        verbose    = flag.Bool("v", false, "verbose")
        version    = flag.Bool("V", false, "show version")
        rpc        = flag.Bool("rpc", false, "start rpc server")
        progress   = flag.Bool("progress", true, "show progress bar")
        kind       = flag.String("k", "amf", "what kind of content to scan in dlog[amf|xxx]")
        cpuprofile = flag.String("cpuprofile", "", "write cpu profile to a file for pprof")
        memprofile = flag.String("memprofile", "", "write cpu profile to a file for pprof")
        nworkers   = flag.Int("n", 10, "how many concurrent workers permitted")
        debug      = flag.Bool("d", false, "debug mode")
        mapper     = flag.String("mapper", "", "let a runnable script be the mapper")
        reducer    = flag.String("reducer", "", "let a runnable script be the reducer")
        conf       = flag.String("conf", "conf/dlogmon.ini", "conf file path")
        tick       = flag.Int("tick", 0, "tick in ms")
        filemode   = flag.Bool("filemode", false, "input is plain text file(s)")
        trace      = flag.Bool("t", false, "trace each func call")
    )

    flag.Usage = func() {
        const usage = `dlogmon - kaixin dlog parser and report generator

Flags:
`
        fmt.Fprint(os.Stderr, usage)
        flag.PrintDefaults()
    }
    flag.Parse()

    option := new(Option)
    option.debug = *debug
    option.progress = *progress
    option.mapper = *mapper
    option.reducer = *reducer
    option.cpuprofile = *cpuprofile
    option.memprofile = *memprofile
    option.kind = *kind
    option.version = *version
    option.rpc = *rpc
    option.filemode = *filemode
    option.tick = *tick
    option.Nworkers = uint8(*nworkers)
    option.verbose = *verbose
    option.conf, _ = loadConf(*conf)
    option.trace = *trace
    if option.trace {
        T.Enable()
    }
    if *f != "" {
        option.files = []string{*f}
        option.Timespan = fmt.Sprintf("%s worker => dlog file: %s", *kind, *f)
        return option, nil
    }

    // 根据指定的时间范围判断分析哪些文件

    // day
    dir := *d
    if dir == "" {
        // default today
        now := time.Now()
        year, month, day := now.Date()
        if now.Hour() < 10 {
            yestoday := now.AddDate(0, 0, -1)
            day = yestoday.Day()
        }
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

    basedir, e := option.conf.String("default", "dlogbase")
    if e != nil {
        basedir = DEFAULT_DLOG_BASE_DIR
    }
    globs := make([]string, 0)
    for i := h1; i <= h2; i++ {
        globs = append(globs, fmt.Sprintf("%s%s/*.%s-%02d*", basedir, dir, dir, i))
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

    if len(option.files) < 1 {
        return option, errors.New("Fatal error: empty dlog")
    }

    option.Timespan = fmt.Sprintf("%s worker => 20%s %02d:00-%02d:00[%d dlog files]", *kind,
        dir, h1, h2+1, len(option.files))

    return option, nil
}
