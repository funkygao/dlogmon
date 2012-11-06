package mr

func less(i, j interface{}, asc bool) bool {
    switch i.(type) {
    case string:
        i, j := i.(string), j.(string)
        if asc {
            return i < j
        } else {
            return i > j
        }
    case float64:
        i, j := i.(float64), j.(float64)
        if asc {
            return i < j
        } else {
            return i > j
        }
    case int:
        i, j := i.(int), j.(int)
        if asc {
            return i < j
        } else {
            return i > j
        }
    case int64:
        i, j := i.(int64), j.(int64)
        if asc {
            return i < j
        } else {
            return i > j
        }
    }

    return false
}
