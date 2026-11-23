package scraper

// MitchellHashimotoScraper scrapes Mitchell Hashimoto's blog (mitchellh.com/writing).
type MitchellHashimotoScraper struct {
}

var mitchellHashimotoSite = siteConfig{
	title:    `h1.nimbus:not(.uppercase)`,
	body:     `div.prose`,
	links:    `main div.divide-y a`,
	linkText: `h1`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *MitchellHashimotoScraper) ScrapeLinks(url string) (map[string]string, error) {
	return mitchellHashimotoSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *MitchellHashimotoScraper) ScrapeArticle(url string) (string, error) {
	return mitchellHashimotoSite.article(url)
}

func (s *MitchellHashimotoScraper) ScrapeTitle(url string) (string, error) {
	return mitchellHashimotoSite.scrapeTitle(url)
}

func (s *MitchellHashimotoScraper) ScrapeFilename(url string) (string, error) {
	return mitchellHashimotoSite.filename(url), nil
}
