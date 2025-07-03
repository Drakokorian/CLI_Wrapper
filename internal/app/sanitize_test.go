package app

import "testing"

func TestSanitizePrompt(t *testing.T) {
	cases := []struct{ in, out string }{
		{"hello", "hello"},
		{"-danger", "_danger"},
		{"\ntrim\r", "trim"},
	}
	for _, c := range cases {
		got := SanitizePrompt(c.in)
		if got != c.out {
			t.Errorf("expected %q got %q", c.out, got)
		}
	}
}

func TestSanitizeModel(t *testing.T) {
	if _, err := SanitizeModel("../bad"); err == nil {
		t.Error("expected error")
	}
	if _, err := SanitizeModel("bad model"); err == nil {
		t.Error("expected error")
	}
	cases := []string{
		"good-model_1",
		"gpt-3.5-turbo",
		"model.version_2",
	}
	for _, c := range cases {
		if m, err := SanitizeModel(c); err != nil || m != c {
			t.Errorf("unexpected result %v %v", m, err)
		}
	}
}
