package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

type CloudflareScraper struct {
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (c *CloudflareScraper) ScrapeArticle(url string) (string, error) {
	collector := colly.NewCollector()

	var markdown string

	// title
	collector.OnHTML("article.post-full > h1", func(e *colly.HTMLElement) {
		markdown += fmt.Sprintf("# %s\n\n", strings.TrimSpace(e.Text))
	})

	// article body — use the first div.post-content only (skip the boilerplate footer)
	found := false
	collector.OnHTML("section.post-full-content div.post-content", func(e *colly.HTMLElement) {
		if found {
			return
		}
		found = true
		markdown += parseCloudflareContent(e)
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func (c *CloudflareScraper) ScrapeTitle(url string) (string, error) {
	collector := colly.NewCollector()

	var title string

	collector.OnHTML("article.post-full > h1", func(e *colly.HTMLElement) {
		title = strings.TrimSpace(e.Text)
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	return title, nil
}

func (c *CloudflareScraper) ScrapeFilename(url string) (string, error) {
	return getBasenameFromURL(url), nil
}

func parseCloudflareContent(e *colly.HTMLElement) string {
	builder := strings.Builder{}

	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		if !child.DOM.Parent().IsSelection(e.DOM) {
			return
		}

		switch child.Name {
		case "div":
			// headings are wrapped in <div class="flex anchor relative">
			if child.DOM.HasClass("flex") && child.DOM.HasClass("anchor") {
				child.ForEach("h2, h3, h4", func(_ int, heading *colly.HTMLElement) {
					switch heading.Name {
					case "h2":
						builder.WriteString(fmt.Sprintf("## %s\n\n", strings.TrimSpace(heading.Text)))
					case "h3":
						builder.WriteString(fmt.Sprintf("### %s\n\n", strings.TrimSpace(heading.Text)))
					case "h4":
						builder.WriteString(fmt.Sprintf("#### %s\n\n", strings.TrimSpace(heading.Text)))
					}
				})
			}
		case "p":
			text := strings.TrimSpace(parseCloudflareInline(child))
			if text != "" {
				builder.WriteString(fmt.Sprintf("%s\n\n", text))
			}
		case "ul":
			builder.WriteString(parseCloudflareList(child, false))
		case "ol":
			builder.WriteString(parseCloudflareList(child, true))
		case "pre":
			builder.WriteString(parseCloudflareCodeBlock(child))
		case "blockquote":
			text := strings.TrimSpace(child.Text)
			if text != "" {
				builder.WriteString(fmt.Sprintf("> %s\n\n", text))
			}
		case "figure":
			if child.DOM.HasClass("kg-image-card") {
				child.ForEach("img", func(_ int, img *colly.HTMLElement) {
					src := img.Attr("src")
					alt := img.Attr("alt")
					if src != "" {
						builder.WriteString(fmt.Sprintf("![%s](%s)\n\n", alt, src))
					}
				})
			}
		}
	})

	return builder.String()
}

// parseCloudflareCodeBlock renders a <pre> element as a markdown fenced code block.
func parseCloudflareCodeBlock(e *colly.HTMLElement) string {
	builder := strings.Builder{}
	lang := parseCloudflareCodeLang(e)
	code := strings.TrimRight(e.Text, "\n")
	if lang != "" {
		builder.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", lang, code))
	} else {
		builder.WriteString(fmt.Sprintf("```\n%s\n```\n\n", code))
	}
	return builder.String()
}

// parseCloudflareCodeLang extracts the language identifier from a <pre> element's class.
// e.g. class="language-typescript" → "typescript"
func parseCloudflareCodeLang(e *colly.HTMLElement) string {
	class := e.Attr("class")
	for _, part := range strings.Fields(class) {
		if strings.HasPrefix(part, "language-") {
			return strings.TrimPrefix(part, "language-")
		}
	}
	return ""
}

// parseCloudflareList renders a <ul> or <ol> element as markdown.
func parseCloudflareList(e *colly.HTMLElement, ordered bool) string {
	builder := strings.Builder{}
	index := 0
	e.ForEach("li", func(_ int, li *colly.HTMLElement) {
		if !li.DOM.Parent().IsSelection(e.DOM) {
			return
		}
		index++

		// List items may contain <p> children or bare text.
		bulletText := ""
		hasBlockChildren := li.DOM.Children().FilterFunction(func(_ int, s *goquery.Selection) bool {
			name := goquery.NodeName(s)
			return name == "p" || name == "div" || name == "ul" || name == "ol"
		}).Length() > 0

		if hasBlockChildren {
			li.ForEach("p", func(_ int, p *colly.HTMLElement) {
				if !p.DOM.Parent().IsSelection(li.DOM) {
					return
				}
				if bulletText == "" {
					bulletText = strings.TrimSpace(parseCloudflareInline(p))
				}
			})
		} else {
			bulletText = strings.TrimSpace(parseCloudflareInline(li))
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

// parseCloudflareInline renders the inline content of an element as markdown,
// preserving <code> as backtick spans, <a> as markdown links, <b> as bold, and <i> as italic.
func parseCloudflareInline(e *colly.HTMLElement) string {
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
				linkText := parseCloudflareInlineNodes(sel.Contents().Nodes)
				if href != "" {
					builder.WriteString(fmt.Sprintf("[%s](%s)", linkText, href))
				} else {
					builder.WriteString(linkText)
				}
			case "b", "strong":
				text := parseCloudflareInlineNodes(sel.Contents().Nodes)
				trimmed := strings.TrimSpace(text)
				if trimmed != "" {
					builder.WriteString(fmt.Sprintf("**%s**", trimmed))
				}
			case "i", "em":
				text := parseCloudflareInlineNodes(sel.Contents().Nodes)
				trimmed := strings.TrimSpace(text)
				if trimmed != "" {
					builder.WriteString(fmt.Sprintf("*%s*", trimmed))
				}
			case "u":
				builder.WriteString(sel.Text())
			default:
				builder.WriteString(sel.Text())
			}
		}
	}
	return builder.String()
}

// parseCloudflareInlineNodes renders a slice of HTML nodes as inline markdown,
// used for content inside elements like <a> where we still want <code> preserved.
func parseCloudflareInlineNodes(nodes []*html.Node) string {
	builder := strings.Builder{}
	for _, node := range nodes {
		switch node.Type {
		case html.TextNode:
			builder.WriteString(node.Data)
		case html.ElementNode:
			switch node.Data {
			case "code":
				var codeText strings.Builder
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						codeText.WriteString(c.Data)
					}
				}
				builder.WriteString(fmt.Sprintf("`%s`", codeText.String()))
			case "u":
				// <u> inside <a> — just emit the text
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						builder.WriteString(c.Data)
					}
				}
			default:
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
