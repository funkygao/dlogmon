package stream

var streamMode = []string{
    "ExecPipe",
    "LzopFile",
    "PlainText",
}

func (this StreamMode) String() string {
    return streamMode[this]
}

func (this StreamMode) Valids() (modes []string) {
    modes = make([]string, LAST_MODE+1)
    for i := FIRST_MODE; i <= LAST_MODE; i++ {
        modes[i] = streamMode[i]
    }

    return
}
