package scraper

// ProtonScraper scrapes the Proton blog (proton.me/blog).
type ProtonScraper struct {
}

var protonSite = siteConfig{
	title: `article header h1`,
	body:  `div.post-content`,
	skip: []string{
		`.sr-only`,
		`div.text-center`,
		`p.has-text-align-center`,
	},
	links: `div.blog-item-title a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *ProtonScraper) ScrapeLinks(url string) (map[string]string, error) {
	return protonSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *ProtonScraper) ScrapeArticle(url string) (string, error) {
	return protonSite.article(url)
}

func (s *ProtonScraper) ScrapeTitle(url string) (string, error) {
	return protonSite.scrapeTitle(url)
}

func (s *ProtonScraper) ScrapeFilename(url string) (string, error) {
	return protonSite.filename(url), nil
}
