package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

type JapanTimesScraper struct {
}

const japanTimesTitleSelector = "div.jt-top-article-details h1.title-article"

// subscriber-only articles render only the first few paragraphs, with the
// body marked with this class
const japanTimesPaywallClass = "blurred-text"

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (j *JapanTimesScraper) ScrapeArticle(url string) (string, error) {
	collector := colly.NewCollector()

	var title string
	var subtitle string
	var image string
	var body string
	var truncated bool

	collector.OnHTML(japanTimesTitleSelector, func(e *colly.HTMLElement) {
		if title == "" {
			title = strings.TrimSpace(e.Text)
		}
	})

	collector.OnHTML("div.jt-top-article-details h3.subtitle-article", func(e *colly.HTMLElement) {
		if subtitle == "" {
			subtitle = strings.TrimSpace(e.Text)
		}
	})

	// the gallery may hold several images; only the leading one is used
	collector.OnHTML("div.jt-article-details div.article-image-gallery", func(e *colly.HTMLElement) {
		if image == "" {
			image = parseJapanTimesImage(e)
		}
	})

	// article body; related articles, the newsletter sign-up and the
	// subscription pitch are outside of #jtarticle
	collector.OnHTML("div#jtarticle div.article-body", func(e *colly.HTMLElement) {
		body += parseJapanTimesContent(e)
		if e.DOM.HasClass(japanTimesPaywallClass) {
			truncated = true
		}
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	if title != "" {
		fmt.Fprintf(&builder, "# %s\n\n", title)
	}
	if subtitle != "" {
		fmt.Fprintf(&builder, "*%s*\n\n", subtitle)
	}
	builder.WriteString(image)
	builder.WriteString(body)
	if truncated {
		builder.WriteString("*The rest of this article is available to subscribers only.*\n\n")
	}

	return builder.String(), nil
}

func (j *JapanTimesScraper) ScrapeTitle(url string) (string, error) {
	collector := colly.NewCollector()

	var title string

	collector.OnHTML(japanTimesTitleSelector, func(e *colly.HTMLElement) {
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

func (j *JapanTimesScraper) ScrapeFilename(url string) (string, error) {
	return getBasenameFromURL(url), nil
}

func parseJapanTimesContent(e *colly.HTMLElement) string {
	builder := strings.Builder{}

	// ad slots (<div>) are skipped; as they are placed inside <p> in the
	// source, the parser closes the paragraph early and leaves an empty <p>
	// behind
	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		if !child.DOM.Parent().IsSelection(e.DOM) {
			return
		}

		switch child.Name {
		case "p":
			text := strings.TrimSpace(parseJapanTimesInline(child))
			if text != "" {
				fmt.Fprintf(&builder, "%s\n\n", text)
			}
		case "h2":
			fmt.Fprintf(&builder, "## %s\n\n", strings.TrimSpace(child.Text))
		case "h3":
			fmt.Fprintf(&builder, "### %s\n\n", strings.TrimSpace(child.Text))
		case "h4":
			fmt.Fprintf(&builder, "#### %s\n\n", strings.TrimSpace(child.Text))
		case "ul":
			builder.WriteString(parseJapanTimesList(child, false))
		case "ol":
			builder.WriteString(parseJapanTimesList(child, true))
		case "blockquote":
			text := strings.TrimSpace(child.Text)
			if text != "" {
				fmt.Fprintf(&builder, "> %s\n\n", text)
			}
		}
	})

	return builder.String()
}

// parseJapanTimesImage renders the first image of the gallery as markdown,
// linking the original upload rather than the resized CDN copy, followed by
// its caption (with credit) in italics if present.
func parseJapanTimesImage(e *colly.HTMLElement) string {
	link := e.DOM.Find("a.fancybox").First()
	src := link.AttrOr("href", "")
	if src == "" {
		src = link.Find("img").First().AttrOr("src", "")
	}
	if src == "" || strings.HasPrefix(src, "data:") {
		return ""
	}

	builder := strings.Builder{}
	fmt.Fprintf(&builder, "![](%s)\n\n", src)

	caption := strings.Join(strings.Fields(link.AttrOr("data-caption", "")), " ")
	if caption != "" {
		fmt.Fprintf(&builder, "*%s*\n\n", caption)
	}

	return builder.String()
}

// parseJapanTimesList renders a <ul> or <ol> element as markdown.
func parseJapanTimesList(e *colly.HTMLElement, ordered bool) string {
	builder := strings.Builder{}
	index := 0
	e.ForEach("li", func(_ int, li *colly.HTMLElement) {
		if !li.DOM.Parent().IsSelection(e.DOM) {
			return
		}
		index++

		bulletText := strings.TrimSpace(parseJapanTimesInline(li))
		if ordered {
			fmt.Fprintf(&builder, "%d. %s\n", index, bulletText)
		} else {
			fmt.Fprintf(&builder, "* %s\n", bulletText)
		}
	})
	builder.WriteString("\n")
	return builder.String()
}

// parseJapanTimesInline renders the inline content of an element as markdown,
// preserving <a> as markdown links, <strong>/<b> as bold, <em>/<i> as italic
// and <code> as backtick spans.
func parseJapanTimesInline(e *colly.HTMLElement) string {
	return parseJapanTimesInlineNodes(e.DOM.Contents())
}

func parseJapanTimesInlineNodes(sel *goquery.Selection) string {
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
				linkText := parseJapanTimesInlineNodes(s.Contents())
				if href != "" {
					fmt.Fprintf(&builder, "[%s](%s)", linkText, href)
				} else {
					builder.WriteString(linkText)
				}
			case "b", "strong":
				builder.WriteString(wrapJapanTimesInline(parseJapanTimesInlineNodes(s.Contents()), "**"))
			case "i", "em":
				builder.WriteString(wrapJapanTimesInline(parseJapanTimesInlineNodes(s.Contents()), "*"))
			case "code":
				fmt.Fprintf(&builder, "`%s`", s.Text())
			case "br":
				builder.WriteString("\n")
			case "script", "style", "noscript", "div":
			default:
				builder.WriteString(parseJapanTimesInlineNodes(s.Contents()))
			}
		}
	})
	return builder.String()
}

// wrapJapanTimesInline wraps text with the marker, keeping leading and
// trailing whitespace outside the marker (e.g. "<em> word</em>" → " *word*").
func wrapJapanTimesInline(text, marker string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return text
	}
	leading := text[:strings.Index(text, trimmed)]
	trailing := text[len(leading)+len(trimmed):]
	return leading + marker + trimmed + marker + trailing
}
