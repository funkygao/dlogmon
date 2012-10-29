package dlog

import "testing"

func TestMaxConcurrentWorkersGreaterThanThousand(t *testing.T) {
	if MAX_CONCURRENT_WORKERS < 1000 {
		t.Error("<1000")
	}
}
