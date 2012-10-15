// line reader
package dlog

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "os/exec"
    "strings"
    "strconv"
)

var lineValidatorRegexes = [][]string{{"AMF_SLOW", "100.123", "PHP.CDlog"}, {"Q=DLog.log"}}

// AMF_SLOW dlog analyzer
type AmfDlog struct {
    Dlog
}

// a single line meta info
type amfRequest struct {
    Request
    class, method, args string
    time int16
}

// parse a line into meta info
// ret -> valid line?
func (req *amfRequest) ParseLine(line string) {
    // major parts seperated by space
    parts := strings.Split(line, " ")

    // uri related
    uriInfo := strings.Split(parts[5], "+")
    if len(uriInfo) < 3 {
        log.Fatal(line)
    }
    req.http_method, req.uri, req.rid = uriInfo[0], uriInfo[1], uriInfo[2]

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
        log.Fatal(line, err)
    }
    req.time = int16(time)
    req.class = callInfo[3]
    if len(callInfo) > 10 {
        req.method = callInfo[10]
    }
}

// better printable Request
func (req *amfRequest) String() string {
    return fmt.Sprintf("{http:%s uri:%s rid:%s class:%s method:%s time:%d args:%s}",
        req.http_method, req.uri, req.rid, req.class, req.method, req.time, req.args)
}

func (dlog AmfDlog) ReadLines() {
    run := exec.Command(LZOP_CMD, LZOP_OPTION, dlog.filename)
    out, err := run.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }
    if err := run.Start(); err != nil {
        log.Fatal(err)
    }

    if dlog.options.mapper != "" {
        //mapper := exec.Command(dlog.options.mapper)
    }

    inputReader := bufio.NewReader(out)
    lineCount := 0
    for {
        line, err := inputReader.ReadString(EOL)
        if err != nil {
            if err != io.EOF {
                log.Fatal(err)
            }

            break
        }

        lineCount += 1

        if !dlog.IsLineValid(line) {
            continue
        }

        // extract info from this line
        dlog.OperateLine(line)
    }

    if err := run.Wait(); err != nil {
        log.Fatal(err)
    }

    dlog.chLines <- lineCount
}

func (dlog AmfDlog) IsLineValid(line string) bool {
    // must exists
    for _, regex := range lineValidatorRegexes[0] {
        if !strings.Contains(line, regex) {
            return false
        }
    }

    // must not exists
    for _, regex := range lineValidatorRegexes[1] {
        if strings.Contains(line, regex) {
            return false
        }
    }

    return true
}

// operate on a valid dlog line
func (dlog *AmfDlog) OperateLine(line string) {
    req := new(amfRequest)
    req.ParseLine(line)

    dlog.lock.Lock()
    defer dlog.lock.Unlock()
    fmt.Printf("%6d%25s %35s   %s\n", req.time, req.class, req.method, req.uri)
}

