package scanner

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

const (
	TestScanJobPrintExpResult = `{"stat":{"url":{"Scheme":"http","Opaque":"",` +
		`"User":null,"Host":"www.example.com","Path":"/faq","RawQuery":"","Fragment":""},` +
		`"urlType":"html","parentURL":{"Scheme":"http","Opaque":"","User":null,"Host":` +
		`"www.example.com","Path":"","RawQuery":"","Fragment":""},` +
		`"startTime":"0001-01-01T00:00:00Z","endTime":"0001-01-01T00:00:00Z",` +
		`"canonical":false,"metaCount":0,"metaSizedErr":false,"titleCount":0,"titleSizedErr":` +
		`false,"altTagsErr":false,"h1Count":0,"status":0},"body":null,"children":[]}`
)

func TestScanJobNew(t *testing.T) {
	t.Parallel()
	rootURL, _ := url.Parse("http://www.example.com")
	url, _ := url.Parse("http://www.example.com/faq")
	job := scanJobNew(url, "html", rootURL)
	if job.Stat.URL != url {
		t.Errorf("Invalid URL")
	}
	if job.Stat.URLType != "html" {
		t.Errorf("Invalid URL type")
	}
	if job.Stat.ParentURL != rootURL {
		t.Errorf("Invalid Parent URL")
	}
	if len(job.Children) != 0 {
		t.Errorf("Invalid children")
	}
	if fmt.Sprint(reflect.TypeOf(job.Stat)) != "*scanner.Stats" {
		t.Errorf("*scanner.Stats not initialized.")
	}
	if job.Body != nil {
		t.Errorf("Body not initialized correcly.")
	}

	if fmt.Sprint(reflect.TypeOf(job.Children)) != "[]*scanner.scanJobChild" {
		t.Errorf("[]*scanJobChild not initialized.")
	}

}

func TestScanJobPrint(t *testing.T) {
	t.Parallel()
	rootURL, _ := url.Parse("http://www.example.com")
	url, _ := url.Parse("http://www.example.com/faq")
	job := scanJobNew(url, "html", rootURL)
	if fmt.Sprint(job) != TestScanJobPrintExpResult {
		t.Errorf("Invalid Print of scan job")
	}
}
