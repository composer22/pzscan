package scanner

import (
	"fmt"
	"net/url"
	"testing"
)

const (
	TestStatPrintExpResult = `{"url":{"Scheme":"http","Opaque":"","User":null,"Host":` +
		`"www.example.com","Path":"/faq","RawQuery":"","Fragment":""},"urlType":"html",` +
		`"parentURL":{"Scheme":"http","Opaque":"","User":null,"Host":"www.example.com",` +
		`"Path":"","RawQuery":"","Fragment":""},"startTime":"0001-01-01T00:00:00Z",` +
		`"endTime":"0001-01-01T00:00:00Z","canonical":false,"metaExist":0,"metaSizedErr":` +
		`false,"titleExist":0,"titleSizedErr":false,"altTagsErr":false,"h1Exist":0,"status":0}`
)

func TestStatsNew(t *testing.T) {
	t.Parallel()
	rootURL, _ := url.Parse("http://www.example.com")
	url, _ := url.Parse("http://www.example.com/faq")
	stat := StatsNew(url, "html", rootURL)
	if stat.URL != url {
		t.Errorf("Invalid URL")
	}
	if stat.URLType != "html" {
		t.Errorf("Invalid URL type")
	}
	if stat.ParentURL != rootURL {
		t.Errorf("Invalid Parent URL")
	}
}

func TestStatsPrint(t *testing.T) {
	t.Parallel()
	rootURL, _ := url.Parse("http://www.example.com")
	url, _ := url.Parse("http://www.example.com/faq")
	stat := StatsNew(url, "html", rootURL)
	if fmt.Sprint(stat) != TestStatPrintExpResult {
		t.Errorf("Invalid Print of Stats")
	}
}
