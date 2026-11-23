package scraper

// BellingcatScraper scrapes Bellingcat (www.bellingcat.com).
type BellingcatScraper struct {
}

var bellingcatSite = siteConfig{
	title: `div.singular__content__text__title h1`,
	body:  `div.singular__content__text__content`,
	skip: []string{
		`blockquote.instagram-media`,
	},
	links: `a.featured-post__title, a.teaser__title--default`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *BellingcatScraper) ScrapeLinks(url string) (map[string]string, error) {
	return bellingcatSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *BellingcatScraper) ScrapeArticle(url string) (string, error) {
	return bellingcatSite.article(url)
}

func (s *BellingcatScraper) ScrapeTitle(url string) (string, error) {
	return bellingcatSite.scrapeTitle(url)
}

func (s *BellingcatScraper) ScrapeFilename(url string) (string, error) {
	return bellingcatSite.filename(url), nil
}
