// File: pkg/crawler/crawler_test.go

package crawler

import "testing"

// TestNew verifies that New returns a valid (non-nil) instance of Crawler.
func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Error("Expected New() to return a non-nil Crawler instance")
	}
}

// TestFetchURL verifies that FetchURL returns an empty string and nil error
// regardless of the input URL, as it is currently a stub.
func TestFetchURL(t *testing.T) {
	c := New()
	content, err := c.FetchURL("http://example.com")
	if err != nil {
		t.Errorf("Expected no error from FetchURL, got: %v", err)
	}
	if content != "" {
		t.Errorf("Expected empty content from FetchURL, got: %q", content)
	}
}
