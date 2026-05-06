package engine

// Holds count of findings (total and by severity)
type Summary struct {
	Total      int `json:"total"`
	BySeverity         map[Severity]int `json:"by_severity"`
}
