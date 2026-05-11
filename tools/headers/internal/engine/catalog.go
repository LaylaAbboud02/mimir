package engine

import "net/http"

type CatalogEntry struct {
	ID          string   `json:"id"`
	Header      string   `json:"header"`
	Status      Status   `json:"status"`
	Severity    Severity `json:"severity"`
	Message     string   `json:"message"`
	Remediation string   `json:"remediation"`
	Check       func(http.Header) (matched bool, observed *string)
}

var Catalog = []CatalogEntry{
	{
		ID:          "x-powered-by-disclosure",
		Header:      "X-Powered-By",
		Status:      StatusPresent,
		Severity:    SeverityInfo,
		Message:     "The X-Powered-By header exposes the application framework or stack.",
		Remediation: "Remove the X-Powered-By header from HTTP responses or disable it in your application server configuration.",
		Check: func(headers http.Header) (matched bool, observed *string) {
			val := headers.Get("X-Powered-By")
			if val == "" {
				return false, nil
			}
			return true, &val
		},
	},
}
