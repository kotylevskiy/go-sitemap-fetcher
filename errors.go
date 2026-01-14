package gositemapfetcher

import (
	"fmt"
	"net/url"
)

// ErrNilYield indicates a nil yield callback was provided.
type ErrNilYield struct{}

func (e *ErrNilYield) Error() string {
	return "yield callback is nil"
}

// ErrInvalidURL indicates the input website or sitemap URL is invalid.
type ErrInvalidURL struct {
	URL string
	Err error
}

func (e *ErrInvalidURL) Error() string {
	if e.URL == "" {
		return fmt.Sprintf("invalid URL: %v", e.Err)
	}
	return fmt.Sprintf("invalid URL %q: %v", e.URL, e.Err)
}

func (e *ErrInvalidURL) Unwrap() error {
	return e.Err
}

// ErrNoSitemaps indicates that no sitemap URLs were discovered.
type ErrNoSitemaps struct {
	URL *url.URL
}

func (e *ErrNoSitemaps) Error() string {
	if e.URL == nil {
		return "no sitemaps discovered"
	}
	return fmt.Sprintf("no sitemaps discovered for %s", e.URL)
}

// ErrHTTPStatus indicates an unexpected HTTP status while fetching a sitemap.
type ErrHTTPStatus struct {
	URL        *url.URL
	StatusCode int
	Status     string
}

func (e *ErrHTTPStatus) Error() string {
	if e.URL == nil {
		return fmt.Sprintf("unexpected HTTP status %d", e.StatusCode)
	}
	return fmt.Sprintf("unexpected HTTP status %d for %s", e.StatusCode, e.URL)
}

// ErrSitemapParse indicates a failure while parsing sitemap XML.
type ErrSitemapParse struct {
	URL *url.URL
	Err error
}

func (e *ErrSitemapParse) Error() string {
	if e.URL == nil {
		return fmt.Sprintf("sitemap parse failed: %v", e.Err)
	}
	return fmt.Sprintf("sitemap parse failed for %s: %v", e.URL, e.Err)
}

func (e *ErrSitemapParse) Unwrap() error {
	return e.Err
}

// ErrMaxDepth indicates the sitemap index depth limit was exceeded.
type ErrMaxDepth struct {
	MaxDepth int
	URL      *url.URL
}

func (e *ErrMaxDepth) Error() string {
	if e.URL == nil {
		return fmt.Sprintf("max depth %d exceeded", e.MaxDepth)
	}
	return fmt.Sprintf("max depth %d exceeded at %s", e.MaxDepth, e.URL)
}

// ErrMaxSitemaps indicates the sitemap count limit was exceeded.
type ErrMaxSitemaps struct {
	MaxSitemaps int
}

func (e *ErrMaxSitemaps) Error() string {
	return fmt.Sprintf("max sitemaps %d exceeded", e.MaxSitemaps)
}

// ErrMaxURLs indicates the URL count limit was exceeded.
type ErrMaxURLs struct {
	MaxURLs int
}

func (e *ErrMaxURLs) Error() string {
	return fmt.Sprintf("max URLs %d exceeded", e.MaxURLs)
}

// ErrYield wraps a failure returned by the yield callback.
type ErrYield struct {
	Err error
}

func (e *ErrYield) Error() string {
	return fmt.Sprintf("yield callback failed: %v", e.Err)
}

func (e *ErrYield) Unwrap() error {
	return e.Err
}
