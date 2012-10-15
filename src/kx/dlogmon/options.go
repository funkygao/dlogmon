package main

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
    REGEX_SEP = "|"
)

// CLI options object
type Options struct {
    files []string
    regexs[] string
}

// parse CLI options
func parseFlags() *Options {
    options := new(Options)

    d := flag.String("D", "", "day of dlog[default today] e,g 121005")
    h := flag.String("H", "10", "hour of dlog[default 10] e,g 9-11")
    f := flag.String("f", "", "specify a single dlog file to analyze")
    regex := flag.String("e", "", "line pattern e,g PHP.CDlog|AMF_SLOW")
    flag.Parse()

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

    h1, err = strconv.Atoi(hp[0])
    if err != nil {
        panic(err)
    }
    h2, err = strconv.Atoi(hp[1])
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

    // regex list
    if strings.Contains(*regex, REGEX_SEP) {
        parts := strings.Split(*regex, REGEX_SEP)
        for _, p := range parts {
            options.regexs = append(options.regexs, p)
        }
    } else {
        options.regexs = append(options.regexs, *regex) // a single regex
    }

    return options
}

