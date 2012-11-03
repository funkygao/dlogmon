package stream

const (
    EXEC_PIPE StreamMode = iota
    LZOP_FILE
    PLAIN_FILE

    FIRST_MODE = EXEC_PIPE
    LAST_MODE  = PLAIN_FILE
)
