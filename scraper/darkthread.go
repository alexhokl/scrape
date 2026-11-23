package scraper

// DarkthreadScraper scrapes 黑暗執行緒 (blog.darkthread.net).
type DarkthreadScraper struct {
}

var darkthreadSite = siteConfig{
	title: `h1.post-title`,
	body:  `article.post div[itemprop~="articleBody"]`,
	skip: []string{
		`p.hide`,
	},
	links: `article.post header h1 a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *DarkthreadScraper) ScrapeLinks(url string) (map[string]string, error) {
	return darkthreadSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *DarkthreadScraper) ScrapeArticle(url string) (string, error) {
	return darkthreadSite.article(url)
}

func (s *DarkthreadScraper) ScrapeTitle(url string) (string, error) {
	return darkthreadSite.scrapeTitle(url)
}

func (s *DarkthreadScraper) ScrapeFilename(url string) (string, error) {
	return darkthreadSite.filename(url), nil
}
