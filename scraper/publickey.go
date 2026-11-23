package scraper

// PublickeyScraper scrapes Publickey (www.publickey1.jp).
type PublickeyScraper struct {
}

var publickeySite = siteConfig{
	title: `#maincol h1`,
	body:  `#maincol div.entrybody`,
	links: `div.newentries li a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *PublickeyScraper) ScrapeLinks(url string) (map[string]string, error) {
	return publickeySite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *PublickeyScraper) ScrapeArticle(url string) (string, error) {
	return publickeySite.article(url)
}

func (s *PublickeyScraper) ScrapeTitle(url string) (string, error) {
	return publickeySite.scrapeTitle(url)
}

func (s *PublickeyScraper) ScrapeFilename(url string) (string, error) {
	return publickeySite.filename(url), nil
}
