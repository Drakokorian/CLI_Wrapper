package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"cli-wrapper/internal/app"
	"cli-wrapper/internal/logging"
)

func TestEndpoints(t *testing.T) {
	logger, _ := logging.NewWithPath(filepath.Join(t.TempDir(), "log.txt"))
	mgr := app.NewSessionManager(t.TempDir(), logger, 1)
	defer mgr.Close()

	srv := New(mgr, logger)
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
}
