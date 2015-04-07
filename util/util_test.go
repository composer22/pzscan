package util

import (
	"strings"
	"testing"
)

func TestGetLinks(t *testing.T) {
	body := strings.NewReader(`some useless text and an <a href="www.plumdistrict.com/faq.html">www.foo.com</>`)
	result := GetLinks(body)
	for _, v := range result {
		t.Logf(v.Host)
		t.Logf(v.Path)
	}

}
