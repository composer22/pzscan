package scanner

import (
	"io"
	"net/url"

	"golang.org/x/net/html"
)

const (
	metaDescriptionMin = 130
	metaDescriptionMax = 155
	titleMin           = 57
	titleMax           = 68
)

// AnalyzePage parses a page and analyzes it for SEO purposes, placing the results in stats.
func AnalyzePage(body io.Reader, stats *Stats) []*url.URL {
	urls := []*url.URL{} // ChildrenResults.
	p := html.NewTokenizer(body)
	for {
		tType := p.Next()
		switch tType {
		case html.ErrorToken:
			return urls
		case html.StartTagToken:
			token := p.Token()
			switch token.DataAtom.String() {
			case "a":
				anchorFound(&token, urls)
			case "link":
				canonicalFound(&token, stats)
			case "meta":
				metaDescriptions(&token, stats)
			case "title":
				checkTitle(p, stats)
			case "img":
				checkImages(&token, stats)
			case "h1":
				checkH1(stats)
			default:
			}
		default: // NOP
		}
	}
}

// anchorFound will scan an anchor element for href link children
func anchorFound(token *html.Token, urls []*url.URL) {
	// Check each attribute
	for _, attr := range token.Attr {
		switch attr.Key {
		case "href":
			u, err := url.Parse(attr.Val)
			if err == nil {
				urls = append(urls, u)
			}
		}
	}
}

// canonicalFound will scan a link element for a rel="canonical"
func canonicalFound(token *html.Token, stats *Stats) {
	// Check each attribute
	for _, attr := range token.Attr {
		switch attr.Key {
		case "rel":
			if attr.Val == "canonical" {
				stats.Canonical = true // Canonical found
			}
		}
	}
}

// metaDescriptions will scan a meta element for description information and set stats
func metaDescriptions(token *html.Token, stats *Stats) {
	var descriptionFound bool = false
	var content string = ""

	// Validate the name and content.
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

	// Set stats.
	if descriptionFound {
		stats.MetaCount++
		if len(content) <= metaDescriptionMin || len(content) >= metaDescriptionMax {
			stats.MetaSizedErr = true
		}
	}
}

// checkTitle will scan a title element for content and set stats.
func checkTitle(p *html.Tokenizer, stats *Stats) {
	textNode := p.Next() // get the text
	title := textNode.String()

	stats.TitleCount++

	// Check size
	if len(title) < titleMin || len(title) > titleMax {
		stats.TitleSizedErr = true
	}
}

// checkImages will scan an img element for an alt tag.
func checkImages(token *html.Token, stats *Stats) {
	altFound := false

	// Check each attribute
	for _, attr := range token.Attr {
		switch attr.Key {
		case "alt":
			if len(attr.Val) > 0 {
				altFound = true
			}

		}
	}
	if !altFound {
		stats.AltTagsErr = true // Valid alt not found for this image.
	}
}

// checkH1 will record h1 stats.
func checkH1(stats *Stats) {
	stats.H1Count++
}
