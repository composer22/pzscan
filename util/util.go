package util

import (
	"io"
	"net/url"

	"golang.org/x/net/html"
)

// GetLinks parses a page for anchor tags and grabs the URLS.
func GetLinks(body io.Reader) []*url.URL {
	urls := []*url.URL{} // Results.
	p := html.NewTokenizer(body)
	for {
		tType := p.Next()
		switch tType {
		case html.ErrorToken:
			return urls
		case html.StartTagToken:
			token := p.Token()
			// Anchor found?
			if token.DataAtom.String() == "a" {
				// Check each attribute
				for _, attr := range token.Attr {
					// If href, then get value and store it.
					if attr.Key == "href" {
						u, err := url.Parse(attr.Val)
						if err == nil {
							urls = append(urls, u)
						}
					}
				}
			}
		default: // NOP
		}
	}
}
