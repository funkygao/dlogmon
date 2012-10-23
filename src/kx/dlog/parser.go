package dlog

import (
    "strconv"
    "strings"
)

func (this *amfRequest) parseLine(line string) {
    // major parts seperated by space
    parts := strings.Split(line, " ")

    // uri related
    uriInfo := strings.Split(parts[5], "+")
    if len(uriInfo) < 3 {
        panic(uriInfo)
    }
    this.http_method, this.uri, this.rid = uriInfo[0], uriInfo[1], uriInfo[2]

    // class call and args related
    callRaw := strings.Replace(parts[6], "{", "", -1)
    callRaw = strings.Replace(callRaw, "}", "", -1)
    callRaw = strings.Replace(callRaw, "\"", "", -1)
    callRaw = strings.Replace(callRaw, "[", "", -1)
    callRaw = strings.Replace(callRaw, "]", "", -1)
    callRaw = strings.Replace(callRaw, ",", ":", -1)
    callInfo := strings.Split(callRaw, ":")
    time, err := strconv.Atoi(callInfo[1])
    if err != nil {
        println(line)
        panic(err)
    }
    this.time = int16(time)
    this.class = callInfo[3]
    if len(callInfo) > 10 {
        this.method = callInfo[10]
    }
}
