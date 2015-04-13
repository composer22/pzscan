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
	   <h1>Only One</hi>
	   <img src="example.jpg" alt="valid test image">
		<a href="/page2">Link to page 2</a>
	  </body>
	</html>
`
	testInvalidPage = `
	<html>
	  <head>
       <!-- missing canonical -->
	    <link rel="stylesheet" href="example.css">
		<meta name="description" content="short description">
		<meta name="description" content="too many?">
		 <script src="example.js"></script>
		<title>Short Title</title>
	  </head>
	  <body>
	   <h1>First no error</hi>
	   <h1>Second too many</hi>
	  <img src="example.jpg" xaltm="missing alt">
	  </body>
	</html>
`
)

var (
	testScannerLookup = map[string]struct {
		urlType       string
		canonical     bool
		metaCount     int
		metaSizedErr  bool
		titleCount    int
		titleSizedErr bool
		altTagsErr    bool
		h1Count       int
		status        int
	}{
		"html":  {"html", true, 1, false, 1, false, false, 1, 200},
		"img":   {"img", false, 0, false, 0, false, false, 0, 200},
		"html2": {"html", false, 2, true, 1, true, true, 2, 200},
		"css":   {"css", false, 0, false, 0, false, false, 0, 200},
		"js":    {"js", false, 0, false, 0, false, false, 0, 200},
	}
)

func TestScanNew(t *testing.T) {
	t.Parallel()
	var tTimeEmpty time.Time
	s := New(testRootURL, testMaxRunMin, testMaxWorkers)

	tURL, _ := url.Parse(fmt.Sprintf("http://%s", testRootURL))
	if s.RootURL.String() != tURL.String() {
		t.Errorf("RootURL not initialized.")
	}
	if len(s.Tests) != 0 {
		t.Errorf("Tests not initialized.")
	}
	if s.MaxRunMin != testMaxRunMin {
		t.Errorf("MaxRunMin not initialized.")
	}
	if s.MaxWorkers != testMaxWorkers {
		t.Errorf("MaxWorkers not initialized.")
	}
	if s.StartTime != tTimeEmpty {
		t.Errorf("StartTime not initialized.")
	}

	if s.ExpireTime != tTimeEmpty {
		t.Errorf("ExpireTime not initialized.")
	}

	if s.EndTime != tTimeEmpty {
		t.Errorf("EndTime not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(s.mu)) != "sync.Mutex" {
		t.Errorf("sync.Mutex not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(s.wg)) != "sync.WaitGroup" {
		t.Errorf("sync.WaitGroup not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(s.stopOnce)) != "sync.Once" {
		t.Errorf("sync.Once not initialized.")
	}
	if s.log == nil {
		t.Errorf("logger.Logger not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(s.jobq)) != "chan *scanner.scanJob" {
		t.Errorf("chan *scanner.scanJob not initialized.")
	}
	if fmt.Sprint(reflect.TypeOf(s.doneCh)) != "chan *scanner.scanJob" {
		t.Errorf("chan *scanner.scanJob not initialized.")
	}
}

func TestScanRun(t *testing.T) {
	hPage1 := func(w http.ResponseWriter, r *http.Request) {
		pg := fmt.Sprintf(testValidPageTemplate,
			strings.Repeat("*", metaDescriptionMin),
			strings.Repeat("*", titleMin))
		io.WriteString(w, pg)
	}
	hPage2 := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, testInvalidPage)
	}
	hImage := func(w http.ResponseWriter, r *http.Request) {
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", hPage1)
	mux.HandleFunc("/page2", hPage2)
	mux.HandleFunc("/example.jpg", hImage)

	srvr := httptest.NewServer(mux)
	u, _ := url.Parse(fmt.Sprint(srvr.URL))
	scnr := New(u.Host, testMaxRunMin, testMaxWorkers)
	scnr.Run()
	srvr.Close()

	for _, chdrn := range scnr.Tests {
		for _, stat := range chdrn {
			k := stat.URLType
			if k == "html" && stat.URL.Path == "/page2" {
				k = "html2"
			}
			expResult := testScannerLookup[k]
			if stat.URLType != expResult.urlType {
				t.Errorf("Invalid URL type returned.")
			}
			if stat.Canonical != expResult.canonical {
				t.Errorf("Canonical tested incorrectly.")
			}
			if stat.MetaCount != expResult.metaCount {
				t.Errorf("MetaCount tested incorrectly.")
			}
			if stat.MetaSizedErr != expResult.metaSizedErr {
				t.Errorf("MetaSizedErr tested incorrectly.")
			}
			if stat.TitleCount != expResult.titleCount {
				t.Errorf("TitleCount tested incorrectly.")
			}
			if stat.TitleSizedErr != expResult.titleSizedErr {
				t.Errorf("TitleSizedErr tested incorrectly.")
			}
			if stat.AltTagsErr != expResult.altTagsErr {
				t.Errorf("AltTagsErr tested incorrectly.")
			}
			if stat.H1Count != expResult.h1Count {
				t.Errorf("H1Count tested incorrectly.")
			}
			if stat.StatusCode != http.StatusOK {
				t.Errorf("Status Code returned tested incorrectly.")
			}
		}
	}
	scnr.Stop()
}

func TestScanVersionAndExit(t *testing.T) {
	t.Parallel()
	t.Skip("Cannot test due to Exit()")
}

func TestScanStop(t *testing.T) {
	t.Parallel()
	t.Skip("Covered by TestScanRun")
}

func TestScanHandleSignals(t *testing.T) {
	t.Parallel()
	t.Skip("Cannot test due to Exit()")
}

func TestScanEvaluate(t *testing.T) {
	t.Parallel()
	t.Skip("Covered by TestScanRun")
}
