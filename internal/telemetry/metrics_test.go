package telemetry

import (
	"context"
	"os"
	"testing"
	"time"
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

func TestPollProcess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch, err := PollProcess(ctx, int32(os.Getpid()), 10*time.Millisecond)
	if err != nil {
		t.Fatalf("poll: %v", err)
	}
	select {
	case u := <-ch:
		if u.CPU < 0 || u.CPU > 100 || u.Memory < 0 || u.Memory > 100 {
			t.Fatalf("invalid usage: %#v", u)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("timeout")
	}
	cancel()
	for range ch {
	}
}
