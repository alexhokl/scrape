package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

type UnseenJapanScraper struct {
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (u *UnseenJapanScraper) ScrapeArticle(url string) (string, error) {
	collector := colly.NewCollector()

	var markdown string

	// title
	collector.OnHTML("article.uj-article h1.uj-article__title", func(e *colly.HTMLElement) {
		markdown += fmt.Sprintf("# %s\n\n", strings.TrimSpace(e.Text))
	})

	// article body
	collector.OnHTML("article.uj-article div.uj-article__body", func(e *colly.HTMLElement) {
		markdown += parseUnseenJapanContent(e)
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func (u *UnseenJapanScraper) ScrapeTitle(url string) (string, error) {
	collector := colly.NewCollector()

	var title string

	collector.OnHTML("article.uj-article h1.uj-article__title", func(e *colly.HTMLElement) {
		title = strings.TrimSpace(e.Text)
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	return title, nil
}

func (u *UnseenJapanScraper) ScrapeFilename(url string) (string, error) {
	return getBasenameFromURL(url), nil
}

func parseUnseenJapanContent(e *colly.HTMLElement) string {
	builder := strings.Builder{}

	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		if !child.DOM.Parent().IsSelection(e.DOM) {
			return
		}

		switch child.Name {
		case "p":
			// skip the membership pitch injected into the article body
			if child.DOM.HasClass("uj-member-pitch") {
				return
			}
			text := strings.TrimSpace(parseUnseenJapanInline(child))
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
			// images are wrapped in <div class="wp-block-image"><figure>...</figure></div>
			if child.DOM.HasClass("wp-block-image") {
				builder.WriteString(parseUnseenJapanImage(child))
			}
		case "figure":
			builder.WriteString(parseUnseenJapanImage(child))
		case "ul":
			builder.WriteString(parseUnseenJapanList(child, false))
		case "ol":
			builder.WriteString(parseUnseenJapanList(child, true))
		case "blockquote":
			text := strings.TrimSpace(child.Text)
			if text != "" {
				fmt.Fprintf(&builder, "> %s\n\n", text)
			}
		}
	})

	return builder.String()
}

// parseUnseenJapanImage renders the first <img> in the element as markdown,
// followed by its <figcaption> in italics if present.
// Images are lazy-loaded, so the real URL is in data-src and src holds a
// data: placeholder.
func parseUnseenJapanImage(e *colly.HTMLElement) string {
	img := e.DOM.Find("img").First()
	if img.Length() == 0 {
		return ""
	}

	src := img.AttrOr("data-src", "")
	if src == "" {
		src = img.AttrOr("src", "")
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

// parseUnseenJapanList renders a <ul> or <ol> element as markdown.
func parseUnseenJapanList(e *colly.HTMLElement, ordered bool) string {
	builder := strings.Builder{}
	index := 0
	e.ForEach("li", func(_ int, li *colly.HTMLElement) {
		if !li.DOM.Parent().IsSelection(e.DOM) {
			return
		}
		index++

		bulletText := strings.TrimSpace(parseUnseenJapanInline(li))
		if ordered {
			fmt.Fprintf(&builder, "%d. %s\n", index, bulletText)
		} else {
			fmt.Fprintf(&builder, "* %s\n", bulletText)
		}
	})
	builder.WriteString("\n")
	return builder.String()
}

// parseUnseenJapanInline renders the inline content of an element as markdown,
// preserving <a> as markdown links, <strong>/<b> as bold, <em>/<i> as italic
// and <code> as backtick spans.
func parseUnseenJapanInline(e *colly.HTMLElement) string {
	return parseUnseenJapanInlineNodes(e.DOM.Contents())
}

func parseUnseenJapanInlineNodes(sel *goquery.Selection) string {
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
				linkText := parseUnseenJapanInlineNodes(s.Contents())
				if href != "" {
					fmt.Fprintf(&builder, "[%s](%s)", linkText, href)
				} else {
					builder.WriteString(linkText)
				}
			case "b", "strong":
				builder.WriteString(wrapUnseenJapanInline(parseUnseenJapanInlineNodes(s.Contents()), "**"))
			case "i", "em":
				builder.WriteString(wrapUnseenJapanInline(parseUnseenJapanInlineNodes(s.Contents()), "*"))
			case "code":
				fmt.Fprintf(&builder, "`%s`", s.Text())
			case "br":
				builder.WriteString("\n")
			case "script", "style", "noscript":
			default:
				builder.WriteString(parseUnseenJapanInlineNodes(s.Contents()))
			}
		}
	})
	return builder.String()
}

// wrapUnseenJapanInline wraps text with the marker, keeping leading and
// trailing whitespace outside the marker (e.g. "<em> Kaze</em>" → " *Kaze*").
func wrapUnseenJapanInline(text, marker string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return text
	}
	leading := text[:strings.Index(text, trimmed)]
	trailing := text[len(leading)+len(trimmed):]
	return leading + marker + trimmed + marker + trailing
}
