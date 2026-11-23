package scraper

// JuliaEvansScraper scrapes Julia Evans's blog (jvns.ca).
type JuliaEvansScraper struct {
}

var juliaEvansSite = siteConfig{
	title: `h1.entry-title`,
	body:  `article.hentry main`,
	links: `p.article-list a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *JuliaEvansScraper) ScrapeLinks(url string) (map[string]string, error) {
	return juliaEvansSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *JuliaEvansScraper) ScrapeArticle(url string) (string, error) {
	return juliaEvansSite.article(url)
}

func (s *JuliaEvansScraper) ScrapeTitle(url string) (string, error) {
	return juliaEvansSite.scrapeTitle(url)
}

func (s *JuliaEvansScraper) ScrapeFilename(url string) (string, error) {
	return juliaEvansSite.filename(url), nil
}
