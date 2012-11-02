package main

import (
    "fmt"
    "unsafe"
)

func main() {
    var (
        b bool
        u8 uint8
        i8 int8
        i int
        i64 int64
        s string = "ab"
        sl []string = []string{"a", "b"}
        f func()
    )

    fmt.Println("bool", unsafe.Sizeof(b))
    fmt.Println("uint8", unsafe.Sizeof(u8))
    fmt.Println("int8", unsafe.Sizeof(i8))
    fmt.Println("int", unsafe.Sizeof(i))
    fmt.Println("int64", unsafe.Sizeof(i64))
    fmt.Println("string", unsafe.Sizeof(s))
    fmt.Println("slice", unsafe.Sizeof(sl))
    fmt.Println("func", unsafe.Sizeof(f))
}
