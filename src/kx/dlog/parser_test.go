package dlog

import (
	"testing"
)

func TestTrimAllRune(t *testing.T) {
	got := trimAllRune("acd,-mmq", []rune{'a', 'm', ','})
	expected := "cd-q"
	if got != expected {
		t.Error(got + " != " + expected)
	}
}
