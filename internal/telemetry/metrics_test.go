package telemetry

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	u, err := Read()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if u.CPU < 0 || u.CPU > 100 {
		t.Fatalf("cpu out of range: %v", u.CPU)
	}
	if u.Memory < 0 || u.Memory > 100 {
		t.Fatalf("mem out of range: %v", u.Memory)
	}
}

func TestReadProcess(t *testing.T) {
	u, err := ReadProcess(int32(os.Getpid()))
	if err != nil {
		t.Fatalf("read process: %v", err)
	}
	if u.CPU < 0 || u.CPU > 100 {
		t.Fatalf("cpu out of range: %v", u.CPU)
	}
	if u.Memory < 0 || u.Memory > 100 {
		t.Fatalf("mem out of range: %v", u.Memory)
	}
}
