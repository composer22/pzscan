package scanner

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"
)

const (
	testBodyNegativeContent = `some useless text without any valid content.`
	testBodyTemplateSource  = `some useless text and <%s %s %s> and some other text.`

	testBodyAnchorSource = `<p>some useless text</p> and an <a href="www.example2.com/faq.html">www.foo.com</>`
	testBodyAnchorResult = "www.example2.com/faq.html"
)

var (
	testURLRoot, _ = url.Parse("http://example.com")

	descMin                 = make([]byte, metaDescriptionMin)
	descMax                 = make([]byte, metaDescriptionMax)
	descMinErr              = make([]byte, metaDescriptionMin-1)
	descMaxErr              = make([]byte, metaDescriptionMax+1)
	testBodyMetaValidLow    = ""
	testBodyMetaValidHigh   = ""
	testBodyMetaInvalidLow  = ""
	testBodyMetaInvalidHigh = ""
	testBodyCanonical       = []struct {
		tag            string
		attr           string
		expectedResult bool
		message        string
	}{
		{"link", `rel="canonical"`, true, "Link and attr should have been found"},
		{"kink", `rel="canonical"`, false, "Link should not have been found."},
		{"link", `rex="canonical"`, false, "canonical attr should not have been found."},
		{"link", `rel="Kanonical"`, false, "canonical attr should not have been found."},
	}

	testBodyDescription = []struct {
		tag            string
		attr1          string
		attr2          string
		content        string
		expectedSize   int
		expectedResult bool
		message        string
	}{}
)

func init() {
	rand.Read(descMin)
	rand.Read(descMax)
	rand.Read(descMinErr)
	rand.Read(descMaxErr)
	testBodyMetaValidLow = fmt.Sprintf("%X", descMin)
	testBodyMetaValidHigh = fmt.Sprintf("%X", descMax)
	testBodyMetaInvalidLow = fmt.Sprintf("%X", descMinErr)
	testBodyMetaInvalidHigh = fmt.Sprintf("%X", descMaxErr)

	testBodyDescription = []struct {
		tag            string
		attr1          string
		attr2          string
		content        string
		expectedSize   int
		expectedResult bool
		message        string
	}{
		{"meta", `name="description"`, "content", testBodyMetaValidLow, 1, true,
			"Valid low size description should not have been flagged as error."},
		{"meta", `name="description"`, "content", testBodyMetaValidHigh, 1, true,
			"Valid high size description should not have been flagged as error."},
		{"meta", `name="description"`, "content", testBodyMetaValidLow, 0, false,
			"Invalid low size description should have been flagged as error."},
		{"meta", `name="description"`, "content", testBodyMetaValidHigh, 0, false,
			"Invalid high size description should have been flagged as error."},
		{"Xeta", `name="description"`, "content", testBodyMetaValidLow, 0, true,
			"Invalid tag should have flagged a count error."},
		{"meta", `Xame="description"`, "content", testBodyMetaValidLow, 0, true,
			"Invalid attr1 should have flagged a count error."},
		{"meta", `name="Xescription"`, "content", testBodyMetaValidLow, 0, true,
			"Invalid attr1 value should have flagged a count error."},
		{"meta", `name="description"`, "content", testBodyMetaValidLow, 0, true,
			"Invalid attr2 should have flagged a count error."},
	}
}
func TestBodyAnalyzerAnchor(t *testing.T) {
	job := scanJobNew(testURLRoot, "html", nil)
	job.Body = ioutil.NopCloser(bytes.NewBufferString(testBodyAnchorSource))
	a := bodyAnalyzerNew(job)
	a.analyzeBody()
	found := false
	for _, c := range a.ScanJob.Children {
		if c.URL.String() == testBodyAnchorResult {
			found = true
		}
	}
	if !found {
		t.Errorf("Ancor tag should have been found.")
	}

	job = scanJobNew(testURLRoot, "html", nil)
	job.Body = ioutil.NopCloser(bytes.NewBufferString(testBodyNegativeContent))
	a = bodyAnalyzerNew(job)
	a.analyzeBody()
	if len(a.ScanJob.Children) > 0 {
		t.Errorf("Ancor tag should not have been found.")
	}
}

func TestBodyAnalyzerCanonical(t *testing.T) {
	for _, tc := range testBodyCanonical {
		source := fmt.Sprintf(testBodyTemplateSource, tc.tag, tc.attr, "")
		job := scanJobNew(testURLRoot, "html", nil)
		job.Body = ioutil.NopCloser(bytes.NewBufferString(source))
		a := bodyAnalyzerNew(job)
		a.analyzeBody()
		if a.ScanJob.Stat.Canonical != tc.expectedResult {
			t.Errorf(tc.message)
		}
	}
}

func TestBodyAnalyzerMeta(t *testing.T) {
	for _, tc := range testBodyDescription {
		content := fmt.Sprintf(`%s="%s"`, tc.attr2, tc.content)
		source := fmt.Sprintf(testBodyTemplateSource, tc.tag, tc.attr1, content)
		job := scanJobNew(testURLRoot, "html", nil)
		t.Logf("%d", len(content))
		job.Body = ioutil.NopCloser(bytes.NewBufferString(source))
		a := bodyAnalyzerNew(job)
		a.analyzeBody()
		if a.ScanJob.Stat.MetaCount != tc.expectedSize ||
			a.ScanJob.Stat.MetaSizedErr != tc.expectedResult {
			t.Errorf(tc.message)
		}
	}
}
