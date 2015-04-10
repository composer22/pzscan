package scanner

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

const (
	testRootURL           = "example2.com"
	testMaxRunMin         = 1
	testMaxWorkers        = 4
	testValidPageTemplate = `
	<html>
	  <head>
	        <link rel="canonical" href="http://example2.com">
			<meta name="description" content="%s">
			<title>%s</title>
	  </head>
	  <body>
	   <h1>First</hi>
	   <img href="example.jpg" alt="valid test image">
		<a href="/page2">Link to page 2</a>
	  </body>
	</html>
`
	testInvalidPage = `
	<html>
	  <head>
			<meta name="description" content="too short a text">
			<title>Test Page 2</title>
	  </head>
	  <body>
	   <h1>First no error</hi>
	   <h1>Second causes error</hi>
	  <img href="example.jpg">
	  </body>
	</html>
`
)

func TestScanNew(t *testing.T) {
	t.Parallel()
	var testTimeEmpty time.Time
	s := New(testRootURL, testMaxRunMin, testMaxWorkers)

	if fmt.Sprint(reflect.TypeOf(s.mu)) != "sync.Mutex" {
		t.Errorf("sync.Mutex not initialized.")
	}
	if s.log == nil {
		t.Errorf("logger.Logger not initialized.")
	}
	testURL, _ := url.Parse(fmt.Sprintf("http://%s", testRootURL))

	if s.RootURL.String() != testURL.String() {
		t.Errorf("RootURL not initialized.")
	}
	if len(s.Tests) != 0 {
		t.Errorf("Tests not initialized.")
	}
	if s.MaxRunMin != testMaxRunMin {
		t.Errorf("MaxRunMin not initialized.")
	}
	if s.StartTime != testTimeEmpty {
		t.Errorf("StartTime not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(s.jobq)) != "chan *scanner.scanJob" {
		t.Errorf("chan *scanner.scanJob not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(s.doneCh)) != "chan *scanner.scanJob" {
		t.Errorf("chan *scanner.scanJob not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(s.wg)) != "sync.WaitGroup" {
		t.Errorf("sync.WaitGroup not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(s.stopOnce)) != "sync.Once" {
		t.Errorf("sync.Once not initialized.")
	}
}

func TestScanRun(t *testing.T) {
	handlerPage1 := func(w http.ResponseWriter, r *http.Request) {
		page := fmt.Sprintf(testValidPageTemplate,
			strings.Repeat("*", metaDescriptionMin),
			strings.Repeat("*", titleMin))
		io.WriteString(w, page)
	}
	handlerPage2 := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, testInvalidPage)
	}
	handlerImage := func(w http.ResponseWriter, r *http.Request) {
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerPage1)
	mux.HandleFunc("/page2", handlerPage2)
	mux.HandleFunc("/example.jpg", handlerImage)

	server := httptest.NewServer(mux)
	u, _ := url.Parse(fmt.Sprint(server.URL))
	scanner := New(u.Host, testMaxRunMin, testMaxWorkers)
	scanner.Run()
	server.Close()
	// TODO Continue testing.  We have images that aren't loading.
	//	fmt.Println(scanner.Tests)
}

func TestScanStop(t *testing.T) {
	t.Parallel()
	t.Skipf("TODO")
}

func TestScanDump(t *testing.T) {
	t.Parallel()
	t.Skipf("TODO")
}

func TestScanHandleSignals(t *testing.T) {
	t.Parallel()
	t.Skipf("Cannot test due to exit point.")
}

func TestScanEvaluate(t *testing.T) {
	t.Parallel()
	t.Skipf("TODO")
}
