package netapi

import (
    "net"
    "net/rpc"
)

func StartServer() error {
    server := rpc.NewServer()
    l, e := net.Listen(PROTO, ADDRS)
    if e != nil {
        return e
    }
    go server.Accept(l)

    return nil
}

func Register(i interface{}) error {
    return rpc.Register(i)
}
