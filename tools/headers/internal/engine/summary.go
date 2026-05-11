package engine

// Holds count of findings (total and by severity)
type Summary struct {
	Total      int `json:"total"`
	BySeverity map[Severity]int `json:"by_severity"`
}

func computeSummary(findings []Finding) Summary {
	by := make(map[Severity]int)
	by[SeverityHigh] = 0
	by[SeverityMedium] = 0
	by[SeverityLow] = 0
	by[SeverityInfo] = 0
	for _, f := range findings {
		by[f.Severity]++
	}
	return Summary{
		Total:      len(findings),
		BySeverity: by,
	}
}
