package main

import (
    "fmt"
    "github.com/funkygao/cmd"
    "kx/mr"
)

type dlogmonCli struct {
}

var cli cmd.Cmd
var result mr.KeyValue

func init() {
    cli = cmd.New(new(dlogmonCli))
    cli.Intro = DLOGMON_INTRO
    cli.Prompt = DLOOGMON_PROMPT
}

func cmdloop(reduceResult mr.KeyValue) {
    result = reduceResult
    cli.Cmdloop()
}

func (this dlogmonCli) Help() {
    fmt.Println(`Available commands:
group sort top raw

Use "help [command]" for more information about a command.`)
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

func (this dlogmonCli) Help_raw() {
    fmt.Println(`raw
output raw reducer's data`)
}

func (this dlogmonCli) Do_raw() {
    for k, v := range result {
        fmt.Printf("key=> %#v\n", k)
        fmt.Printf("val=> %#v\n", v)
        fmt.Println()
    }
}

func (this dlogmonCli) Do_top(n string) {
}

func (this dlogmonCli) Do_sort(col string) {
}

func (this dlogmonCli) Do_group(group string) {
}
