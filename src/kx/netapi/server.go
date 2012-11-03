package netapi

import (
    "net"
    "net/http"
    "net/rpc"
)

func StartServer(i interface{}) error {
    if e := rpc.Register(i); e != nil {
        return e
    }
    rpc.HandleHTTP()
    l, e := net.Listen(PROTO, ADDRS)
    if e != nil {
        return e
    }
    go http.Serve(l, nil)

    return nil
}
