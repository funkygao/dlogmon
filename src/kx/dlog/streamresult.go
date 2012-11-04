package dlog

import (
    "encoding/json"
)

var emptyStreamResult StreamResult = ""

// Empty or invalid?
func (this StreamResult) Empty() bool {
    return this == "\n" || this == ""
}

// Decode string into a record
func (this StreamResult) Decode(record interface{}) (err error) {
    err = json.Unmarshal([]byte(this), record)
    return
}
