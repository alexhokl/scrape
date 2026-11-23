package scraper

// RickStrahlScraper scrapes Rick Strahl's Web Log (weblog.west-wind.com).
type RickStrahlScraper struct {
}

var rickStrahlSite = siteConfig{
	title: `article header h2.page-header`,
	body:  `#ArticleBody`,
	skip: []string{
		`div.document-outline`,
		`#ArticleBody > div[style]`,
		`img.da-content-image`,
	},
	links: `h2.page-header a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *RickStrahlScraper) ScrapeLinks(url string) (map[string]string, error) {
	return rickStrahlSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *RickStrahlScraper) ScrapeArticle(url string) (string, error) {
	return rickStrahlSite.article(url)
}

func (s *RickStrahlScraper) ScrapeTitle(url string) (string, error) {
	return rickStrahlSite.scrapeTitle(url)
}

func (s *RickStrahlScraper) ScrapeFilename(url string) (string, error) {
	return rickStrahlSite.filename(url), nil
}
