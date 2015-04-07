package scanner

import (
	"encoding/json"
	"net/url"
	"time"
)

// Stats is a construct that hold information on the scanning of a URL.
type Stats struct {
	URL        *url.URL  `json:"url"`       // The URL we scanned.
	StartTime  time.Time `json:"startTime"` // The start time of the scan.
	EndTime    time.Time `json:"endTime"`   // The end time of the scan.
	StatusCode int       `json:"status"`    // The status code we returned from the scan.
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (s *Stats) String() string {
	result, _ := json.Marshal(s)
	return string(result)
}
