package sanitize

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

func Sanitize(input string) string {
	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()
	p.AllowElements("p")
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowElements("ul", "ol", "li")
	p.AllowElements("strong")
	p.AllowElements("table", "tr", "td", "th")
	p.AllowElements("pre", "code")
	p.AllowAttrs("class").Matching(regexp.MustCompile(`^language-go$`)).OnElements("code")
	return p.Sanitize(input)
}
