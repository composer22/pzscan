package scanner

import (
	"net/url"

	"golang.org/x/net/html"
)

const (
	metaDescriptionMin = 131
	metaDescriptionMax = 154
	titleMin           = 57
	titleMax           = 68
)

// bodyAnalyzer is used to analyze a body of html text returned from a scan.
// Updates job statistics and finds addional URLs that need scanning.
type bodyAnalyzer struct {
	ScanJob *scanJob
}

// bodyAnalyzerNew regurns a new instance of a bodyAnalyzer
func bodyAnalyzerNew(j *scanJob) *bodyAnalyzer {
	return &bodyAnalyzer{ScanJob: j}
}

// AnalyzePage parses a page and analyzes it for SEO purposes, placing the results in stats.
func (a *bodyAnalyzer) analyzeBody() {
	p := html.NewTokenizer(a.ScanJob.Body)
	for {
		tt := p.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			token := p.Token()
			switch token.DataAtom.String() {
			case "a":
				a.anchorFound(&token)
			case "link":
				a.canonicalFound(&token)
			case "meta":
				a.metaDescriptions(&token)
			case "title":
				a.checkTitle(p)
			case "img":
				a.checkImages(&token)
			case "h1":
				a.checkH1()
			case "script":
				a.checkJS(&token)
			default:
			}
		default: // NOP
		}
	}
}

// anchorFound will scan an anchor element for new URLs
func (a *bodyAnalyzer) anchorFound(token *html.Token) {
	for _, attr := range token.Attr {
		if attr.Key == "href" {
			u, err := url.Parse(attr.Val)
			if err == nil {
				a.ScanJob.Children = append(a.ScanJob.Children, &scanJobChild{
					URL:     u,
					URLType: "html",
				})
			}
		}
	}
}

// canonicalFound will scan a link element for a rel="canonical" and set stats
// or if a CSS file was identified, record it as a new job.
func (a *bodyAnalyzer) canonicalFound(token *html.Token) {
	var csFound bool
	var href string
	for _, attr := range token.Attr {
		switch attr.Key {
		case "rel":
			switch attr.Val {
			case "canonical":
				a.ScanJob.Stat.Canonical = true // Canonical found
			case "stylesheet":
				csFound = true
			}
		case "href":
			href = attr.Val
		}
	}

	if csFound && href != "" {
		u, err := url.Parse(href)
		if err == nil {
			a.ScanJob.Children = append(a.ScanJob.Children, &scanJobChild{
				URL:     u,
				URLType: "css",
			})
		}
	}
}

// metaDescriptions will scan a meta element for description information and set stats
func (a *bodyAnalyzer) metaDescriptions(token *html.Token) {
	var descriptionFound bool
	var content string

	for _, attr := range token.Attr {
		switch attr.Key {
		case "name":
			if attr.Val == "description" {
				descriptionFound = true
			}
		case "content":
			content = attr.Val
		}
	}

	if descriptionFound {
		a.ScanJob.Stat.MetaCount++
		if len(content) < metaDescriptionMin || len(content) > metaDescriptionMax {
			a.ScanJob.Stat.MetaSizedErr = true
		}
	}
}

// checkTitle will scan a title element for content and set stats.
func (a *bodyAnalyzer) checkTitle(p *html.Tokenizer) {
	var title string

	// Try and get the text
	tokenType := p.Next()
	if tokenType == html.TextToken {
		token := p.Token()
		title = token.String()
	}

	a.ScanJob.Stat.TitleCount++
	if len(title) < titleMin || len(title) > titleMax {
		a.ScanJob.Stat.TitleSizedErr = true
	}
}

// checkImages will scan an img element for an alt tag and sets stats.
func (a *bodyAnalyzer) checkImages(token *html.Token) {
	var altFound bool

	for _, attr := range token.Attr {
		switch attr.Key {
		case "alt":
			if len(attr.Val) > 0 {
				altFound = true
			}
		case "src":
			u, err := url.Parse(attr.Val)
			if err == nil && u.Path != "" {
				a.ScanJob.Children = append(a.ScanJob.Children, &scanJobChild{
					URL:     u,
					URLType: "img",
				})
			}
		}
	}

	if !altFound {
		a.ScanJob.Stat.AltTagsErr = true // Valid alt not found for this image.
	}
}

// checkH1 will record h1 stats.
func (a *bodyAnalyzer) checkH1() {
	a.ScanJob.Stat.H1Count++
}

// checkJS will scan a script eelement for an src tag and sets stats.
func (a *bodyAnalyzer) checkJS(token *html.Token) {
	for _, attr := range token.Attr {
		if attr.Key == "src" {
			u, err := url.Parse(attr.Val)
			if err == nil && u.Path != "" {
				a.ScanJob.Children = append(a.ScanJob.Children, &scanJobChild{
					URL:     u,
					URLType: "js",
				})
			}
		}
	}
}
