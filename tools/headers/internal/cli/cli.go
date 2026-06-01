package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/LaylaAbboud02/mimir/tools/headers/internal/engine"
)

// Version is the current tool version. Build pipelines may override this via
// -ldflags "-X github.com/.../cli.Version=1.2.3" when cutting a release.
const Version = "0.1.0-dev"

// Run is the program entry point. It parses args, fetches the target URL,
// evaluates the response headers, and writes the result to stdout.
// It returns an exit code: 0 on success, 2 on operational error.
func Run(args []string, stdout, stderr io.Writer) int {
	options, err := ParseFlags(args, stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 2
	}
	if options.ShowVersion {
		fmt.Fprintf(stdout, "mimir-headers v%s\n", Version)
		return 0
	}

	validURL, err := ValidateURL(options.URL, options.AllowHTTP)
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 2
	}

	userAgent := options.UserAgent
	if userAgent == "" {
		userAgent = "mimir-headers/" + Version
	}

	fetchOptions := FetchOptions{
		Timeout: options.Timeout,
		MaxRedirects: options.MaxRedirects,
		AllowHTTP: options.AllowHTTP,
		UserAgent: userAgent,
	}

	fetchResult, err := Fetch(context.Background(), validURL, fetchOptions)
	if err != nil {
		var fetchErr *FetchError
		if errors.As(err, &fetchErr) {
			fmt.Fprintf(stderr, "error: %s: %s\n", fetchErr.Kind, fetchErr.Message)
		} else {
			fmt.Fprintf(stderr, "error: %s\n", err)
		}
		return 2
	}

	engineInput := engine.Input{
		Headers: fetchResult.Headers,
	}

	result := engine.Evaluate(engineInput)

	if err := WriteTextOutput(stdout, result, fetchResult); err != nil {
		fmt.Fprintf(stderr, "error: writing output: %s\n", err)
		return 2
	}

	return 0
}