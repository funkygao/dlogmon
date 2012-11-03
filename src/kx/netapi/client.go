package netapi

import (
    "errors"
    "net/rpc"
    "reflect"
    "sync"
)

var client *rpc.Client

func initClient(addr string) error {
    l := new(sync.Mutex)
    l.Lock()
    defer l.Unlock()

    var e error
    client, e = rpc.DialHTTP(PROTO, addr)
    return e
}

// Call a remote server's method
// Param reply must be pointer
func Call(addr, method string, args interface{}, reply interface{}) error {
    if client == nil {
        initClient(addr)
    }

    if v := reflect.ValueOf(reply); v.Kind() != reflect.Ptr {
        return errors.New("reply must be pointer")
    }

    return client.Call(method, args, reply)
}
