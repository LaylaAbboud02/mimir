package cli

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetch_SuccessfulFetch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "test-agent/1.0" {
			t.Errorf("User-Agent: got %q, want %q", r.Header.Get("User-Agent"), "test-agent/1.0")
		}
		w.Header().Set("X-Custom", "hello")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	result, err := Fetch(context.Background(), ts.URL, FetchOptions{
		Timeout:      5 * time.Second,
		MaxRedirects: 5,
		UserAgent:    "test-agent/1.0",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.StatusCode != http.StatusOK {
		t.Errorf("status: got %d, want 200", result.StatusCode)
	}
	if result.Headers.Get("X-Custom") != "hello" {
		t.Errorf("X-Custom: got %q, want %q", result.Headers.Get("X-Custom"), "hello")
	}
	if result.RedirectCount != 0 {
		t.Errorf("redirect count: got %d, want 0", result.RedirectCount)
	}
	if result.Elapsed < 0 {
		t.Errorf("elapsed is negative: %v", result.Elapsed)
	}
	if result.FinalURL == "" {
		t.Error("FinalURL must not be empty")
	}
}

func TestFetch_RedirectChainAllowed(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, ts.URL+"/final", http.StatusFound)
			return
		}
		w.Header().Set("X-Final", "yes")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	result, err := Fetch(context.Background(), ts.URL, FetchOptions{
		Timeout:      5 * time.Second,
		MaxRedirects: 5,
		UserAgent:    "test-agent/1.0",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RedirectCount != 1 {
		t.Errorf("redirect count: got %d, want 1", result.RedirectCount)
	}
	want := ts.URL + "/final"
	if result.FinalURL != want {
		t.Errorf("FinalURL: got %q, want %q", result.FinalURL, want)
	}
}

func TestFetch_RedirectChainExceedsLimit(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, ts.URL+"/loop", http.StatusFound)
	}))
	defer ts.Close()

	_, err := Fetch(context.Background(), ts.URL, FetchOptions{
		Timeout:      5 * time.Second,
		MaxRedirects: 3,
		UserAgent:    "test-agent/1.0",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var fe *FetchError
	if !errors.As(err, &fe) {
		t.Fatalf("expected *FetchError, got %T: %v", err, err)
	}
	if fe.Kind != KindRedirect {
		t.Errorf("error kind: got %q, want %q", fe.Kind, KindRedirect)
	}
}

// TestFetch_MaxRedirectsZero checks that MaxRedirects=0 causes any redirect
// to be treated as an error rather than following it. This is the chosen
// behavior: MaxRedirects=0 means "no redirects at all", returning a redirect
// error rather than the redirect response itself, keeping error semantics
// consistent with exceeding the redirect cap.
func TestFetch_MaxRedirectsZero(t *testing.T) {
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, ts.URL+"/other", http.StatusFound)
	}))
	defer ts.Close()

	_, err := Fetch(context.Background(), ts.URL, FetchOptions{
		Timeout:      5 * time.Second,
		MaxRedirects: 0,
		UserAgent:    "test-agent/1.0",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var fe *FetchError
	if !errors.As(err, &fe) {
		t.Fatalf("expected *FetchError, got %T: %v", err, err)
	}
	if fe.Kind != KindRedirect {
		t.Errorf("error kind: got %q, want %q", fe.Kind, KindRedirect)
	}
}

func TestFetch_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	_, err := Fetch(context.Background(), ts.URL, FetchOptions{
		Timeout:      50 * time.Millisecond,
		MaxRedirects: 5,
		UserAgent:    "test-agent/1.0",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var fe *FetchError
	if !errors.As(err, &fe) {
		t.Fatalf("expected *FetchError, got %T: %v", err, err)
	}
	if fe.Kind != KindTimeout {
		t.Errorf("error kind: got %q, want %q", fe.Kind, KindTimeout)
	}
}

// TestFetch_NonSuccessStatus verifies that a non-2xx response is not an error:
// Fetch succeeds and FetchResult.StatusCode reflects the server's status,
// because header evaluation is performed regardless of status code (spec E10).
func TestFetch_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Error-Header", "present")
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	result, err := Fetch(context.Background(), ts.URL, FetchOptions{
		Timeout:      5 * time.Second,
		MaxRedirects: 5,
		UserAgent:    "test-agent/1.0",
	})

	if err != nil {
		t.Fatalf("unexpected error for non-2xx response: %v", err)
	}
	if result.StatusCode != http.StatusInternalServerError {
		t.Errorf("status: got %d, want 500", result.StatusCode)
	}
	if result.Headers.Get("X-Error-Header") != "present" {
		t.Error("expected X-Error-Header to be present in FetchResult.Headers")
	}
}
