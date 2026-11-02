package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

type ProPublicaScraper struct {
}

// the title is rendered twice (for small and large screens), so only the
// first match is used
const proPublicaTitleSelector = "h1.p-opener__hed"

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (p *ProPublicaScraper) ScrapeArticle(url string) (string, error) {
	collector := colly.NewCollector()

	var title string
	var dek string
	var body string

	collector.OnHTML(proPublicaTitleSelector, func(e *colly.HTMLElement) {
		if title == "" {
			title = strings.TrimSpace(e.Text)
		}
	})

	collector.OnHTML("h2.p-opener__dek", func(e *colly.HTMLElement) {
		if dek == "" {
			dek = strings.TrimSpace(e.Text)
		}
	})

	// article body; the newsletter pitch, co-publish notes, corrections and
	// recirculation are in <aside> elements outside of entry-content
	collector.OnHTML("div.article-body div.entry-content", func(e *colly.HTMLElement) {
		body += parseProPublicaContent(e)
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	if title != "" {
		fmt.Fprintf(&builder, "# %s\n\n", title)
	}
	if dek != "" {
		fmt.Fprintf(&builder, "*%s*\n\n", dek)
	}
	builder.WriteString(body)

	return builder.String(), nil
}

func (p *ProPublicaScraper) ScrapeTitle(url string) (string, error) {
	collector := colly.NewCollector()

	var title string

	collector.OnHTML(proPublicaTitleSelector, func(e *colly.HTMLElement) {
		if title == "" {
			title = strings.TrimSpace(e.Text)
		}
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	return title, nil
}

func (p *ProPublicaScraper) ScrapeFilename(url string) (string, error) {
	return getBasenameFromURL(url), nil
}

func parseProPublicaContent(e *colly.HTMLElement) string {
	builder := strings.Builder{}

	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		if !child.DOM.Parent().IsSelection(e.DOM) {
			return
		}

		switch child.Name {
		case "p":
			text := strings.TrimSpace(parseProPublicaInline(child))
			if text != "" {
				fmt.Fprintf(&builder, "%s\n\n", text)
			}
		case "h2":
			fmt.Fprintf(&builder, "## %s\n\n", strings.TrimSpace(child.Text))
		case "h3":
			fmt.Fprintf(&builder, "### %s\n\n", strings.TrimSpace(child.Text))
		case "h4":
			fmt.Fprintf(&builder, "#### %s\n\n", strings.TrimSpace(child.Text))
		case "div":
			// a lead-in wraps a heading and the content following it
			if child.DOM.HasClass("wp-block-propublica-lead-in") {
				builder.WriteString(parseProPublicaContent(child))
			}
		case "figure":
			builder.WriteString(parseProPublicaImage(child))
		case "ul":
			builder.WriteString(parseProPublicaList(child, false))
		case "ol":
			builder.WriteString(parseProPublicaList(child, true))
		case "blockquote":
			text := strings.TrimSpace(child.Text)
			if text != "" {
				fmt.Fprintf(&builder, "> %s\n\n", text)
			}
		}
	})

	return builder.String()
}

// parseProPublicaImage renders the first <img> in the element as markdown,
// followed by its <figcaption> in italics if present.
func parseProPublicaImage(e *colly.HTMLElement) string {
	img := e.DOM.Find("img").First()
	if img.Length() == 0 {
		return ""
	}

	src := img.AttrOr("src", "")
	if src == "" {
		src = img.AttrOr("data-src", "")
	}
	if src == "" || strings.HasPrefix(src, "data:") {
		return ""
	}

	builder := strings.Builder{}
	fmt.Fprintf(&builder, "![%s](%s)\n\n", img.AttrOr("alt", ""), src)

	caption := strings.TrimSpace(e.DOM.Find("figcaption").First().Text())
	if caption != "" {
		fmt.Fprintf(&builder, "*%s*\n\n", caption)
	}

	return builder.String()
}

// parseProPublicaList renders a <ul> or <ol> element as markdown.
func parseProPublicaList(e *colly.HTMLElement, ordered bool) string {
	builder := strings.Builder{}
	index := 0
	e.ForEach("li", func(_ int, li *colly.HTMLElement) {
		if !li.DOM.Parent().IsSelection(e.DOM) {
			return
		}
		index++

		bulletText := strings.TrimSpace(parseProPublicaInline(li))
		if ordered {
			fmt.Fprintf(&builder, "%d. %s\n", index, bulletText)
		} else {
			fmt.Fprintf(&builder, "* %s\n", bulletText)
		}
	})
	builder.WriteString("\n")
	return builder.String()
}

// parseProPublicaInline renders the inline content of an element as markdown,
// preserving <a> as markdown links, <strong>/<b> as bold, <em>/<i> as italic
// and <code> as backtick spans.
func parseProPublicaInline(e *colly.HTMLElement) string {
	return parseProPublicaInlineNodes(e.DOM.Contents())
}

func parseProPublicaInlineNodes(sel *goquery.Selection) string {
	builder := strings.Builder{}
	sel.Each(func(_ int, s *goquery.Selection) {
		node := s.Get(0)
		switch node.Type {
		case html.TextNode:
			builder.WriteString(node.Data)
		case html.ElementNode:
			switch node.Data {
			case "a":
				href, _ := s.Attr("href")
				linkText := parseProPublicaInlineNodes(s.Contents())
				if href != "" {
					fmt.Fprintf(&builder, "[%s](%s)", linkText, href)
				} else {
					builder.WriteString(linkText)
				}
			case "b", "strong":
				builder.WriteString(wrapProPublicaInline(parseProPublicaInlineNodes(s.Contents()), "**"))
			case "i", "em":
				builder.WriteString(wrapProPublicaInline(parseProPublicaInlineNodes(s.Contents()), "*"))
			case "code":
				fmt.Fprintf(&builder, "`%s`", s.Text())
			case "br":
				builder.WriteString("\n")
			case "script", "style", "noscript":
			default:
				builder.WriteString(parseProPublicaInlineNodes(s.Contents()))
			}
		}
	})
	return builder.String()
}

// wrapProPublicaInline wraps text with the marker, keeping leading and
// trailing whitespace outside the marker (e.g. "<em> word</em>" → " *word*").
func wrapProPublicaInline(text, marker string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return text
	}
	leading := text[:strings.Index(text, trimmed)]
	trailing := text[len(leading)+len(trimmed):]
	return leading + marker + trimmed + marker + trailing
}
