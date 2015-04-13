package scanner

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

const (
	TestStatPrintExpResult = `{"url":{"Scheme":"http","Opaque":"","User":null,"Host":` +
		`"www.example.com","Path":"/faq","RawQuery":"","Fragment":""},"urlType":"html",` +
		`"parentURL":{"Scheme":"http","Opaque":"","User":null,"Host":"www.example.com",` +
		`"Path":"","RawQuery":"","Fragment":""},"startTime":"0001-01-01T00:00:00Z",` +
		`"endTime":"0001-01-01T00:00:00Z","canonical":false,"metaCount":0,"metaSizedErr":` +
		`false,"titleCount":0,"titleSizedErr":false,"altTagsErr":false,"h1Count":0,"status":0}`
)

func TestStatsNew(t *testing.T) {
	t.Parallel()
	ru, _ := url.Parse("http://www.example.com")
	u, _ := url.Parse("http://www.example.com/faq")
	stat := StatsNew(u, "html", ru)
	if stat.URL != u {
		t.Errorf("Invalid URL")
	}
	if stat.URLType != "html" {
		t.Errorf("Invalid URL type")
	}
	if stat.ParentURL != ru {
		t.Errorf("Invalid Parent URL")
	}
	if fmt.Sprint(reflect.TypeOf(stat.URL)) != "*url.URL" {
		t.Errorf("*url.URL not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(stat.URLType)) != "string" {
		t.Errorf("string not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(stat.ParentURL)) != "*url.URL" {
		t.Errorf("*url.URL not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(stat.StartTime)) != "time.Time" {
		t.Errorf("time.Time not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(stat.EndTime)) != "time.Time" {
		t.Errorf("time.Time not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(stat.Canonical)) != "bool" {
		t.Errorf("bool expected.")
	}
	if fmt.Sprint(reflect.TypeOf(stat.MetaCount)) != "int" {
		t.Errorf("int expected.")
	}
	if fmt.Sprint(reflect.TypeOf(stat.MetaSizedErr)) != "bool" {
		t.Errorf("bool expected.")
	}
	if fmt.Sprint(reflect.TypeOf(stat.TitleCount)) != "int" {
		t.Errorf("int expected.")
	}
	if fmt.Sprint(reflect.TypeOf(stat.TitleSizedErr)) != "bool" {
		t.Errorf("bool expected.")
	}
	if fmt.Sprint(reflect.TypeOf(stat.AltTagsErr)) != "bool" {
		t.Errorf("bool expected.")
	}
	if fmt.Sprint(reflect.TypeOf(stat.H1Count)) != "int" {
		t.Errorf("int expected.")
	}
	if fmt.Sprint(reflect.TypeOf(stat.StatusCode)) != "int" {
		t.Errorf("int expected.")
	}
}

func TestStatsPrint(t *testing.T) {
	t.Parallel()
	ru, _ := url.Parse("http://www.example.com")
	u, _ := url.Parse("http://www.example.com/faq")
	s := StatsNew(u, "html", ru)
	rslt := fmt.Sprint(s)
	if rslt != TestStatPrintExpResult {
		t.Errorf("Invalid Print of Stats. Expected: %s\n\nReceived: %s\n",
			TestStatPrintExpResult, rslt)
	}
}
