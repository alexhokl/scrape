package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

type MicrosoftLearnScraper struct {
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (g *MicrosoftLearnScraper) ScrapeArticle(url string) (string, error) {
	c := colly.NewCollector()

	var markdown string

	// title
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		markdown += fmt.Sprintf("# %s\n\n", e.Text)
	})

	// article
	c.OnHTML("div.content", func(e *colly.HTMLElement) {
		e.ForEach("*", func(_ int, child *colly.HTMLElement) {
			if child.DOM.Parent().IsSelection(e.DOM) {
				switch child.Name {
				case "p":
					markdown += fmt.Sprintf("%s\n\n", parseMicrosoftParagraph(child))
				case "ul":
					child.ForEach("li", func(_ int, li *colly.HTMLElement) {
						markdown += fmt.Sprintf("* %s\n", parseMicrosoftParagraph(li))
					})
					markdown += fmt.Sprintln()
				case "div":
					if child.DOM.HasClass("NOTE") {
						child.ForEach("p", func(_ int, p *colly.HTMLElement) {
							if p.DOM.HasClass("alert") {
								return
							}
							markdown += fmt.Sprintf("> %s\n\n", parseMicrosoftParagraph(p))
						})
						markdown += fmt.Sprintln()
					}
				case "h2":
					markdown += fmt.Sprintf("## %s\n\n", child.Text)
				case "h3":
					markdown += fmt.Sprintf("### %s\n\n", child.Text)
				}
			}
		})
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func (g *MicrosoftLearnScraper) ScrapeTitle(url string) (string, error) {
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

func (g *MicrosoftLearnScraper) ScrapeFilename(url string) (string, error) {
	title, err := g.ScrapeTitle(url)
	if err != nil {
		return "", err
	}
	fileName := generateFileNameFromTitle(title)
	return fileName, nil
}

func parseMicrosoftParagraph(p *colly.HTMLElement) string {
	builder := strings.Builder{}

	p.ForEach("*", func(_ int, child *colly.HTMLElement) {
		// check if child is a text node
		switch child.DOM.Nodes[0].Type {
		case html.TextNode:
			builder.WriteString(child.Text)
		case html.ElementNode:
			if child.Name == "span" {
				// it is likely an image in Microsoft Learn
				builder.WriteString("_image_")
				return
			}
			// for other element nodes, we can just bold the text
			builder.WriteString(fmt.Sprintf("**%s**", parseMicrosoftParagraph(child)))
		}
	})

	return builder.String()
}
