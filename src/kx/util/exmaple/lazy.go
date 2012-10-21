package main

import (
    "kx/util"
    "time"
)

func main() {
    evalFunc := func(x util.Any) (util.Any, util.Any) {
        v := x.(int)
        return v, v + 1
    }

    lf := util.BuildLazyIntEvaluator(evalFunc, 0)

    for {
        println(lf())
        time.Sleep(1e8)
    }

}
