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
	if m, err := SanitizeModel("good-model_1"); err != nil || m != "good-model_1" {
		t.Errorf("unexpected result %v %v", m, err)
	}
}
