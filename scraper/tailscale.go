package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

type TailscaleScraper struct {
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (t *TailscaleScraper) ScrapeArticle(url string) (string, error) {
	c := colly.NewCollector()

	var markdown string

	// title
	c.OnHTML("article#main-content > h1", func(e *colly.HTMLElement) {
		markdown += fmt.Sprintf("# %s\n\n", strings.TrimSpace(e.Text))
	})

	// article body
	c.OnHTML("article#main-content > div.ts-prose", func(e *colly.HTMLElement) {
		markdown += parseTailscaleContent(e)
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func (t *TailscaleScraper) ScrapeTitle(url string) (string, error) {
	c := colly.NewCollector()

	var title string

	c.OnHTML("article#main-content > h1", func(e *colly.HTMLElement) {
		title = strings.TrimSpace(e.Text)
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return title, nil
}

func (t *TailscaleScraper) ScrapeFilename(url string) (string, error) {
	return getBasenameFromURL(url), nil
}

func parseTailscaleContent(e *colly.HTMLElement) string {
	builder := strings.Builder{}

	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		if !child.DOM.Parent().IsSelection(e.DOM) {
			return
		}

		switch child.Name {
		case "h2":
			builder.WriteString(fmt.Sprintf("## %s\n\n", parseTailscaleHeading(child)))
		case "h3":
			builder.WriteString(fmt.Sprintf("### %s\n\n", parseTailscaleHeading(child)))
		case "h4":
			builder.WriteString(fmt.Sprintf("#### %s\n\n", parseTailscaleHeading(child)))
		case "p":
			text := strings.TrimSpace(parseTailscaleInline(child))
			if text != "" {
				builder.WriteString(fmt.Sprintf("%s\n\n", text))
			}
		case "ul":
			builder.WriteString(parseTailscaleList(child, false))
		case "ol":
			builder.WriteString(parseTailscaleList(child, true))
		case "div":
			if child.DOM.HasClass("note") {
				child.ForEach("*", func(_ int, noteChild *colly.HTMLElement) {
					if !noteChild.DOM.Parent().IsSelection(child.DOM) {
						return
					}
					switch noteChild.Name {
					case "p":
						text := strings.TrimSpace(parseTailscaleInline(noteChild))
						if text != "" {
							builder.WriteString(fmt.Sprintf("> %s\n\n", text))
						}
					case "div":
						if noteChild.DOM.HasClass("group") && noteChild.DOM.HasClass("relative") && noteChild.DOM.HasClass("overflow-hidden") {
							builder.WriteString(parseTailscaleCodeBlock(noteChild))
						}
					}
				})
			} else if child.DOM.HasClass("group") && child.DOM.HasClass("relative") && child.DOM.HasClass("overflow-hidden") {
				builder.WriteString(parseTailscaleCodeBlock(child))
			}
		case "pre":
			// standalone pre (not inside a div.group wrapper)
			builder.WriteString(parseTailscaleCodeBlock(child))
		}
	})

	return builder.String()
}

// parseTailscaleCodeBlock renders a code block element as a markdown fenced code block.
// It accepts either a div.group.relative.overflow-hidden wrapper or a bare <pre> element.
func parseTailscaleCodeBlock(e *colly.HTMLElement) string {
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
	lang := parseTailscaleCodeLang(pre)
	code := strings.TrimRight(pre.Text, "\n")
	if lang != "" {
		builder.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", lang, code))
	} else {
		builder.WriteString(fmt.Sprintf("```\n%s\n```\n\n", code))
	}
	return builder.String()
}

// parseTailscaleHeading extracts the heading text.
// Tailscale headings wrap the text in <a><span id="inner-text">...</span></a>,
// so we prefer that span's text when available, falling back to the element text.
func parseTailscaleHeading(e *colly.HTMLElement) string {
	text := e.ChildText("span#inner-text")
	if text == "" {
		text = e.Text
	}
	return strings.TrimSpace(text)
}

// parseTailscaleList renders a <ul> or <ol> element as markdown.
// Each list item's direct children are examined: <p> elements provide the
// bullet text, and any div.group code block wrappers are emitted as fenced
// code blocks after the bullet.
func parseTailscaleList(e *colly.HTMLElement, ordered bool) string {
	builder := strings.Builder{}
	index := 0
	e.ForEach("li", func(_ int, li *colly.HTMLElement) {
		// only process direct children of this list element
		if !li.DOM.Parent().IsSelection(e.DOM) {
			return
		}
		index++

		// Collect bullet text from direct <p> children (or the li text when
		// there are no child block elements).
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
					bulletText = strings.TrimSpace(parseTailscaleInline(p))
				}
			})
		} else {
			bulletText = strings.TrimSpace(parseTailscaleInline(li))
		}

		if ordered {
			builder.WriteString(fmt.Sprintf("%d. %s\n", index, bulletText))
		} else {
			builder.WriteString(fmt.Sprintf("* %s\n", bulletText))
		}

		// Emit any code blocks that are direct children of this <li>.
		li.ForEach("div", func(_ int, d *colly.HTMLElement) {
			if !d.DOM.Parent().IsSelection(li.DOM) {
				return
			}
			if d.DOM.HasClass("group") && d.DOM.HasClass("relative") && d.DOM.HasClass("overflow-hidden") {
				builder.WriteString(parseTailscaleCodeBlock(d))
			}
		})
	})
	builder.WriteString("\n")
	return builder.String()
}

// parseTailscaleCodeLang extracts the language identifier from a <pre> element's class.
// e.g. class="refractor language-shell" → "shell"
func parseTailscaleCodeLang(e *colly.HTMLElement) string {
	class := e.Attr("class")
	for _, part := range strings.Fields(class) {
		if strings.HasPrefix(part, "language-") {
			return strings.TrimPrefix(part, "language-")
		}
	}
	return ""
}

// parseTailscaleInline renders the inline content of an element as markdown,
// preserving <code> as backtick spans and <a> as markdown links.
// Text nodes are emitted as-is; all other elements fall back to their text content.
func parseTailscaleInline(e *colly.HTMLElement) string {
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
				linkText := parseTailscaleInlineNodes(sel.Contents().Nodes)
				if href != "" {
					builder.WriteString(fmt.Sprintf("[%s](%s)", linkText, href))
				} else {
					builder.WriteString(linkText)
				}
			default:
				builder.WriteString(sel.Text())
			}
		}
	}
	return builder.String()
}

// parseTailscaleInlineNodes renders a slice of HTML nodes as inline markdown,
// used for content inside elements like <a> where we still want <code> preserved.
func parseTailscaleInlineNodes(nodes []*html.Node) string {
	builder := strings.Builder{}
	for _, node := range nodes {
		switch node.Type {
		case html.TextNode:
			builder.WriteString(node.Data)
		case html.ElementNode:
			if node.Data == "code" {
				// collect text content of the code node
				var codeText strings.Builder
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						codeText.WriteString(c.Data)
					}
				}
				builder.WriteString(fmt.Sprintf("`%s`", codeText.String()))
			} else {
				// fallback: emit text content
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
