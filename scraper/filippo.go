package scraper

// FilippoScraper scrapes Filippo Valsorda's blog (words.filippo.io).
type FilippoScraper struct {
}

var filippoSite = siteConfig{
	title: `article h1`,
	body:  `article section`,
	links: `article h1 a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *FilippoScraper) ScrapeLinks(url string) (map[string]string, error) {
	return filippoSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *FilippoScraper) ScrapeArticle(url string) (string, error) {
	return filippoSite.article(url)
}

func (s *FilippoScraper) ScrapeTitle(url string) (string, error) {
	return filippoSite.scrapeTitle(url)
}

func (s *FilippoScraper) ScrapeFilename(url string) (string, error) {
	return filippoSite.filename(url), nil
}
