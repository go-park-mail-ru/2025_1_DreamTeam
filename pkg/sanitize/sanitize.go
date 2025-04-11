package sanitize

import "github.com/microcosm-cc/bluemonday"

func Sanitize(input string) string {
	p := bluemonday.NewPolicy()

	p.AllowStandardURLs()
	p.AllowElements("p")
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowElements("ul", "ol", "li")
	p.AllowElements("strong")
	p.AllowElements("table", "tr", "td", "th")
	return p.Sanitize(input)
}
