package main

import (
    "kx/util"
    "time"
)

func main() {
    evalFunc := func(x util.Any) (retVal, nextVal util.Any) {
        v := x.(int)
        return v, v + 2
    }

    f := util.BuildLazyIntEvaluator(evalFunc, 0)

    for {
        println(f())
        time.Sleep(1e6)
    }

}
