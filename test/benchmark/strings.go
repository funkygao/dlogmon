// benchmark strings.Contains
package main

import (
    "strings"
    "time"
    "fmt"
    fs "index/suffixarray"
)

const (
    LINE = ">121016-140231 192.168.100.123 3309 KProxy KXI.SQA /SAMPLE:1/S T=0.000 9999/127.0.0.1:26630 1065 Q=LCache.get X{CALLER^GET+3g.kaixin001.com/parking/index.php+65ced944} {key^uobject_s_user_info:70839223:70839223; expire^3600} A=0 {value^~601!C:14:'CUObjectResult':573:`7B`03&result`03&fields`02#uid%email$nick)real_name&gender(birthday$city%mtime(hometown$logo(privacy1(privacy2*lunarbirth(privacy3(privacy4%state)stateTime%atime`01$rows`02`02·×ã¡@5sunjiang21539@tom.com"
    LOOPS = 1000000
)

func benchFast(line string, substr string) time.Duration {
    index := fs.New([]byte(line))
    start := time.Now()
    for i:=0; i<LOOPS; i++ {
        index.Lookup([]byte(substr), 1)
    }
    end := time.Now()
    delta := end.Sub(start)
    fmt.Printf("%10s: %20s\t%16s %10s\n", "fast", substr, delta, delta/LOOPS)
    return delta
}

func bench(line string, substr string) time.Duration {
    start := time.Now()
    for i:=0; i<LOOPS; i++ {
        strings.Contains(line, substr)
    }
    end := time.Now()
    delta := end.Sub(start)
    fmt.Printf("%10s: %20s\t%16s %10s\n", "contains", substr, delta, delta/LOOPS)
    return delta
}

func main() {
    fmt.Printf("%s\n%s\n", time.Now(), strings.Repeat("=", 80))

    strs := []string{
        "AMF_SLOW",
        "PHP.CDlog",
        "Q=DLog.log",
        "123.123",
        "123.123.123.123",
        "100.123"}
    for _, s := range strs {
        bench(LINE, s)
        benchFast(LINE, s)
    }

    fmt.Printf("%s\n%s\n", strings.Repeat("=", 80), time.Now())
}

