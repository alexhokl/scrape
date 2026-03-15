package scraper

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

type GrafanaScraper struct {
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (g *GrafanaScraper) ScrapeArticle(url string) (string, error) {
	c := colly.NewCollector()

	var markdown string

	// title
	c.OnHTML("main h1", func(e *colly.HTMLElement) {
		if markdown == "" {
			markdown += fmt.Sprintf("# %s\n\n", strings.TrimSpace(e.Text))
		}
	})

	// article body
	c.OnHTML("div.rich-text", func(e *colly.HTMLElement) {
		markdown += parseGrafanaContent(e)
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func (g *GrafanaScraper) ScrapeTitle(url string) (string, error) {
	c := colly.NewCollector()

	var title string

	c.OnHTML("main h1", func(e *colly.HTMLElement) {
		if title == "" {
			title = strings.TrimSpace(e.Text)
		}
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return title, nil
}

func (g *GrafanaScraper) ScrapeFilename(url string) (string, error) {
	return getBasenameFromURL(url), nil
}

func parseGrafanaContent(e *colly.HTMLElement) string {
	builder := strings.Builder{}

	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		if !child.DOM.Parent().IsSelection(e.DOM) {
			return
		}

		switch child.Name {
		case "h2":
			builder.WriteString(fmt.Sprintf("## %s\n\n", parseGrafanaHeading(child)))
		case "h3":
			builder.WriteString(fmt.Sprintf("### %s\n\n", parseGrafanaHeading(child)))
		case "h4":
			builder.WriteString(fmt.Sprintf("#### %s\n\n", parseGrafanaHeading(child)))
		case "p":
			text := strings.TrimSpace(parseGrafanaInline(child))
			if text != "" {
				builder.WriteString(fmt.Sprintf("%s\n\n", text))
			}
		case "ul":
			builder.WriteString(parseGrafanaList(child, false))
		case "ol":
			builder.WriteString(parseGrafanaList(child, true))
		case "div":
			if child.DOM.HasClass("relative") {
				builder.WriteString(parseGrafanaCodeBlock(child))
			}
		case "pre":
			builder.WriteString(parseGrafanaCodeBlock(child))
		case "img":
			builder.WriteString(parseGrafanaImage(child))
		}
	})

	return builder.String()
}

// parseGrafanaHeading extracts the heading text.
// Grafana headings contain an anchor <a> link followed by a <span> with the visible text.
func parseGrafanaHeading(e *colly.HTMLElement) string {
	// The visible heading text is in the last <span> child (not inside <a>)
	text := ""
	e.ForEach("span", func(_ int, span *colly.HTMLElement) {
		// skip spans inside the anchor link
		if span.DOM.Closest("a").Length() > 0 {
			return
		}
		if t := strings.TrimSpace(span.Text); t != "" {
			text = t
		}
	})
	if text == "" {
		text = strings.TrimSpace(e.Text)
	}
	return text
}

// parseGrafanaCodeBlock renders a code block element as a markdown fenced code block.
// It accepts either a div.relative wrapper containing a <pre> or a bare <pre> element.
func parseGrafanaCodeBlock(e *colly.HTMLElement) string {
	builder := strings.Builder{}
	var pre *colly.HTMLElement
	if e.Name == "pre" {
		pre = e
	} else {
		e.ForEach("pre", func(_ int, p *colly.HTMLElement) {
			if pre == nil {
				pre = p
			}
		})
	}
	if pre == nil {
		return ""
	}
	code := strings.TrimRight(pre.Text, "\n")
	builder.WriteString(fmt.Sprintf("```\n%s\n```\n\n", code))
	return builder.String()
}

// parseGrafanaList renders a <ul> or <ol> element as markdown.
// Grafana list items wrap their content in a <div> child.
func parseGrafanaList(e *colly.HTMLElement, ordered bool) string {
	builder := strings.Builder{}
	index := 0
	e.ForEach("li", func(_ int, li *colly.HTMLElement) {
		// only process direct children of this list element
		if !li.DOM.Parent().IsSelection(e.DOM) {
			return
		}
		index++

		bulletText := ""
		// Grafana wraps list item text in a <div> child
		li.ForEach("div", func(_ int, d *colly.HTMLElement) {
			if !d.DOM.Parent().IsSelection(li.DOM) {
				return
			}
			if bulletText == "" {
				bulletText = strings.TrimSpace(parseGrafanaInline(d))
			}
		})
		if bulletText == "" {
			bulletText = strings.TrimSpace(parseGrafanaInline(li))
		}

		if ordered {
			builder.WriteString(fmt.Sprintf("%d. %s\n", index, bulletText))
		} else {
			builder.WriteString(fmt.Sprintf("* %s\n", bulletText))
		}
	})
	builder.WriteString("\n")
	return builder.String()
}

// parseGrafanaImage renders an <img> element as a markdown image.
// Grafana blog pages serve images through the Next.js image proxy
// (/mw/_next/image/?url=<encoded>&w=...); the real image URL is extracted
// from the "url" query parameter when present.
func parseGrafanaImage(e *colly.HTMLElement) string {
	src := e.Attr("src")
	alt := e.Attr("alt")
	if src == "" {
		return ""
	}
	imgURL := parseGrafanaImageURL(src)
	return fmt.Sprintf("![%s](%s)\n\n", alt, imgURL)
}

// parseGrafanaImageURL extracts the real image URL from a Next.js proxy src.
// If the src contains a "url" query parameter, its decoded value is returned.
// Otherwise the original src is returned unchanged.
func parseGrafanaImageURL(src string) string {
	parsed, err := url.Parse(src)
	if err != nil {
		return src
	}
	if raw := parsed.Query().Get("url"); raw != "" {
		return raw
	}
	return src
}

// parseGrafanaInline renders the inline content of an element as markdown,
// preserving <code> as backtick spans and <a> as markdown links.
func parseGrafanaInline(e *colly.HTMLElement) string {
	builder := strings.Builder{}
	for _, node := range e.DOM.Contents().Nodes {
		switch node.Type {
		case html.TextNode:
			builder.WriteString(node.Data)
		case html.ElementNode:
			sel := e.DOM.FindNodes(node)
			switch node.Data {
			case "code":
				builder.WriteString(fmt.Sprintf("`%s`", sel.Text()))
			case "a":
				href, _ := sel.Attr("href")
				linkText := parseGrafanaInlineNodes(sel.Contents().Nodes)
				if href != "" {
					builder.WriteString(fmt.Sprintf("[%s](%s)", linkText, href))
				} else {
					builder.WriteString(linkText)
				}
			case "strong":
				builder.WriteString(fmt.Sprintf("**%s**", sel.Text()))
			default:
				builder.WriteString(sel.Text())
			}
		}
	}
	return builder.String()
}

// parseGrafanaInlineNodes renders a slice of HTML nodes as inline markdown.
func parseGrafanaInlineNodes(nodes []*html.Node) string {
	builder := strings.Builder{}
	for _, node := range nodes {
		switch node.Type {
		case html.TextNode:
			builder.WriteString(node.Data)
		case html.ElementNode:
			if node.Data == "code" {
				var codeText strings.Builder
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						codeText.WriteString(c.Data)
					}
				}
				builder.WriteString(fmt.Sprintf("`%s`", codeText.String()))
			} else {
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						builder.WriteString(c.Data)
					}
				}
			}
		}
	}
	return builder.String()
}
