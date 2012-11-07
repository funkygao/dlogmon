package main

import (
    "fmt"
    "github.com/funkygao/cmd"
    "kx/dlog"
    "kx/mr"
    T "kx/trace"
    "path/filepath"
    "runtime"
    "strconv"
)

// Current CLI params
type dlogmonCli struct {
    group, sortCol string
    top            int
    groups         map[string]string
    result         mr.KeyValue // reducer's result, maybe grouped by key
    worker         dlog.IWorker
}

var cli *dlogmonCli

func cliCmdloop(worker dlog.IWorker, reduceResult mr.KeyValue) {
    cli = new(dlogmonCli)
    cli.result = reduceResult
    cli.worker = worker
    cli.groups = make(map[string]string)
    for i, g := range cli.result.Groups() {
        cli.groups[strconv.Itoa(i)] = g
    }

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

func (this dlogmonCli) Help_groups() {
    fmt.Println("groups\nshow all groups")
}

func (this dlogmonCli) Help_save() {
    fmt.Println("save {filename}\nsave current parsed result to a file")
}

func (this dlogmonCli) Help_load() {
    fmt.Println("load {filename}\nload parsed result from a file")
}

func (this dlogmonCli) Help_hist() {
    fmt.Println("hist\nloadable history parsed results")
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
    fmt.Println(this.worker)
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

func (this dlogmonCli) Do_show() {
    this.render()
}

func (this dlogmonCli) Do_gc() {
    runtime.GC()
}

func (this *dlogmonCli) Do_group(group string) {
    this.group = this.groups[group]
    this.render()
}

func (this dlogmonCli) Do_groups() {
    for k, v := range this.groups {
        fmt.Printf("%2s => %s\n", k, v)
    }
}

func (this dlogmonCli) Do_status() {
    fmt.Printf("group=%s, sort column=%s, top=%d\n", this.group, this.sortCol, this.top)
    fmt.Printf("mem:%v, goroutines:%v\n", T.MemAlloced(), runtime.NumGoroutine())
}

func (this dlogmonCli) Do_save(filename string) {
}

func (this dlogmonCli) Do_load(filename string) {
}

func (this dlogmonCli) Do_hist() {
    fmt.Println("Available history files:")
    files, _ := filepath.Glob(dlog.VarDir + "/*.gob")
    var found bool
    for _, file := range files {
        fmt.Println(filepath.Base(file))
        found = true
    }
    if !found {
        fmt.Println("not found")
    }
}
