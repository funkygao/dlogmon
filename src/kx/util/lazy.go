package util

type Any interface{}

type EvalFunc func(Any) (Any, Any)

func BuildLazyEvaluator(evalFunc EvalFunc, initVal Any) func() Any {
    retValChan := make(chan Any)

    loopFuc := func() {
        var nextVal Any = initVal
        var retVal Any
        for {
            retVal, nextVal = evalFunc(nextVal)

            retValChan <- retVal
        }
    }

    retFunc := func() Any {
        return <-retValChan
    }

    go loopFuc()
    return retFunc
}

func BuildLazyIntEvaluator(evalFunc EvalFunc, initState Any) func() int {
    ef := BuildLazyEvaluator(evalFunc, initState)
    return func() int {
        return ef().(int)
    }
}
