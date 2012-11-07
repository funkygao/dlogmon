package stats

type Counter map[string] int

func NewCounter(vals []string) Counter {
    c := make(Counter)
    for _, v := range vals {
        if _, found := c[v]; !found {
            c[v] = 1
        } else {
            c[v] ++
        }
    }
    return c
}
