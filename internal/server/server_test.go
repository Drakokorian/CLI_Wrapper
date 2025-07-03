package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/net/websocket"

	"cli-wrapper/internal/app"
	"cli-wrapper/internal/history"
	"cli-wrapper/internal/logging"
)

func TestEndpoints(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "state"), 0o755)
	logger, _ := logging.NewWithPath("info", filepath.Join(dir, "log.txt"))
	cfg := &app.Config{Concurrency: 1, Theme: "light", CPUThreshold: 50, MemoryThreshold: 50, PollInterval: 2, LogLevel: "info", LogPath: filepath.Join(dir, "log.txt"), WorkingDir: dir}
	hist, _ := history.New(dir)
	defer hist.Close()
	mgr := app.NewSessionManager(logger, 1, cfg, hist)
	defer mgr.Close()

	srv := New(mgr, logger, dir, cfg, hist)
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

	if err := hist.AddMetric(history.Metric{ID: "sid", CPU: 1, Memory: 2}); err != nil {
		t.Fatalf("add metric: %v", err)
	}
	resp, err = http.Get(ts.URL + "/resource/session/sid")
	if err != nil {
		t.Fatalf("session metrics: %v", err)
	}
	var metrics []history.Metric
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil || len(metrics) != 1 {
		t.Fatalf("decode metrics: %v len=%d", err, len(metrics))
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

	// config endpoints
	wd := filepath.Join(dir, "work")
	if err := os.Mkdir(wd, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	resp, err = http.Get(ts.URL + "/config")
	if err != nil {
		t.Fatalf("config get: %v", err)
	}
	var cfgResp map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&cfgResp); err != nil {
		t.Fatalf("decode config: %v", err)
	}
	if cfgResp["concurrency"].(float64) != 1 || cfgResp["logLevel"].(string) == "" {
		t.Fatalf("expected concurrency 1")
	}

	postBody := map[string]any{"concurrency": 2, "workingDir": wd, "logLevel": "debug", "logPath": filepath.Join(dir, "new.log")}
	body, _ := json.Marshal(postBody)
	resp, err = http.Post(ts.URL+"/config", "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("post config: %v status %d", err, resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&cfgResp); err != nil {
		t.Fatalf("decode post config: %v", err)
	}
	if cfgResp["concurrency"].(float64) != 2 || cfgResp["logLevel"].(string) != "debug" {
		t.Fatalf("expected concurrency 2")
	}

	// history endpoints
	data := `[{"id":"1","prompt":"hi","response":"hello","model":"openai","success":true}]`
	resp, err = http.Post(ts.URL+"/history", "application/json", bytes.NewBufferString(data))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("import history: %v status %d", err, resp.StatusCode)
	}

	resp, err = http.Get(ts.URL + "/history")
	if err != nil {
		t.Fatalf("get history: %v", err)
	}
	var recs []history.Record
	if err := json.NewDecoder(resp.Body).Decode(&recs); err != nil || len(recs) != 1 {
		t.Fatalf("decode history: %v len=%d", err, len(recs))
	}

	resp, err = http.Get(ts.URL + "/history/search?q=hello")
	if err != nil {
		t.Fatalf("search history: %v", err)
	}
	var recs2 []history.Record
	if err := json.NewDecoder(resp.Body).Decode(&recs2); err != nil || len(recs2) == 0 {
		t.Fatalf("search decode: %v len=%d", err, len(recs2))
	}

	resp, err = http.Get(ts.URL + "/history/export")
	if err != nil {
		t.Fatalf("export history: %v", err)
	}
	var exp []history.Record
	if err := json.NewDecoder(resp.Body).Decode(&exp); err != nil || len(exp) != 1 {
		t.Fatalf("export decode: %v len=%d", err, len(exp))
	}

	// billing endpoint using stubbed CLI
	script := filepath.Join(dir, "openai")
	content := "#!/bin/sh\necho '{\"total_granted\":1}'"
	if err := os.WriteFile(script, []byte(content), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", dir+string(os.PathListSeparator)+oldPath)
	defer os.Setenv("PATH", oldPath)

	resp, err = http.Get(ts.URL + "/billing")
	if err != nil {
		t.Fatalf("billing: %v", err)
	}
	var billResp struct {
		Tool  string           `json:"tool"`
		URL   string           `json:"url"`
		Usage app.BillingUsage `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&billResp); err != nil {
		t.Fatalf("decode billing: %v", err)
	}
	if billResp.Tool != "openai" || billResp.URL == "" || billResp.Usage.TotalGranted != 1 {
		t.Fatalf("unexpected billing response %#v", billResp)
	}
}

func TestStreamChars(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "state"), 0o755)
	logger, _ := logging.NewWithPath("info", filepath.Join(dir, "log.txt"))
	cfg := &app.Config{Concurrency: 1, Theme: "light", CPUThreshold: 50, MemoryThreshold: 50, PollInterval: 2, LogLevel: "info", LogPath: filepath.Join(dir, "log.txt"), WorkingDir: dir}
	hist, _ := history.New(dir)
	defer hist.Close()
	mgr := app.NewSessionManager(logger, 1, cfg, hist)
	defer mgr.Close()

	script := filepath.Join(dir, "echoer")
	if err := os.WriteFile(script, []byte("#!/bin/sh\necho hi"), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", dir+string(os.PathListSeparator)+oldPath)
	defer os.Setenv("PATH", oldPath)

	srv := New(mgr, logger, dir, cfg, hist)
	ts := httptest.NewServer(srv.mux)
	defer ts.Close()

	id, err := mgr.AddSession("echoer", "", []string{""})
	if err != nil {
		t.Fatalf("add session: %v", err)
	}

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/streamchars?id=" + id
	ws, err := websocket.Dial(wsURL, "", ts.URL)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer ws.Close()

	var out bytes.Buffer
	for i := 0; i < 3; i++ {
		var msg string
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			t.Fatalf("receive: %v", err)
		}
		out.WriteString(msg)
	}
	if out.String() != "hi\n" {
		t.Fatalf("got %q", out.String())
	}
}
