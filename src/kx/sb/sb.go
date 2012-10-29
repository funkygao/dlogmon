package sb

import "bytes"

type StringBuilder struct {
	*bytes.Buffer
}

func NewStringBuilder(s string) *StringBuilder {
	return &StringBuilder{bytes.NewBuffer([]byte(s))}
}

func (this *StringBuilder) Append(s string) {
	this.WriteString(s)
}
