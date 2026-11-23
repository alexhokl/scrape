package scraper

// JessFrazelleScraper scrapes Jessie Frazelle's blog (blog.jessfraz.com).
type JessFrazelleScraper struct {
}

var jessFrazelleSite = siteConfig{
	title: `article h1`,
	body:  `article`,
	skip: []string{
		`h1`,
		`div.time`,
		`aside.toc`,
		`a.footnote-return`,
	},
	links: `article:not(.publication) h3 a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *JessFrazelleScraper) ScrapeLinks(url string) (map[string]string, error) {
	return jessFrazelleSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *JessFrazelleScraper) ScrapeArticle(url string) (string, error) {
	return jessFrazelleSite.article(url)
}

func (s *JessFrazelleScraper) ScrapeTitle(url string) (string, error) {
	return jessFrazelleSite.scrapeTitle(url)
}

func (s *JessFrazelleScraper) ScrapeFilename(url string) (string, error) {
	return jessFrazelleSite.filename(url), nil
}
