package sb

import "testing"

func TestStringBuilder(t *testing.T) {
    sb := NewStringBuilder("we")
    sb.Append("are")
    sb.Append("here")
    if sb.String() != "wearehere" {
        t.Error("fail")
    }
}
