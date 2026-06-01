package cli

import (
	"bytes"
	"testing"
)

func TestRun_InvalidURL(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := Run([]string{"example.com"}, &stdout, &stderr)
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
	if stderr.Len() == 0 {
		t.Error("expected error message on stderr, got nothing")
	}
	if stdout.Len() != 0 {
		t.Errorf("expected stdout to be empty on error, got %q", stdout.String())
	}
}