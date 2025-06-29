package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"cli-wrapper/internal/app"
	"cli-wrapper/internal/logging"
)

func TestEndpoints(t *testing.T) {
	dir := t.TempDir()
	logger, _ := logging.NewWithPath(filepath.Join(dir, "log.txt"))
	cfg := &app.Config{Concurrency: 1, Theme: "light", CPUThreshold: 50, MemoryThreshold: 50, PollInterval: 2}
	mgr := app.NewSessionManager(dir, logger, 1, cfg)
	defer mgr.Close()

	srv := New(mgr, logger, dir, cfg)
	ts := httptest.NewServer(srv.mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/sessions")
	if err != nil {
		t.Fatalf("sessions: %v", err)
	}
	var ids []string
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		t.Fatalf("decode sessions: %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("expected 0 sessions")
	}

	resp, err = http.Get(ts.URL + "/resource")
	if err != nil {
		t.Fatalf("resource: %v", err)
	}
	var usage map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
		t.Fatalf("decode resource: %v", err)
	}
	if usage["cpu"] < 0 || usage["memory"] < 0 {
		t.Fatalf("unexpected usage: %v", usage)
	}

	resp, err = http.Get(ts.URL + "/theme")
	if err != nil {
		t.Fatalf("theme: %v", err)
	}
	var themeResp map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&themeResp); err != nil {
		t.Fatalf("decode theme: %v", err)
	}
	if themeResp["theme"] != "light" {
		t.Fatalf("got %s want light", themeResp["theme"])
	}

	buf := bytes.NewBufferString(`{"theme":"dark"}`)
	resp, err = http.Post(ts.URL+"/theme", "application/json", buf)
	if err != nil {
		t.Fatalf("post theme: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&themeResp); err != nil {
		t.Fatalf("decode post theme: %v", err)
	}
	if themeResp["theme"] != "dark" {
		t.Fatalf("got %s want dark", themeResp["theme"])
	}
}
