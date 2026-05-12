package engine

import (
	"net/http"
	"strconv"
	"strings"
)

type CatalogEntry struct {
	ID          string   `json:"id"`
	Header      string   `json:"header"`
	Status      Status   `json:"status"`
	Severity    Severity `json:"severity"`
	Message     string   `json:"message"`
	Remediation string   `json:"remediation"`
	Check       func(http.Header) (matched bool, observed *string)
}

// parseMaxAge extracts the max-age directive value from a
// Strict-Transport-Security header value string. Returns the integer value and
// true if the directive is present and parses cleanly; returns 0 and false
// otherwise.
func parseMaxAge(headerVal string) (int64, bool) {
	for _, part := range strings.Split(headerVal, ";") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		if !strings.EqualFold(strings.TrimSpace(kv[0]), "max-age") {
			continue
		}
		n, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
		if err != nil {
			return 0, false
		}
		return n, true
	}
	return 0, false
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
	{
		ID:          "hsts-missing",
		Header:      "Strict-Transport-Security",
		Status:      StatusMissing,
		Severity:    SeverityMedium,
		Message:     "The Strict-Transport-Security header is not set on this response; browsers have no instruction to enforce HTTPS exclusively for this domain.",
		Remediation: "Add Strict-Transport-Security with a starting value of \"max-age=31536000; includeSubDomains\" and increase the max-age once the site is confirmed HTTPS-only.",
		Check: func(headers http.Header) (matched bool, observed *string) {
			if len(headers.Values("Strict-Transport-Security")) == 0 {
				return true, nil
			}
			return false, nil
		},
	},
	{
		ID:          "hsts-disabled",
		Header:      "Strict-Transport-Security",
		Status:      StatusWeak,
		Severity:    SeverityMedium,
		Message:     "The Strict-Transport-Security header is present but sets max-age=0, which actively revokes HSTS protection and instructs browsers to stop enforcing HTTPS for this domain.",
		Remediation: "Set max-age to a positive value such as \"max-age=31536000; includeSubDomains\" to enforce HTTPS consistently.",
		Check: func(headers http.Header) (matched bool, observed *string) {
			vals := headers.Values("Strict-Transport-Security")
			if len(vals) == 0 {
				return false, nil
			}
			joined := strings.Join(vals, ", ")
			age, ok := parseMaxAge(joined)
			if !ok || age != 0 {
				return false, nil
			}
			return true, &joined
		},
	},
	{
		ID:          "hsts-short-max-age",
		Header:      "Strict-Transport-Security",
		Status:      StatusWeak,
		Severity:    SeverityMedium,
		Message:     "The Strict-Transport-Security max-age is set but shorter than the recommended 180 days (15552000 seconds), reducing the window of protection against downgrade attacks.",
		Remediation: "Increase max-age to at least one year (31536000 seconds), for example \"max-age=31536000; includeSubDomains\".",
		Check: func(headers http.Header) (matched bool, observed *string) {
			vals := headers.Values("Strict-Transport-Security")
			if len(vals) == 0 {
				return false, nil
			}
			joined := strings.Join(vals, ", ")
			age, ok := parseMaxAge(joined)
			if !ok || age <= 0 || age >= 15552000 {
				return false, nil
			}
			return true, &joined
		},
	},
}
