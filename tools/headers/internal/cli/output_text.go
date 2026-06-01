// PLACEHOLDER OUTPUT — this format will be replaced in a later session with
// the advisor-voice text output described in the spec §7.1.
// Do not invest effort here.
package cli

import (
	"fmt"
	"io"

	"github.com/LaylaAbboud02/mimir/tools/headers/internal/engine"
)

// WriteTextOutput writes a minimal plaintext summary of the engine result and
// fetch metadata to w. The format is intentional placeholder output only.
func WriteTextOutput(w io.Writer, result engine.Result, fetchResult *FetchResult) error {
	_, err := fmt.Fprintf(w, "mimir-headers %s — %s (%d, %dms)\n\n",
		Version,
		fetchResult.FinalURL,
		fetchResult.StatusCode,
		fetchResult.Elapsed.Milliseconds(),
	)
	if err != nil {
		return err
	}

	for _, finding := range result.Findings {
		if finding.ObservedValue != nil {
			_, err = fmt.Fprintf(w, "[%s] %s  %s: %s\n",
				finding.Severity,
				finding.ID,
				finding.Header,
				*finding.ObservedValue,
			)
		} else {
			_, err = fmt.Fprintf(w, "[%s] %s  %s\n",
				finding.Severity,
				finding.ID,
				finding.Header,
			)
		}
		if err != nil {
			return err
		}
	}

	// Summary count at the bottom.
	_, err = fmt.Fprintf(w, "\n%d finding(s)\n", result.Summary.Total)
	return err
}