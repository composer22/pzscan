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
	ru, _ := url.Parse("http://www.example.com")
	u, _ := url.Parse("http://www.example.com/faq")
	j := scanJobNew(u, "html", ru)
	if j.Stat.URL != u {
		t.Errorf("Invalid URL.")
	}
	if j.Stat.URLType != "html" {
		t.Errorf("Invalid URL type.")
	}
	if j.Stat.ParentURL != ru {
		t.Errorf("Invalid Parent URL.")
	}
	if len(j.Children) != 0 {
		t.Errorf("Invalid children.")
	}
	if fmt.Sprint(reflect.TypeOf(j.Stat)) != "*scanner.Stats" {
		t.Errorf("*scanner.Stats not initialized.")
	}
	if j.Body != nil {
		t.Errorf("Body not initialized correcly.")
	}

	if fmt.Sprint(reflect.TypeOf(j.Children)) != "[]*scanner.scanJobChild" {
		t.Errorf("[]*scanJobChild not initialized.")
	}
}

func TestScanJobPrint(t *testing.T) {
	t.Parallel()
	ru, _ := url.Parse("http://www.example.com")
	u, _ := url.Parse("http://www.example.com/faq")
	j := scanJobNew(u, "html", ru)
	if fmt.Sprint(j) != TestScanJobPrintExpResult {
		t.Errorf("Invalid Print of scan job.")
	}
}
