package scanner

import (
	"encoding/json"
	"net/url"
)

// scanJob is a transport packet that represents URL that needs processing.
type scanJob struct {
	Stat    *Stats     `json:"stat"`     // Stats from the scan.
	Childen []*url.URL `json:"children"` // Child URLs found on the page.
}

// scanJobNew is a factory for creating a new job instance.
func scanJobNew(u *url.URL) *scanJob {
	return &scanJob{
		Stat:    &Stats{URL: u},
		Childen: []*url.URL{},
	}
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (s *scanJob) String() string {
	result, _ := json.Marshal(s)
	return string(result)
}
