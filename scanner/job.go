package scanner

import (
	"encoding/json"
	"io"
	"net/url"
)

// scanJobChild represents any URL found on an html page (ex html, img, css, js etc.
type scanJobChild struct {
	URL     *url.URL `json:"url"`     // The URL found.
	URLType string   `json:"urlType"` // The URL type ex: html, img, css, js etc.
}

// scanJob is a transport packet that represents URL that needs processing.
type scanJob struct {
	Stat     *Stats          `json:"stat"`     // Stats from the scan.
	Body     io.ReadCloser   `json:"body"`     // Body returned from the scan.
	Children []*scanJobChild `json:"children"` // Child URLs found on the page.
}

// scanJobNew is a factory for creating a new job instance.
func scanJobNew(u *url.URL, urlType string, parent *url.URL) *scanJob {
	return &scanJob{
		Stat:     StatsNew(u, urlType, parent),
		Children: []*scanJobChild{},
	}
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (s *scanJob) String() string {
	result, _ := json.Marshal(s)
	return string(result)
}
