package main

import (
    "fmt"
    "github.com/funkygao/cmd"
)

type dlogmonCli struct {
}

var cli cmd.Cmd

func init() {
    cli = cmd.New(new(dlogmonCli))
    cli.Intro = DLOGMON_INTRO
    cli.Prompt = DLOOGMON_PROMPT
}

func cmdloop() {
    cli.Cmdloop()
}

func (this dlogmonCli) Help() {
    fmt.Println(`Available commands:
group sort top

Use "help [topic]" for more information about that topic.`)
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
