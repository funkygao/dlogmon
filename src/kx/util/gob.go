package util

import (
    "bytes"
    "encoding/gob"
)

var (
    buf bytes.Buffer
    enc = gob.NewEncoder(&buf)
    dec = gob.NewDecoder(&buf)
)

func EncodeStrSlice(s []string) (r string, e error) {
    //buf.Reset()
    if e = enc.Encode(s); e != nil {
        return
    }
    r = buf.String()
    return
}

func DecodeStrToSlice(s string) (r []string, e error) {
    //buf.Reset()
    buf.WriteString(s)
    buf.Write([]byte(s))
    r = make([]string, 0)
    if e = dec.Decode(&r); e != nil {
        panic(e)
        return
    }
    return
}
