package scanner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"
)

const (
	testBodyNegativeContent = `some useless text without any valid content.`
	testBodyTemplateSource  = `some useless text and <%s %s %s> and some other text.`
	testBodyTemplateTag     = `Some more useless text with a <%s>%s</%s> stuck inside.`
	testBodyAnchorSource    = `<p>some useless text</p> and an <a href="www.example2.com/faq.html">www.foo.com</>`
	testBodyAnchorResult    = "www.example2.com/faq.html"
)

var (
	testURLRoot, _ = url.Parse("http://example.com")

	testBodyCanonical = []struct {
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
	}{
		{"meta", `name="description"`, "content", strings.Repeat("*", metaDescriptionMin), 1, false,
			"Valid low size description should not have been flagged as error."},
		{"meta", `name="description"`, "content", strings.Repeat("*", metaDescriptionMax), 1, false,
			"Valid high size description should not have been flagged as error."},
		{"meta", `name="description"`, "content", strings.Repeat("*", metaDescriptionMin-1), 1, true,
			"Invalid low size description should have been flagged as error."},
		{"meta", `name="description"`, "content", strings.Repeat("*", metaDescriptionMax+1), 1, true,
			"Invalid high size description should have been flagged as error."},
		{"Xeta", `name="description"`, "content", strings.Repeat("*", metaDescriptionMin), 0, false,
			"Invalid tag should have flagged a count error."},
		{"meta", `Xame="description"`, "content", strings.Repeat("*", metaDescriptionMin), 0, false,
			"Invalid attr1 should have flagged a count error."},
		{"meta", `name="Xescription"`, "content", strings.Repeat("*", metaDescriptionMin), 0, false,
			"Invalid attr1 value should have flagged a count error."},
		{"meta", `name="description"`, "Kontent", strings.Repeat("*", metaDescriptionMin), 1, true,
			"Invalid attr2 should have flagged a count error."},
	}

	testBodyTitle = []struct {
		tag            string
		text           string
		expectedSize   int
		expectedResult bool
		message        string
	}{
		{"Title", strings.Repeat("*", titleMin), 1, false, "Valid small title should have passed."},
		{"Title", strings.Repeat("*", titleMax), 1, false, "Valid large title should have passed."},
		{"Title", strings.Repeat("*", titleMin-1), 1, true, "Invalid low size should have been found."},
		{"Title", "", 1, true, "Invalid missing text should have been found."},
		{"Title", strings.Repeat("*", titleMax+1), 1, true, "Invalid high size should have been found."},
		{"Xitle", strings.Repeat("*", titleMin), 0, false, "Missing title should not have triggered err."},
	}

	testBodyImg = []struct {
		tag            string
		attr1          string
		attr2          string
		expectedResult bool
		expectedFind   int
		message        string
	}{
		{"img", `alt="some text"`, `src="smiley.gif"`, false, 1, "Valid image should not report error."},
		{"Xmg", `alt="some text"`, `src="smiley.gif"`, false, 0, "Invalid tag should have been reported."},
		{"img", `Xlt="some text"`, `src="smiley.gif"`, true, 1, "Invalid attr should have been reported."},
		{"img", `alt=""`, `src="smiley.gif"`, true, 1, "Missing alt text should have been reported."},
		{"img", `alt="some text"`, `Xrc="smiley.gif"`, false, 0, "Missing src url should be child."},
		{"img", `alt="some text"`, `src=""`, false, 0, "Missing src url should be child."},
	}

	testBodyH1 = []struct {
		text           string
		expectedResult int
		message        string
	}{
		{fmt.Sprintf(testBodyTemplateTag, "h1", "", "h1"), 1, "h1 element should have been found."},
		{fmt.Sprintf(testBodyTemplateTag, "X1", "", "X1"), 0, "h1 element should not have been found."},
		{strings.Repeat(fmt.Sprintf(testBodyTemplateTag, "h1", "", "h1"), 2), 2, "Multiple h1 elements should have been found."},
		{strings.Repeat(fmt.Sprintf(testBodyTemplateTag, "X1", "", "X1"), 2), 0, "Multiple h1 elements should not have been found."},
	}
)

func TestBodyAnalyzerAnchor(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	for _, tc := range testBodyDescription {
		content := fmt.Sprintf(`%s="%s"`, tc.attr2, tc.content)
		source := fmt.Sprintf(testBodyTemplateSource, tc.tag, tc.attr1, content)
		job := scanJobNew(testURLRoot, "html", nil)
		job.Body = ioutil.NopCloser(bytes.NewBufferString(source))
		a := bodyAnalyzerNew(job)
		a.analyzeBody()
		if a.ScanJob.Stat.MetaCount != tc.expectedSize ||
			a.ScanJob.Stat.MetaSizedErr != tc.expectedResult {
			t.Errorf(tc.message)
		}
	}
}

func TestBodyAnalyzerTitle(t *testing.T) {
	t.Parallel()
	for i, tc := range testBodyTitle {
		source := fmt.Sprintf(testBodyTemplateTag, tc.tag, tc.text, tc.tag)
		job := scanJobNew(testURLRoot, "html", nil)
		job.Body = ioutil.NopCloser(bytes.NewBufferString(source))
		a := bodyAnalyzerNew(job)
		a.analyzeBody()
		if a.ScanJob.Stat.TitleCount != tc.expectedSize ||
			a.ScanJob.Stat.TitleSizedErr != tc.expectedResult {
			t.Log(source)
			t.Logf("%d stats: %d %d %t ", i, len(source), a.ScanJob.Stat.TitleCount, a.ScanJob.Stat.TitleSizedErr)
			t.Errorf(tc.message)
		}
	}
}

func TestBodyAnalyzerImg(t *testing.T) {
	t.Parallel()
	for _, tc := range testBodyImg {
		source := fmt.Sprintf(testBodyTemplateSource, tc.tag, tc.attr1, tc.attr2)
		job := scanJobNew(testURLRoot, "html", nil)
		job.Body = ioutil.NopCloser(bytes.NewBufferString(source))
		a := bodyAnalyzerNew(job)
		a.analyzeBody()
		if a.ScanJob.Stat.AltTagsErr != tc.expectedResult ||
			len(a.ScanJob.Children) != tc.expectedFind {
			t.Errorf(tc.message)
		}
	}
}

func TestBodyH1(t *testing.T) {
	t.Parallel()
	for _, tc := range testBodyH1 {
		job := scanJobNew(testURLRoot, "html", nil)
		job.Body = ioutil.NopCloser(bytes.NewBufferString(tc.text))
		a := bodyAnalyzerNew(job)
		a.analyzeBody()
		if a.ScanJob.Stat.H1Count != tc.expectedResult {
			t.Errorf(tc.message)
		}
	}
}
