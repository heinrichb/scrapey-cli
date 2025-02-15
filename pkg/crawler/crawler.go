// File: pkg/crawler/crawler.go

package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/*
Crawler is responsible for fetching HTML content from URLs.

Usage:

	Create an instance of Crawler using New() and then call FetchURL
	to retrieve the HTML content from a specified URL.

Notes:
  - This implementation is currently a stub.
  - Future enhancements may include handling HTTP GET requests,
    concurrency, rate-limiting, timeouts, retries, and robust error handling.
*/
type Crawler struct {
	// Fields for storing configuration or concurrency settings can be added here.
}

/*
New returns a new instance of Crawler.

Usage:

	c := New()
*/
func New() *Crawler {
	return &Crawler{}
}

/*
FetchURL retrieves the HTML content from the specified URL.

Parameters:
  - url: A string representing the URL to fetch.

Returns:
  - A string containing the HTML content (if successful) or an empty string.
  - An error if the fetch operation fails.

Usage:

	content, err := c.FetchURL("http://example.com")
	if err != nil {
	    // Handle error.
	}

Notes:
  - This function is currently a stub and returns an empty string with a nil error.
  - Future implementations will include actual HTTP request handling.
*/
func (c *Crawler) FetchURL(url string) (string, error) {
	// Stub: return placeholder HTML or empty string for now.
	client := &http.Client{
		Timeout: 10 * time.Second, // Set timeout
	}

	// jsonData := `{"key":"value"}`

	// Create a custom request
	// req, err := http.NewRequest("Post", url, bytes.NewBuffer([]byte(jsonData)))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return "", err
	}
	// req.Header.Set("Content-Type", "application/json") // Set headers

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return "", err
	}
	return string(body), nil
}
