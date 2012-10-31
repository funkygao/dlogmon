package stream

var streamMode = []string{
    "ExecPipe",
    "LzopFile",
    "PlainText",
}

func (this StreamMode) String() string {
    return streamMode[this]
}
