package cli

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

// Error kind constants for FetchError.Kind. These values are the stable
// identifiers used by the CLI output and the JSON error shape in §8.3.
const (
	KindDNS        = "dns"
	KindConnection = "connection"
	KindTLS        = "tls"
	KindTimeout    = "timeout"
	KindRedirect   = "redirect"
)

// FetchOptions controls the behavior of Fetch.
type FetchOptions struct {
	Timeout      time.Duration
	MaxRedirects int
	UserAgent    string
	// AllowHTTP permits a cross-protocol redirect from https to http.
	// It does not validate the scheme of the input URL — that is the
	// caller's responsibility.
	AllowHTTP bool
}

// FetchResult holds the response data needed for header evaluation and report
// assembly. The response body is never retained; only headers are kept.
type FetchResult struct {
	FinalURL      string
	RedirectCount int
	StatusCode    int
	Headers       http.Header
	Elapsed       time.Duration
}

// FetchError is a categorized operational error returned by Fetch.
type FetchError struct {
	Kind    string
	Message string
	Err     error
}

// Error implements the error interface.
func (e *FetchError) Error() string {
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

// Unwrap returns the underlying cause, enabling errors.Is and errors.As to
// traverse the chain.
func (e *FetchError) Unwrap() error {
	return e.Err
}

// Fetch performs an HTTP GET to rawURL and returns the response headers.
// The caller must supply a valid http:// or https:// URL; scheme validation
// is not performed here. AllowHTTP controls only cross-protocol redirect
// rejection (spec edge case E9), not the input URL scheme.
func Fetch(ctx context.Context, rawURL string, opts FetchOptions) (*FetchResult, error) {
	redirectCount := 0

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// If not redirects allowed
			if opts.MaxRedirects == 0 {
				return &FetchError{
					Kind:    KindRedirect,
					Message: "redirects disabled (max-redirects=0)",
				}
			}
			// If redirect count exceeded the specified limit
			if redirectCount >= opts.MaxRedirects {
				return &FetchError{
					Kind:    KindRedirect,
					Message: fmt.Sprintf("redirect limit of %d exceeded", opts.MaxRedirects),
				}
			}
			// Reject https to http downgrade unless --allow-http is set (spec E9).
			// Only applies when the immediately preceding request was https.
			if !opts.AllowHTTP && req.URL.Scheme == "http" &&
				len(via) > 0 && via[len(via)-1].URL.Scheme == "https" {
				return &FetchError{
					Kind:    KindRedirect,
					Message: "cross-protocol redirect from https to http rejected",
				}
			}
			redirectCount++
			return nil
		},
	}

	// Context controls cancellation/deadlines for the HTTP request.
	reqCtx := ctx

	// If the user provided a timeout, create a new context that will automatically cancel the request when the timeout is reached.
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		reqCtx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// Create a new HTTP GET request using the request context.
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, &FetchError{Kind: KindConnection, Message: err.Error(), Err: err}
	}

	// If the caller provided a custom User-Agent, add it to the request headers.
	if opts.UserAgent != "" {
		req.Header.Set("User-Agent", opts.UserAgent)
	}

	// Calculate how long the request took.
	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		return nil, classifyError(err)
	}

	// Drain and discard the body so the connection can be reused.
	// The tool only inspects headers; body contents are never needed.
	_, _ = io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	return &FetchResult{
		FinalURL:      resp.Request.URL.String(),
		RedirectCount: redirectCount,
		StatusCode:    resp.StatusCode,
		Headers:       resp.Header,
		Elapsed:       elapsed,
	}, nil
}

// classifyError maps raw http.Client errors onto the five FetchError kinds.
// Unclassifiable errors default to KindConnection.
func classifyError(err error) *FetchError {
	// FetchError from CheckRedirect arrives wrapped in *url.Error; unwrap it.
	var fe *FetchError
	if errors.As(err, &fe) {
		return fe
	}

	// Context deadline or transport-level deadline exceeded.
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, os.ErrDeadlineExceeded) {
		return &FetchError{Kind: KindTimeout, Message: err.Error(), Err: err}
	}

	// DNS resolution failure.
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return &FetchError{Kind: KindDNS, Message: err.Error(), Err: err}
	}

	// TLS record header error: e.g., plain-HTTP endpoint contacted over TLS.
	var tlsRecErr tls.RecordHeaderError
	if errors.As(err, &tlsRecErr) {
		return &FetchError{Kind: KindTLS, Message: err.Error(), Err: err}
	}

	// TLS certificate verification failure (Go 1.20+).
	var certErr *tls.CertificateVerificationError
	if errors.As(err, &certErr) {
		return &FetchError{Kind: KindTLS, Message: err.Error(), Err: err}
	}

	// net.OpError covers connection refused, transport timeouts, and other
	// network failures. Check Timeout() before defaulting to connection.
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		if opErr.Timeout() {
			return &FetchError{Kind: KindTimeout, Message: err.Error(), Err: err}
		}
		return &FetchError{Kind: KindConnection, Message: err.Error(), Err: err}
	}

	return &FetchError{Kind: KindConnection, Message: err.Error(), Err: err}
}
