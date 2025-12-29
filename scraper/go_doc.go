package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

type GoDocScraper struct {
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (g *GoDocScraper) ScrapeArticle(url string) (string, error) {
	c := colly.NewCollector()

	var markdown string

	// title
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		markdown += fmt.Sprintf("# %s\n\n", e.Text)
	})

	// article
	c.OnHTML("article", func(e *colly.HTMLElement) {
		e.ForEach("*", func(_ int, child *colly.HTMLElement) {
			if child.DOM.Parent().IsSelection(e.DOM) {
				switch child.Name {
				case "h2":
					markdown += fmt.Sprintf("## %s\n\n", child.Text)
				case "h3":
					markdown += fmt.Sprintf("### %s\n\n", child.Text)
				case "p":
					markdown += fmt.Sprintf("%s\n\n", parseGoDocParagraph(child))
				case "ul":
					child.ForEach("li", func(_ int, li *colly.HTMLElement) {
						markdown += fmt.Sprintf("* %s\n", parseGoDocParagraph(li))
					})
					markdown += fmt.Sprintln()
				case "div":
					if child.DOM.HasClass("NOTE") {
						child.ForEach("p", func(_ int, p *colly.HTMLElement) {
							if p.DOM.HasClass("alert") {
								return
							}
							markdown += fmt.Sprintf("> %s\n\n", parseGoDocParagraph(p))
						})
						markdown += fmt.Sprintln()
					}
				case "pre":
					markdown += fmt.Sprintln("```")
					child.ForEach("code", func(_ int, code *colly.HTMLElement) {
						markdown += fmt.Sprintln(parseGoDocParagraph(code))
					})
					markdown += fmt.Sprintln("```")
					markdown += fmt.Sprintln()
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

func (g *GoDocScraper) ScrapeTitle(url string) (string, error) {
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

func (g *GoDocScraper) ScrapeFilename(url string) (string, error) {
	title, err := g.ScrapeTitle(url)
	if err != nil {
		return "", err
	}

	return generateFileNameFromTitle(title), nil
}

func parseGoDocParagraph(p *colly.HTMLElement) string {
	builder := strings.Builder{}

	p.ForEach("*", func(_ int, child *colly.HTMLElement) {
		switch child.DOM.Nodes[0].Type {
		case html.TextNode:
			builder.WriteString(child.Text)
		case html.ElementNode:
			if child.Name == "img" {
				// it is likely an image in GoDoc
				builder.WriteString("_image_")
				return
			}
			text := parseGoDocParagraph(child)
			if text == "" {
				return
			}
			if p.Name == "li" && child.Name == "ul" {
				// list in list
				child.ForEach("li", func(_ int, subListItem *colly.HTMLElement) {
					builder.WriteString(fmt.Sprintf("  * %s\n", parseGoDocParagraph(subListItem)))
				})
				return
			}
			builder.WriteString(fmt.Sprintf("**%s**", text))
		}
	})

	return builder.String()
}
