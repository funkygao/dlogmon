package main

import (
    "fmt"
    "github.com/funkygao/cmd"
    "kx/dlog"
    "kx/mr"
    "strconv"
)

// Current CLI params
type dlogmonCli struct {
    group, sortCol string
    top int
}

var (
    cli cmd.Cmd
    dlogCli *dlogmonCli
    result mr.KeyValue  // reducer's result, maybe grouped by key
    worker dlog.IWorker
)

func init() {
    dlogCli = new(dlogmonCli)
    cli = cmd.New(dlogCli)
    cli.Intro = DLOGMON_INTRO
    cli.Prompt = DLOOGMON_PROMPT
}

func cmdloop(w dlog.IWorker, reduceResult mr.KeyValue) {
    result = reduceResult
    worker = w
    cli.Cmdloop()
}

func (this dlogmonCli) Help() {
    fmt.Println(`Available commands:
group sort top raw worker show

Use "help | h [command]" for more information about a command.`)
}

func (this dlogmonCli) Help_group() {
    fmt.Println("group {group_name}")
}

func (this dlogmonCli) Help_sort() {
    fmt.Println("sort {col_name}")
}

func (this dlogmonCli) Help_top() {
    fmt.Println("top {N}")
}

func (this dlogmonCli) Help_show() {
    fmt.Println(`show
show current report`)
}

func (this dlogmonCli) Help_raw() {
    fmt.Println(`raw
output raw reducer's data`)
}

func (this dlogmonCli) Help_worker() {
    fmt.Println(`worker
show current worker info`)
}

func (this dlogmonCli) Do_raw() {
    for k, v := range result {
        fmt.Printf("key=> %#v\n", k)
        fmt.Printf("val=> %#v\n", v)
        fmt.Println()
    }
}

func (this dlogmonCli) Do_worker() {
    fmt.Printf("%#v\n", worker)
}

func (this dlogmonCli) Do_show() {
}

func (this dlogmonCli) Do_top(n string) {
    t, e := strconv.Atoi(n)
    if e != nil {
        panic("top {N}, N must be unsined integer")
    }

    if t > 0 {
        this.top = t // remember this
        result.ExportResult(worker, this.top)
    }
}

func (this dlogmonCli) Do_sort(col string) {
    this.sortCol = col
}

func (this dlogmonCli) Do_group(group string) {
    this.group = group
}
