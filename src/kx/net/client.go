package net

import (
	"net"
)

type Client struct {
	addr string
}

func (this Client) Connect() {
	conn, err := net.Dial("tcp", this.addr)
	if err != nil {
		panic(err)
	}
}
