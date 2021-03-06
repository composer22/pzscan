package scanner

import (
	"encoding/json"
	"net/url"
	"time"
)

// Stats is a construct that hold information on the scanning of a URL.
type Stats struct {
	URL           *url.URL  `json:"url"`           // The URL we scanned.
	URLType       string    `json:"urlType"`       // The type of url ex: html, img, css, js etc..
	ParentURL     *url.URL  `json:"parentURL"`     // The parent where this was located.
	StartTime     time.Time `json:"startTime"`     // The start time of the scan.
	EndTime       time.Time `json:"endTime"`       // The end time of the scan.
	Canonical     bool      `json:"canonical"`     // Did this page contain a canonical link?
	MetaCount     int       `json:"metaCount"`     // Does meta description exist on the page?
	MetaSizedErr  bool      `json:"metaSizedErr"`  // Are meta descriptions the proper size?
	TitleCount    int       `json:"titleCount"`    // Does title exist on the page?
	TitleSizedErr bool      `json:"titleSizedErr"` // Does the title meet size criteria?
	AltTagsErr    bool      `json:"altTagsErr"`    // Did alt tags exist for all images on this page?
	H1Count       int       `json:"h1Count"`       // Does an h1 tag exist on the page and is it unique?
	StatusCode    int       `json:"status"`        // The status code we returned from the scan.
}

// StatsNew is a factory for creating a new Stats instance.
func StatsNew(u *url.URL, ut string, p *url.URL) *Stats {
	return &Stats{
		URL:       u,
		URLType:   ut,
		ParentURL: p,
	}
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (s *Stats) String() string {
	j, _ := json.Marshal(s)
	return string(j)
}
