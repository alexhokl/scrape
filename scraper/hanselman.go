package scraper

// HanselmanScraper scrapes Scott Hanselman's blog (www.hanselman.com/blog).
type HanselmanScraper struct {
}

var hanselmanSite = siteConfig{
	title: `h2.blogTitle`,
	body:  `section.blogBody`,
	skip: []string{
		`.ad-item`,
		`.item-ads`,
		`.item-footer`,
		`div.about-scott`,
		`div.ads-box`,
	},
	links: `h2.blogTitle a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *HanselmanScraper) ScrapeLinks(url string) (map[string]string, error) {
	return hanselmanSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *HanselmanScraper) ScrapeArticle(url string) (string, error) {
	return hanselmanSite.article(url)
}

func (s *HanselmanScraper) ScrapeTitle(url string) (string, error) {
	return hanselmanSite.scrapeTitle(url)
}

func (s *HanselmanScraper) ScrapeFilename(url string) (string, error) {
	return hanselmanSite.filename(url), nil
}
