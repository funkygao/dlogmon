package dlog

import (
    "encoding/json"
)

var emptyStreamResult StreamResult = ""

// Empty or invalid?
func (this StreamResult) Empty() bool {
    return this == "\n" || this == ""
}

// Decode string into a struct
func (this StreamResult) Decode(r interface{}) (err error) {
    err = json.Unmarshal([]byte(this), r)
    return
}
