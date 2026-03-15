package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

type OllamaScraper struct {
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (o *OllamaScraper) ScrapeArticle(url string) (string, error) {
	collector := colly.NewCollector()

	var markdown string

	// title
	collector.OnHTML("article h1", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		if text != "" {
			markdown += fmt.Sprintf("# %s\n\n", text)
		}
	})

	// article body — the prose section contains all content elements
	collector.OnHTML("article section.prose", func(e *colly.HTMLElement) {
		markdown += parseOllamaContent(e)
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func (o *OllamaScraper) ScrapeTitle(url string) (string, error) {
	collector := colly.NewCollector()

	var title string

	collector.OnHTML("article h1", func(e *colly.HTMLElement) {
		title = strings.TrimSpace(e.Text)
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	return title, nil
}

func (o *OllamaScraper) ScrapeFilename(url string) (string, error) {
	return getBasenameFromURL(url), nil
}

func parseOllamaContent(e *colly.HTMLElement) string {
	builder := strings.Builder{}

	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		if !child.DOM.Parent().IsSelection(e.DOM) {
			return
		}

		switch child.Name {
		case "h2":
			builder.WriteString(fmt.Sprintf("## %s\n\n", strings.TrimSpace(child.Text)))
		case "h3":
			builder.WriteString(fmt.Sprintf("### %s\n\n", strings.TrimSpace(child.Text)))
		case "h4":
			builder.WriteString(fmt.Sprintf("#### %s\n\n", strings.TrimSpace(child.Text)))
		case "p":
			// A <p> that contains only an <img> is rendered as an image block.
			if isOllamaImageParagraph(child) {
				child.ForEach("img", func(_ int, img *colly.HTMLElement) {
					src := img.Attr("src")
					alt := img.Attr("alt")
					if src != "" {
						builder.WriteString(fmt.Sprintf("![%s](%s)\n\n", alt, src))
					}
				})
				return
			}
			text := strings.TrimSpace(parseOllamaInline(child))
			if text != "" {
				builder.WriteString(fmt.Sprintf("%s\n\n", text))
			}
		case "ul":
			builder.WriteString(parseOllamaList(child, false))
		case "ol":
			builder.WriteString(parseOllamaList(child, true))
		case "pre":
			builder.WriteString(parseOllamaCodeBlock(child))
		case "blockquote":
			text := strings.TrimSpace(child.Text)
			if text != "" {
				builder.WriteString(fmt.Sprintf("> %s\n\n", text))
			}
		}
	})

	return builder.String()
}

// isOllamaImageParagraph returns true when a <p> element contains only
// image elements (and optional whitespace text nodes).
func isOllamaImageParagraph(e *colly.HTMLElement) bool {
	hasImg := false
	for _, node := range e.DOM.Contents().Nodes {
		switch node.Type {
		case html.TextNode:
			if strings.TrimSpace(node.Data) != "" {
				return false
			}
		case html.ElementNode:
			if node.Data != "img" {
				return false
			}
			hasImg = true
		}
	}
	return hasImg
}

// parseOllamaCodeBlock renders a <pre> element as a markdown fenced code block.
// The language is taken from the class of the inner <code> element.
func parseOllamaCodeBlock(e *colly.HTMLElement) string {
	builder := strings.Builder{}
	lang := ""
	e.ForEach("code", func(_ int, code *colly.HTMLElement) {
		lang = parseOllamaCodeLang(code)
	})
	code := strings.TrimRight(e.Text, "\n")
	if lang != "" {
		builder.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", lang, code))
	} else {
		builder.WriteString(fmt.Sprintf("```\n%s\n```\n\n", code))
	}
	return builder.String()
}

// parseOllamaCodeLang extracts the language identifier from a <code> element's class.
// e.g. class="language-bash" → "bash"
func parseOllamaCodeLang(e *colly.HTMLElement) string {
	class := e.Attr("class")
	for _, part := range strings.Fields(class) {
		if strings.HasPrefix(part, "language-") {
			return strings.TrimPrefix(part, "language-")
		}
	}
	return ""
}

// parseOllamaList renders a <ul> or <ol> element as markdown.
func parseOllamaList(e *colly.HTMLElement, ordered bool) string {
	builder := strings.Builder{}
	index := 0
	e.ForEach("li", func(_ int, li *colly.HTMLElement) {
		if !li.DOM.Parent().IsSelection(e.DOM) {
			return
		}
		index++
		text := strings.TrimSpace(parseOllamaInline(li))
		if ordered {
			builder.WriteString(fmt.Sprintf("%d. %s\n", index, text))
		} else {
			builder.WriteString(fmt.Sprintf("* %s\n", text))
		}
	})
	builder.WriteString("\n")
	return builder.String()
}

// parseOllamaInline renders the inline content of an element as markdown,
// preserving <code> as backtick spans, <a> as markdown links, <b>/<strong>
// as bold, and <i>/<em> as italic.
func parseOllamaInline(e *colly.HTMLElement) string {
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
				linkText := parseOllamaInlineNodes(sel.Contents().Nodes)
				linkText = strings.TrimSpace(linkText)
				if href != "" && linkText != "" {
					builder.WriteString(fmt.Sprintf("[%s](%s)", linkText, href))
				} else {
					builder.WriteString(linkText)
				}
			case "b", "strong":
				text := parseOllamaInlineNodes(sel.Contents().Nodes)
				trimmed := strings.TrimSpace(text)
				if trimmed != "" {
					builder.WriteString(fmt.Sprintf("**%s**", trimmed))
				}
			case "i", "em":
				text := parseOllamaInlineNodes(sel.Contents().Nodes)
				trimmed := strings.TrimSpace(text)
				if trimmed != "" {
					builder.WriteString(fmt.Sprintf("*%s*", trimmed))
				}
			default:
				builder.WriteString(sel.Text())
			}
		}
	}
	return builder.String()
}

// parseOllamaInlineNodes renders a slice of HTML nodes as inline markdown,
// used for content inside elements like <a> where we still want <code> preserved.
func parseOllamaInlineNodes(nodes []*html.Node) string {
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
