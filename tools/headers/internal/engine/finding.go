package engine

// Severity of the finding
type Severity string

const (
	SeverityHigh   Severity = "high"
	SeverityMedium Severity = "medium"
	SeverityLow    Severity = "low"
	SeverityInfo   Severity = "info"
)

// Type of issue with the header (e.g. missing, misconfigured, etc.)
type Status string

const (
	StatusPresent   Status = "present"
	StatusMissing   Status = "missing"
	StatusMalformed Status = "malformed"
	StatusWeak      Status = "weak"
)

// Represents a single header finding
// Includes:
// - ID: Unique identifier for the finding
// - Status: Whether the header is present, missing, or malformed
// - Header: The header name
// - Severity: The severity of the finding
// - ObservedValue: The value of the header observed in the response (can be null)
// - Message: A human-readable message describing the finding
// - Remediation: A remediation for the finding
type Finding struct {
	ID             string   `json:"id"`
	Status         Status   `json:"status"`
	Header         string   `json:"header"`
	Severity       Severity `json:"severity"`
	ObservedValue  *string  `json:"observed_value"`
	Message        string   `json:"message"`
	Remediation string   `json:"remediation"`
}
