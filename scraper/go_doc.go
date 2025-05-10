package scraper

import (
	"fmt"

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
		for _, n := range e.DOM.Children().Nodes {
			if n.Type != html.ElementNode {
				continue
			}
			switch n.Data {
			case "h2":
				markdown += fmt.Sprintf("## %s\n\n", n.FirstChild.Data)
			case "h3":
				markdown += fmt.Sprintf("### %s\n\n", n.FirstChild.Data)
			case "p":
				markdown += fmt.Sprintf("%s\n\n", getTextInNode(n))
			case "ul":
				for li := range n.ChildNodes() {
					if li.Type != html.ElementNode {
						continue
					}
					if li.Data == "li" {
						markdown += fmt.Sprintf("* %s\n", getTextInNode(li))
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
							markdown += fmt.Sprintf("> %s\n\n", getTextInNode(p))
						}
					}
					markdown += fmt.Sprintln()
				}
			case "pre":
				markdown += "```\n"
				for child := range n.ChildNodes() {
					if child.Type == html.ElementNode && child.Data == "code" {
						markdown += fmt.Sprintln(getTextInNode(child))
					}
				}
				markdown += "```\n\n"
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

func getTextInNode(node *html.Node) string {
	markdown := ""
	for child := range node.ChildNodes() {
		switch child.Type {
		case html.TextNode:
			markdown += fmt.Sprint(child.Data)
		case html.ElementNode:
			if child.Data == "img" {
				markdown += "_image_"
				continue
			}
			text := getTextInNode(child)
			if text == "" {
				continue
			}
			if node.Data == "li" && child.Data == "ul" {
				// list in list
				for subListItem := range child.ChildNodes() {
					if subListItem.Type != html.ElementNode {
						continue
					}
					if subListItem.Data == "li" {
						markdown += fmt.Sprintf("  * %s\n", getTextInNode(subListItem))
					}
				}
				continue
			}
			markdown += fmt.Sprintf("**%s**", text)
		default:
		}
	}
	return markdown
}
