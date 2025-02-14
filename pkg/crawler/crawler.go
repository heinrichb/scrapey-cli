package crawler

// Crawler is responsible for fetching HTML content from URLs.
// Will possibly handle concurrency, rate-limits, etc. later.
type Crawler struct {
	// We might store config references or concurrency settings here.
}

// New returns a new instance of Crawler.
func New() *Crawler {
	return &Crawler{}
}

// FetchURL fetches the contents of a given URL.
// We'll eventually handle HTTP GET requests, timeouts, retries, etc.
func (c *Crawler) FetchURL(url string) (string, error) {
	// Stub: return placeholder HTML or empty string for now
	return "", nil
}
