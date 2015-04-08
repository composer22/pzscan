package scanner

import (
	"encoding/json"
	"net/url"
	"time"
)

// Stats is a construct that hold information on the scanning of a URL.
type Stats struct {
	URL           *url.URL  `json:"url"`        // The URL we scanned.
	URLType       string    `json:"urlType"`    // The type of url ex: html, img, css, js etc..
	ParentURL     *url.URL  `json:"parentURL"`  // The parent where this was located.
	StartTime     time.Time `json:"startTime"`  // The start time of the scan.
	EndTime       time.Time `json:"endTime"`    // The end time of the scan.
	Canonical     bool      `json:"canonical"`  // Did this page contain a canonical link?
	MetaCount     int       `json:"metaExist"`  // Does meta description exist on the page?
	MetaSizedErr  bool      `json:"metaSized"`  // Are meta descriptions the proper size?
	TitleCount    int       `json:"titleExist"` // Does title exist on the page?
	TitleSizedErr bool      `json:"titleSized"` // Does the title meet size criteria?
	AltTagsErr    bool      `json:"altTags"`    // Did alt tags exist for all images on this page?
	H1Count       int       `json:"h1Exist"`    // Does an h1 tag exist on the page and is it unique?
	StatusCode    int       `json:"status"`     // The status code we returned from the scan.
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (s *Stats) String() string {
	result, _ := json.Marshal(s)
	return string(result)
}
