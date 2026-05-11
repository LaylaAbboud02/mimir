package engine

import (
	"net/http"
	"testing"
)

func TestEvaluate_XPoweredBy(t *testing.T) {
	express := "Express"
	tests := []struct {
		name              string
		headers           http.Header
		wantIDs           []string
		wantObservedValue *string
	}{
		{
			name:    "without X-Powered-By",
			headers: http.Header{},
		},
		{
			name: "with X-Powered-By",
			headers: http.Header{
				"X-Powered-By": []string{"Express"},
			},
			wantIDs:           []string{"x-powered-by-disclosure"},
			wantObservedValue: &express,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Evaluate(Input{Headers: tt.headers})

			if len(got.Findings) != len(tt.wantIDs) {
				t.Fatalf("Findings: got %d, want %d", len(got.Findings), len(tt.wantIDs))
			}
			for i, id := range tt.wantIDs {
				if got.Findings[i].ID != id {
					t.Errorf("Findings[%d].ID: got %q, want %q", i, got.Findings[i].ID, id)
				}
			}
			if tt.wantObservedValue != nil {
				if got.Findings[0].ObservedValue == nil {
					t.Fatalf("Findings[0].ObservedValue: got nil, want non-nil")
				}
				if *got.Findings[0].ObservedValue != *tt.wantObservedValue {
					t.Errorf("Findings[0].ObservedValue: got %q, want %q", *got.Findings[0].ObservedValue, *tt.wantObservedValue)
				}
			}
		})
	}
}
