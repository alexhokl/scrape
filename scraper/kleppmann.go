package scraper

// KleppmannScraper scrapes Martin Kleppmann's blog (martin.kleppmann.com).
type KleppmannScraper struct {
}

var kleppmannSite = siteConfig{
	title: `#content h1`,
	body:  `#content`,
	skip: []string{
		`h1`,
		`h1 + p`,
		`div.mailing-list-signup`,
	},
	links: `#content li a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *KleppmannScraper) ScrapeLinks(url string) (map[string]string, error) {
	return kleppmannSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *KleppmannScraper) ScrapeArticle(url string) (string, error) {
	return kleppmannSite.article(url)
}

func (s *KleppmannScraper) ScrapeTitle(url string) (string, error) {
	return kleppmannSite.scrapeTitle(url)
}

func (s *KleppmannScraper) ScrapeFilename(url string) (string, error) {
	return kleppmannSite.filename(url), nil
}
