package engine

import (
	"net/http"
	"testing"
)

// Test that when no headers are provided, the result is empty
func TestEvaluate_EmptyHeaders_ReturnsEmptyResult(t *testing.T) {
	got := Evaluate(Input{
		Headers: http.Header{},
	})

	if len(got.Findings) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}