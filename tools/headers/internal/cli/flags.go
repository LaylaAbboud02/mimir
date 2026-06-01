package cli

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"time"
)

// Options holds the parsed CLI flags and the positional URL argument.
type Options struct {
	Timeout      time.Duration
	MaxRedirects int
	AllowHTTP    bool
	UserAgent    string
	ShowVersion  bool
	URL          string
}

// ParseFlags parses args into an Options value. It writes usage and error text
// to stderr. Returns flag.ErrHelp when -h/--help is passed.
func ParseFlags(args []string, stderr io.Writer) (*Options, error) {
	flags := flag.NewFlagSet("mimir-headers", flag.ContinueOnError)
	flags.SetOutput(stderr)

	options := &Options{}

	flags.DurationVar(&options.Timeout, "timeout", 10*time.Second, "Timeout for the HTTP request")
	flags.IntVar(&options.MaxRedirects, "max-redirects", 5, "Maximum number of redirects to follow")
	flags.BoolVar(&options.AllowHTTP, "allow-http", false, "Allow HTTP URLs (not recommended)")
	flags.StringVar(&options.UserAgent, "user-agent", "", "Custom User-Agent string")
	flags.BoolVar(&options.ShowVersion, "version", false, "Show version information")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	if options.ShowVersion {
		return options, nil
	}

	if flags.NArg() != 1 {
		return nil, fmt.Errorf("expected exactly one URL argument; got %d\nUsage: mimir-headers [flags] <url>", flags.NArg())
	}

	options.URL = flags.Arg(0)

	return options, nil

}

// ValidateURL checks that rawURL is a usable http/https URL and returns its
// canonical form. It performs no network I/O.
func ValidateURL(rawURL string, allowHTTP bool) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Must check empty scheme before the http/https check — an empty scheme
	// also satisfies != "http" && != "https", giving a confusing error message.
	if parsedURL.Scheme == "" {
		return "", fmt.Errorf("url must include scheme (e.g. https://%s)", rawURL)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("url scheme must be http or https, got %q", parsedURL.Scheme)
	}

	if parsedURL.Scheme == "http" && !allowHTTP {
		return "", fmt.Errorf("refusing to fetch plain http URL; pass --allow-http to override")
	}

	if parsedURL.Host == "" {
		return "", fmt.Errorf("url has no host")
	}

	return parsedURL.String(), nil
}