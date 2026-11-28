package scraper

// BBCSportScraper scrapes BBC Sport (www.bbc.com/sport).
type BBCSportScraper struct {
}

var bbcSportSite = siteConfig{
	title: `article h1#main-heading`,
	body:  `article`,
	skip: []string{
		`[data-block="headline"]`,
		`[data-block="byline"]`,
		`[data-block="metadata"]`,
		`[data-block="links"]`,
		`[data-block="media"]`,
		`[data-block="topicList"]`,
		`[data-block="promoList"]`,
		`.visually-hidden`,
		`[class*="StyledFigureCopyright"]`,
	},
	links: `div[data-testid="promo"] a[href*="/articles/"]`,
	// BBC returns 403 Forbidden to colly's default User-Agent.
	userAgent: "Mozilla/5.0 (compatible; scrape)",
}

// ScrapeLinks scrapes the article links of an index page and returns a map
// of link text and URL
func (s *BBCSportScraper) ScrapeLinks(url string) (map[string]string, error) {
	return bbcSportSite.scrapeLinks(url)
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (s *BBCSportScraper) ScrapeArticle(url string) (string, error) {
	return bbcSportSite.article(url)
}

func (s *BBCSportScraper) ScrapeTitle(url string) (string, error) {
	return bbcSportSite.scrapeTitle(url)
}

// ScrapeFilename returns a file name derived from the title, as the last
// path segment of a BBC Sport URL is an opaque ID.
func (s *BBCSportScraper) ScrapeFilename(url string) (string, error) {
	title, err := s.ScrapeTitle(url)
	if err != nil {
		return "", err
	}
	return generateFileNameFromTitle(title), nil
}
