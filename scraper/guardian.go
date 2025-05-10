package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type GuardianScraper struct {
}

// ScrapeLinks scrapes links from the specified URL
// and returns a map of link text and URL
func (g *GuardianScraper) ScrapeLinks(url string) (map[string]string, error) {
	c := colly.NewCollector()

	articles := make(map[string]string)

	c.OnHTML("a", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Attr("data-link-name"), "group-0") {
			return
		}
		articles[e.Attr("aria-label")] = e.Request.AbsoluteURL(e.Attr("href"))
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (g *GuardianScraper) ScrapeArticle(url string) (string, error) {
	c := colly.NewCollector()

	var markdown string

	// title
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		markdown += fmt.Sprintf("# %s\n\n", e.Text)
	})

	// subtitle
	c.OnHTML("div[data-gu-name=standfirst] p", func(e *colly.HTMLElement) {
		markdown += fmt.Sprintf("## %s\n\n", e.Text)
	})

	// article body
	c.OnHTML("div.article-body-commercial-selector p", func(e *colly.HTMLElement) {
		markdown += fmt.Sprintf("%s\n\n", e.Text)
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return markdown, nil
}
