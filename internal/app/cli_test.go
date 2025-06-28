package app

import "testing"

func TestDetectCLIToolNone(t *testing.T) {
	if _, err := DetectCLITool(); err == nil {
		t.Fatal("expected error")
	}
}
