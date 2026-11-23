package scraper

// CharmScraper scrapes the Charm blog (charm.land/blog).
type CharmScraper struct {
}

var charmSite = siteConfig{
	title: `article div.title h1`,
	body:  `div.post`,
	skip: []string{
		`a.cart`,
	},
	links: `article h1 a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *CharmScraper) ScrapeLinks(url string) (map[string]string, error) {
	return charmSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *CharmScraper) ScrapeArticle(url string) (string, error) {
	return charmSite.article(url)
}

func (s *CharmScraper) ScrapeTitle(url string) (string, error) {
	return charmSite.scrapeTitle(url)
}

func (s *CharmScraper) ScrapeFilename(url string) (string, error) {
	return charmSite.filename(url), nil
}
