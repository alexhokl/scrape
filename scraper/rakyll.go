package scraper

// RakyllScraper scrapes Jaana Dogan's blog (rakyll.org).
type RakyllScraper struct {
}

var rakyllSite = siteConfig{
	title: `article h1`,
	body:  `article`,
	skip: []string{
		`h1`,
		`p.meta`,
	},
	links: `ul.posts li a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *RakyllScraper) ScrapeLinks(url string) (map[string]string, error) {
	return rakyllSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *RakyllScraper) ScrapeArticle(url string) (string, error) {
	return rakyllSite.article(url)
}

func (s *RakyllScraper) ScrapeTitle(url string) (string, error) {
	return rakyllSite.scrapeTitle(url)
}

func (s *RakyllScraper) ScrapeFilename(url string) (string, error) {
	return rakyllSite.filename(url), nil
}
