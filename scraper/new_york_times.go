package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type NewYorkTimesScraper struct {
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (g *NewYorkTimesScraper) ScrapeArticle(url string) (string, error) {
	c := colly.NewCollector()

	var markdown string

	// title
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		markdown += fmt.Sprintf("# %s\n\n", e.Text)
	})

	// article body
	c.OnHTML("div.article-content-container p", func(e *colly.HTMLElement) {
		markdown += fmt.Sprintf("%s\n\n", e.Text)
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func (g *NewYorkTimesScraper) ScrapeTitle(url string) (string, error) {
	c := colly.NewCollector()

	var title string

	// title
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		title = strings.TrimSpace(e.Text)
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return title, nil
}

func (g *NewYorkTimesScraper) ScrapeFilename(url string) (string, error) {
	title, err := g.ScrapeTitle(url)
	if err != nil {
		return "", err
	}

	return generateFileNameFromTitle(title), nil
}
