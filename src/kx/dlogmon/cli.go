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
    result mr.KeyValue  // reducer's result, maybe grouped by key
    worker dlog.IWorker
}

var cli *dlogmonCli

func cliCmdloop(worker dlog.IWorker, reduceResult mr.KeyValue) {
    cli = new(dlogmonCli)
    cli.result = reduceResult
    cli.worker = worker

    cmd := cmd.New(cli)
    cmd.Intro = DLOGMON_INTRO
    cmd.Prompt = DLOGMON_PROMPT

    // startup the loop
    cmd.Cmdloop()
}

// universal help
func (this dlogmonCli) Help() {
    fmt.Println(DLOGMON_HELP)
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
    for k, v := range this.result {
        fmt.Printf("key=> %#v\n", k)
        fmt.Printf("val=> %#v\n", v)
        fmt.Println()
    }
}

func (this dlogmonCli) Do_worker() {
    fmt.Printf("%#v\n", this.worker)
}

func (this *dlogmonCli) Do_top(n string) {
    t, e := strconv.Atoi(n)
    if e != nil {
        panic("top {N}, N must be unsined integer")
    }

    if t > 0 {
        this.top = t // remember this
        this.render()
    }
}

func (this dlogmonCli) render() {
    this.result.ExportResult(this.worker, this.group, this.sortCol, this.top)
}

func (this *dlogmonCli) Do_sort(col string) {
    this.sortCol = col
    this.render()
}

func (this *dlogmonCli) Do_group(group string) {
    this.group = group
    this.render()
}

func (this dlogmonCli) Do_status() {
    fmt.Printf("group=%s, sort column=%s, top=%d\n", this.group, this.sortCol, this.top)
}
