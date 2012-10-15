package dlog

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    DLOG_BASE_DIR = "/kx/dlog/"
    SPAN_SEP = "-"
)

// CLI options object
type Options struct {
    files []string
    debug bool
    mapper string
    reducer string
}

func (opt *Options) String() string {
    return fmt.Sprintf("{files: %#v debug:%#v mapper:%s reducer:%s}", opt.files, opt.debug, opt.mapper, opt.reducer)
}

func (options *Options) GetFiles() []string {
    return options.files
}

// parse CLI options
func ParseFlags() *Options {
    options := new(Options)

    d := flag.String("D", "", "day of dlog[default today] e,g 121005")
    h := flag.String("H", "10", "hour of dlog[default 10] e,g 9-11")
    f := flag.String("f", "", "specify a single dlog file to analyze")
    debug := flag.Bool("d", false, "debug mode")
    mapper := flag.String("mapper", "", "let a runnable script be the mapper")
    reducer := flag.String("reducer", "", "let a runnable script be the reducer")
    flag.Parse()

    options.debug = *debug
    options.mapper = *mapper
    options.reducer = *reducer

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
    if strings.Contains(*h, SPAN_SEP) {
        hp = strings.SplitN(*h, SPAN_SEP, 2)
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
            options.files = append(options.files, file)
        }
    }

    if *f != "" {
        options.files = []string{*f}
    }

    return options
}

