package engine

import "net/http"

// Input to the engine
type Input struct {
	Headers http.Header `json:"headers"`
}

// Result of the engine evaluation
// Includes:
// - Findings: List of findings (headers with issues)
// - Summary: Summary of the findings (total and by severity)
type Result struct {
	Findings []Finding `json:"findings"`
	Summary  Summary   `json:"summary"`
}

// Evaluates the input and returns the result
func Evaluate(input Input) Result {
	return Result{
		Findings: []Finding{},
		Summary:  Summary{},
	}
}



