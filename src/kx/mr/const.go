package mr

const (
    SORT_BY_KEY   SortType = iota + 1
    SORT_BY_VALUE
    SORT_SECONDARY_KV // first key, then value
    SORT_SECONDARY_VK // first value, then key

    SORT_ORDER_ASC  = 1
    SORT_ORDER_DESC = 2

    KEY_SECONDARY_KV = "skv"
    KEY_SECONDARY_VK = "svk"
)
