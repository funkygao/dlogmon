package net

import (
    "log"
    "net"
)

type Server struct {
    addr string
}

func (this Server) Listen() {
    l, err := net.Listen("tcp", this.addr)
    if err != nil {
        panic(err)
    }

    for {
        conn, err := l.Accept()
        if err != nil {
            log.Println(err)
            continue
        }

        go onRequest(conn)
    }
}

func onRequest(conn net.Conn) {
    conn.Close()
}
