package scraper

// ObservationDeckScraper scrapes Bryan Cantrill's The Observation Deck (bcantrill.dtrace.org).
type ObservationDeckScraper struct {
}

var observationDeckSite = siteConfig{
	title: `h1.header-title`,
	body:  `article div.page-content`,
	links: `h4.post-item-title a`,
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *ObservationDeckScraper) ScrapeLinks(url string) (map[string]string, error) {
	return observationDeckSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *ObservationDeckScraper) ScrapeArticle(url string) (string, error) {
	return observationDeckSite.article(url)
}

func (s *ObservationDeckScraper) ScrapeTitle(url string) (string, error) {
	return observationDeckSite.scrapeTitle(url)
}

func (s *ObservationDeckScraper) ScrapeFilename(url string) (string, error) {
	return observationDeckSite.filename(url), nil
}
