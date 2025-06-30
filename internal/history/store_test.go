package history

import (
	"encoding/json"
	"testing"
)

func TestStoreCRUD(t *testing.T) {
	dir := t.TempDir()
	s, err := New(dir)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	defer s.Close()

	r := Record{ID: "1", Prompt: "hi", Response: "hello", Model: "openai", Success: true}
	if err := s.Add(r); err != nil {
		t.Fatalf("add: %v", err)
	}
	all, err := s.All()
	if err != nil || len(all) != 1 {
		t.Fatalf("all: %v len=%d", err, len(all))
	}

	res, err := s.Search("hello", 10)
	if err != nil || len(res) == 0 {
		t.Fatalf("search: %v len=%d", err, len(res))
	}

	data, err := s.Export()
	if err != nil {
		t.Fatalf("export: %v", err)
	}

	// import into new store
	s2, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("new2: %v", err)
	}
	defer s2.Close()
	if err := s2.Import(data); err != nil {
		t.Fatalf("import: %v", err)
	}
	all2, _ := s2.All()
	if len(all2) != 1 {
		b, _ := json.Marshal(all2)
		t.Fatalf("imported len=%d data=%s", len(all2), string(b))
	}
}

func TestMetrics(t *testing.T) {
	dir := t.TempDir()
	s, err := New(dir)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	defer s.Close()
	m := Metric{ID: "1", CPU: 10, Memory: 20}
	if err := s.AddMetric(m); err != nil {
		t.Fatalf("add metric: %v", err)
	}
	metrics, err := s.Metrics("1")
	if err != nil || len(metrics) != 1 {
		t.Fatalf("metrics: %v len=%d", err, len(metrics))
	}
	if metrics[0].CPU != 10 || metrics[0].Memory != 20 {
		t.Fatalf("unexpected metric %#v", metrics[0])
	}
}
