package scraper

import (
	"fmt"

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
		for _, n := range e.DOM.Children().Nodes {
			if n.Type != html.ElementNode {
				continue
			}
			switch n.Data {
			case "p":
				markdown += fmt.Sprintf("%s\n\n", dumpTextInNode(n))
			case "ul":
				for li := range n.ChildNodes() {
					if li.Type == html.ElementNode && li.Data == "li" {
						markdown += fmt.Sprintf("* %s\n", dumpTextInNode(li))
					}
				}
				markdown += fmt.Sprintln()
			case "div":
				if containsCssClass(n, "NOTE") {
					for p := range n.ChildNodes() {
						if p.Type == html.ElementNode && p.Data == "p" {
							if containsCssClass(p, "alert") {
								continue
							}
							markdown += fmt.Sprintf("> %s\n\n", dumpTextInNode(p))
						}
					}
					markdown += fmt.Sprintln()
				}
			case "h2":
				markdown += fmt.Sprintf("## %s\n\n", n.FirstChild.Data)
			case "h3":
				markdown += fmt.Sprintf("### %s\n\n", n.FirstChild.Data)
			default:
			}
		}
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func dumpTextInNode(node *html.Node) string {
	markdown := ""
	for child := range node.ChildNodes() {
		switch child.Type {
		case html.TextNode:
			markdown += fmt.Sprint(child.Data)
		case html.ElementNode:
			if child.Data == "span" {
				// it is likely an image in Microsoft Learn
				markdown += "_image_"
				continue
			}
			markdown += fmt.Sprintf("**%s**", dumpTextInNode(child))
		default:
		}
	}
	return markdown
}
