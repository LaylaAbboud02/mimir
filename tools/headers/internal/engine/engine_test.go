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
			headers: http.Header{"Strict-Transport-Security": []string{"max-age=31536000"}},
		},
		{
			name: "with X-Powered-By",
			headers: http.Header{
				"X-Powered-By":              []string{"Express"},
				"Strict-Transport-Security": []string{"max-age=31536000"},
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

func TestEvaluate_HSTS(t *testing.T) {
	ptr := func(s string) *string { return &s }
	tests := []struct {
		name              string
		headers           http.Header
		wantIDs           []string
		wantObservedValue *string
	}{
		{
			name:    "missing",
			headers: http.Header{},
			wantIDs: []string{"hsts-missing"},
		},
		{
			name:    "long max-age",
			headers: http.Header{"Strict-Transport-Security": []string{"max-age=31536000"}},
		},
		{
			name:    "long max-age with includeSubDomains",
			headers: http.Header{"Strict-Transport-Security": []string{"max-age=31536000; includeSubDomains"}},
		},
		{
			name:              "disabled",
			headers:           http.Header{"Strict-Transport-Security": []string{"max-age=0"}},
			wantIDs:           []string{"hsts-disabled"},
			wantObservedValue: ptr("max-age=0"),
		},
		{
			name:              "short max-age",
			headers:           http.Header{"Strict-Transport-Security": []string{"max-age=86400"}},
			wantIDs:           []string{"hsts-short-max-age"},
			wantObservedValue: ptr("max-age=86400"),
		},
		{
			name:    "exact boundary",
			headers: http.Header{"Strict-Transport-Security": []string{"max-age=15552000"}},
		},
		{
			name:              "one second under boundary",
			headers:           http.Header{"Strict-Transport-Security": []string{"max-age=15551999"}},
			wantIDs:           []string{"hsts-short-max-age"},
			wantObservedValue: ptr("max-age=15551999"),
		},
		{
			name:    "malformed max-age",
			headers: http.Header{"Strict-Transport-Security": []string{"max-age=abc"}},
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
			if len(tt.wantIDs) > 0 {
				if tt.wantObservedValue != nil {
					if got.Findings[0].ObservedValue == nil {
						t.Fatalf("Findings[0].ObservedValue: got nil, want %q", *tt.wantObservedValue)
					}
					if *got.Findings[0].ObservedValue != *tt.wantObservedValue {
						t.Errorf("Findings[0].ObservedValue: got %q, want %q", *got.Findings[0].ObservedValue, *tt.wantObservedValue)
					}
				} else {
					if got.Findings[0].ObservedValue != nil {
						t.Errorf("Findings[0].ObservedValue: got %q, want nil", *got.Findings[0].ObservedValue)
					}
				}
			}
		})
	}
}
