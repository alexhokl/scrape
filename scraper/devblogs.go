package scraper

// DevBlogsScraper scrapes Microsoft developer blogs (devblogs.microsoft.com).
type DevBlogsScraper struct {
}

var devBlogsSite = siteConfig{
	title: `header h1`,
	body:  `div.entry-content`,
	links: `a.excerpt-title, div.pinned-wrap h3 a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *DevBlogsScraper) ScrapeLinks(url string) (map[string]string, error) {
	return devBlogsSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *DevBlogsScraper) ScrapeArticle(url string) (string, error) {
	return devBlogsSite.article(url)
}

func (s *DevBlogsScraper) ScrapeTitle(url string) (string, error) {
	return devBlogsSite.scrapeTitle(url)
}

func (s *DevBlogsScraper) ScrapeFilename(url string) (string, error) {
	return devBlogsSite.filename(url), nil
}
