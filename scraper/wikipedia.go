package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

type WikipediaScraper struct {
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (w *WikipediaScraper) ScrapeArticle(url string) (string, error) {
	collector := colly.NewCollector()

	var markdown string

	// title — h1#firstHeading is outside div.mw-parser-output
	collector.OnHTML("h1#firstHeading", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		if text != "" {
			markdown += fmt.Sprintf("# %s\n\n", text)
		}
	})

	// article body
	collector.OnHTML("div#mw-content-text > div.mw-parser-output", func(e *colly.HTMLElement) {
		markdown += parseWikipediaContent(e)
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func (w *WikipediaScraper) ScrapeTitle(url string) (string, error) {
	collector := colly.NewCollector()

	var title string

	collector.OnHTML("h1#firstHeading", func(e *colly.HTMLElement) {
		title = strings.TrimSpace(e.Text)
	})

	err := collector.Visit(url)
	if err != nil {
		return "", err
	}

	return title, nil
}

func (w *WikipediaScraper) ScrapeFilename(url string) (string, error) {
	return getBasenameFromURL(url), nil
}

// wikipediaSkipSections contains section headings whose content should be
// excluded from the scraped output.
var wikipediaSkipSections = map[string]bool{
	"References":       true,
	"Further reading":  true,
	"External links":   true,
	"See also":         true,
	"Notes":            true,
	"Citations":        true,
	"Bibliography":     true,
	"Sources":          true,
	"Cited sources":    true,
	"Works cited":      true,
	"Cited references": true,
}

func parseWikipediaContent(e *colly.HTMLElement) string {
	builder := strings.Builder{}

	// Track whether we are inside a section that should be skipped.
	skipping := false

	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		if !child.DOM.Parent().IsSelection(e.DOM) {
			return
		}

		// Section heading wrappers: <div class="mw-heading mw-heading2">
		if child.Name == "div" && child.DOM.HasClass("mw-heading") {
			heading := parseWikipediaHeading(child)
			if heading == "" {
				return
			}

			// Check if this section should be skipped.
			headingText := extractWikipediaHeadingText(child)
			if wikipediaSkipSections[headingText] {
				skipping = true
				return
			}

			// A non-skipped heading resets skipping for mw-heading2 level.
			if child.DOM.HasClass("mw-heading2") {
				skipping = false
			}

			if !skipping {
				builder.WriteString(heading)
			}
			return
		}

		if skipping {
			return
		}

		switch child.Name {
		case "p":
			text := strings.TrimSpace(parseWikipediaInline(child))
			if text != "" {
				builder.WriteString(fmt.Sprintf("%s\n\n", text))
			}
		case "ul":
			// Skip portal boxes and other navigation-related lists.
			if child.DOM.HasClass("portalbox") {
				return
			}
			builder.WriteString(parseWikipediaList(child, false))
		case "ol":
			// Skip reference lists.
			if child.DOM.HasClass("references") {
				return
			}
			builder.WriteString(parseWikipediaList(child, true))
		case "pre":
			code := strings.TrimRight(child.Text, "\n")
			builder.WriteString(fmt.Sprintf("```\n%s\n```\n\n", code))
		case "blockquote":
			text := strings.TrimSpace(child.Text)
			if text != "" {
				builder.WriteString(fmt.Sprintf("> %s\n\n", text))
			}
		case "figure":
			parseWikipediaFigure(child, &builder)
		case "table":
			if child.DOM.HasClass("wikitable") {
				builder.WriteString(parseWikipediaTable(child))
			}
		}
	})

	return builder.String()
}

// parseWikipediaHeading renders a mw-heading div as a markdown heading.
// It returns an empty string when no heading element is found inside the div.
func parseWikipediaHeading(div *colly.HTMLElement) string {
	var result string
	div.ForEach("h2, h3, h4, h5, h6", func(_ int, h *colly.HTMLElement) {
		if result != "" {
			return
		}
		text := wikipediaHeadingText(h)
		if text == "" {
			return
		}
		switch h.Name {
		case "h2":
			result = fmt.Sprintf("## %s\n\n", text)
		case "h3":
			result = fmt.Sprintf("### %s\n\n", text)
		case "h4":
			result = fmt.Sprintf("#### %s\n\n", text)
		case "h5":
			result = fmt.Sprintf("##### %s\n\n", text)
		case "h6":
			result = fmt.Sprintf("###### %s\n\n", text)
		}
	})
	return result
}

// extractWikipediaHeadingText returns the plain text of the heading element
// inside a mw-heading div, stripping edit-section links.
func extractWikipediaHeadingText(div *colly.HTMLElement) string {
	var text string
	div.ForEach("h2, h3, h4, h5, h6", func(_ int, h *colly.HTMLElement) {
		if text != "" {
			return
		}
		text = wikipediaHeadingText(h)
	})
	return text
}

// wikipediaHeadingText extracts the text of a heading element, excluding
// mw-editsection spans.
func wikipediaHeadingText(h *colly.HTMLElement) string {
	clone := h.DOM.Clone()
	clone.Find(".mw-editsection").Remove()
	return strings.TrimSpace(clone.Text())
}

// parseWikipediaFigure extracts an image from a Wikipedia <figure> element.
func parseWikipediaFigure(e *colly.HTMLElement, builder *strings.Builder) {
	e.ForEach("img", func(_ int, img *colly.HTMLElement) {
		src := img.Attr("src")
		alt := img.Attr("alt")
		if src == "" {
			return
		}
		// Wikipedia uses protocol-relative URLs (//upload.wikimedia.org/...)
		if strings.HasPrefix(src, "//") {
			src = "https:" + src
		}
		builder.WriteString(fmt.Sprintf("![%s](%s)\n\n", alt, src))
	})
}

// parseWikipediaList renders a <ul> or <ol> element as markdown.
func parseWikipediaList(e *colly.HTMLElement, ordered bool) string {
	builder := strings.Builder{}
	index := 0
	e.ForEach("li", func(_ int, li *colly.HTMLElement) {
		if !li.DOM.Parent().IsSelection(e.DOM) {
			return
		}
		index++

		text := strings.TrimSpace(parseWikipediaInline(li))
		if ordered {
			builder.WriteString(fmt.Sprintf("%d. %s\n", index, text))
		} else {
			builder.WriteString(fmt.Sprintf("* %s\n", text))
		}
	})
	builder.WriteString("\n")
	return builder.String()
}

// parseWikipediaInline renders the inline content of an element as markdown,
// stripping citation superscripts and edit-section links while preserving
// bold, italic, code, and anchor elements.
func parseWikipediaInline(e *colly.HTMLElement) string {
	builder := strings.Builder{}
	for _, node := range e.DOM.Contents().Nodes {
		switch node.Type {
		case html.TextNode:
			builder.WriteString(node.Data)
		case html.ElementNode:
			sel := e.DOM.FindNodes(node)

			// Skip elements that should be excluded.
			if shouldSkipWikipediaNode(node, sel) {
				continue
			}

			switch node.Data {
			case "code":
				builder.WriteString(fmt.Sprintf("`%s`", sel.Text()))
			case "a":
				href, _ := sel.Attr("href")
				linkText := parseWikipediaInlineNodes(sel.Contents().Nodes)
				linkText = strings.TrimSpace(linkText)
				if linkText == "" {
					continue
				}
				if href != "" {
					// Convert relative wiki links to absolute URLs.
					if strings.HasPrefix(href, "/wiki/") {
						href = "https://en.wikipedia.org" + href
					}
					builder.WriteString(fmt.Sprintf("[%s](%s)", linkText, href))
				} else {
					builder.WriteString(linkText)
				}
			case "b", "strong":
				text := parseWikipediaInlineNodes(sel.Contents().Nodes)
				trimmed := strings.TrimSpace(text)
				if trimmed != "" {
					builder.WriteString(fmt.Sprintf("**%s**", trimmed))
				}
			case "i", "em":
				text := parseWikipediaInlineNodes(sel.Contents().Nodes)
				trimmed := strings.TrimSpace(text)
				if trimmed != "" {
					builder.WriteString(fmt.Sprintf("*%s*", trimmed))
				}
			case "span":
				// Skip mw-editsection spans, shortdescription spans, and
				// metadata spans. For other spans, emit their text.
				if sel.HasClass("mw-editsection") ||
					sel.HasClass("shortdescription") ||
					sel.HasClass("Z3988") {
					continue
				}
				builder.WriteString(parseWikipediaInlineNodes(sel.Contents().Nodes))
			default:
				builder.WriteString(sel.Text())
			}
		}
	}
	return builder.String()
}

// shouldSkipWikipediaNode returns true for elements that should be excluded
// from the markdown output.
func shouldSkipWikipediaNode(node *html.Node, sel *goquery.Selection) bool {
	switch node.Data {
	case "sup":
		// Skip citation references [1], [2] etc and "citation needed" tags.
		if sel.HasClass("reference") || sel.HasClass("noprint") {
			return true
		}
	case "style":
		return true
	case "span":
		if sel.HasClass("mw-editsection") ||
			sel.HasClass("shortdescription") ||
			sel.HasClass("Z3988") {
			return true
		}
	case "div":
		// Skip reference wrappers, navboxes, and other metadata divs.
		if sel.HasClass("reflist") ||
			sel.HasClass("mw-references-wrap") ||
			sel.HasClass("navbox") ||
			sel.HasClass("navbox-styles") ||
			sel.HasClass("shortdescription") ||
			sel.HasClass("spoken-wikipedia") {
			return true
		}
	}
	return false
}

// parseWikipediaInlineNodes renders a slice of HTML nodes as inline markdown,
// used for content inside elements like <a> where we still want inline
// formatting preserved.
func parseWikipediaInlineNodes(nodes []*html.Node) string {
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
			case "b", "strong":
				var text strings.Builder
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						text.WriteString(c.Data)
					}
				}
				trimmed := strings.TrimSpace(text.String())
				if trimmed != "" {
					builder.WriteString(fmt.Sprintf("**%s**", trimmed))
				}
			case "i", "em":
				var text strings.Builder
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						text.WriteString(c.Data)
					}
				}
				trimmed := strings.TrimSpace(text.String())
				if trimmed != "" {
					builder.WriteString(fmt.Sprintf("*%s*", trimmed))
				}
			case "sup":
				// Skip citation superscripts inside links.
				hasRefClass := false
				for _, attr := range node.Attr {
					if attr.Key == "class" && (strings.Contains(attr.Val, "reference") || strings.Contains(attr.Val, "noprint")) {
						hasRefClass = true
						break
					}
				}
				if hasRefClass {
					continue
				}
				// Emit non-reference superscripts as plain text.
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

// parseWikipediaTable renders a Wikipedia table with class "wikitable" as
// a markdown table.
func parseWikipediaTable(table *colly.HTMLElement) string {
	builder := strings.Builder{}
	rowIndex := 0

	table.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
		cells := make([]string, 0)
		isHeader := false

		tr.ForEach("th, td", func(_ int, cell *colly.HTMLElement) {
			if cell.Name == "th" {
				isHeader = true
			}
			text := strings.TrimSpace(cell.Text)
			// Replace newlines within cells with spaces.
			text = strings.ReplaceAll(text, "\n", " ")
			cells = append(cells, text)
		})

		if len(cells) == 0 {
			return
		}

		builder.WriteString("| ")
		builder.WriteString(strings.Join(cells, " | "))
		builder.WriteString(" |\n")

		if isHeader || rowIndex == 0 {
			builder.WriteString("|")
			for range cells {
				builder.WriteString(" --- |")
			}
			builder.WriteString("\n")
		}

		rowIndex++
	})

	builder.WriteString("\n")
	return builder.String()
}
