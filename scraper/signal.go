package scraper

// SignalScraper scrapes the Signal blog (signal.org/blog).
type SignalScraper struct {
}

var signalSite = siteConfig{
	title: `div.blog-post-header h1`,
	body:  `div.blog-post-content div.column`,
	skip: []string{
		`div.social-sharing`,
	},
	links: `div.columns h3 a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *SignalScraper) ScrapeLinks(url string) (map[string]string, error) {
	return signalSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *SignalScraper) ScrapeArticle(url string) (string, error) {
	return signalSite.article(url)
}

func (s *SignalScraper) ScrapeTitle(url string) (string, error) {
	return signalSite.scrapeTitle(url)
}

func (s *SignalScraper) ScrapeFilename(url string) (string, error) {
	return signalSite.filename(url), nil
}
