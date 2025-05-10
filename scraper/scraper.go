package scraper

type LinkScraper interface {
	// ScrapeLinks scrapes links from the specified URL
	// and returns a map of link text and URL
	ScrapeLinks(url string) (map[string]string, error)
}

type ArticleScraper interface {
	// ScrapeArticle scrapes the article content from the specified URL
	// and returns markdown in a string
	ScrapeArticle(url string) (string, error)
}
