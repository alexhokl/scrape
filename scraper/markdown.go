package scraper

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// alwaysSkipped lists elements that never carry article text.
const alwaysSkipped = "script, style, noscript, template, iframe, svg, button, form, input, select, textarea"

// blockElements lists elements rendered as blocks rather than inline text.
var blockElements = map[string]bool{
	"p": true, "h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"ul": true, "ol": true, "li": true, "pre": true, "blockquote": true, "figure": true,
	"table": true, "hr": true, "div": true, "section": true, "article": true, "main": true,
	"header": true, "footer": true, "aside": true, "nav": true, "dl": true, "dt": true,
	"dd": true, "details": true, "summary": true, "figcaption": true, "center": true,
}

var (
	whitespacePattern     = regexp.MustCompile(`[ \t\r\n\f]+`)
	excessNewlinesPattern = regexp.MustCompile(`\n{3,}`)
	codeLanguagePattern   = regexp.MustCompile(`(?:^|\s)(?:language|lang)-([A-Za-z0-9_+#.-]+)`)
)

// markdownConverter renders HTML as markdown, resolving relative links and
// image sources against base.
type markdownConverter struct {
	base *url.URL
}

// htmlToMarkdown renders the content of sel as markdown. Elements matching
// any of the skip selectors are left out. Relative links and image sources
// are resolved against base, which may be nil.
func htmlToMarkdown(sel *goquery.Selection, base *url.URL, skip []string) string {
	body := sel.Clone()
	body.Find(alwaysSkipped).Remove()
	if len(skip) > 0 {
		// select every match before removing any, so that selectors such
		// as "h1 + p" still match after the h1 is removed
		body.Find(strings.Join(skip, ", ")).Remove()
	}

	c := &markdownConverter{base: base}
	markdown := c.blocks(body)
	markdown = excessNewlinesPattern.ReplaceAllString(markdown, "\n\n")
	markdown = strings.TrimSpace(markdown)
	if markdown == "" {
		return ""
	}
	return markdown + "\n"
}

// blocks renders the children of sel as markdown blocks. Consecutive inline
// children are gathered into one paragraph.
func (c *markdownConverter) blocks(sel *goquery.Selection) string {
	builder := strings.Builder{}
	pending := strings.Builder{}

	flush := func() {
		text := normaliseInline(pending.String())
		if text != "" {
			fmt.Fprintf(&builder, "%s\n\n", text)
		}
		pending.Reset()
	}

	sel.Contents().Each(func(_ int, s *goquery.Selection) {
		node := s.Get(0)
		switch node.Type {
		case html.TextNode:
			pending.WriteString(collapseWhitespace(node.Data))
		case html.ElementNode:
			if !c.isBlock(s) {
				pending.WriteString(c.inlineElement(s))
				return
			}
			flush()
			builder.WriteString(c.block(s))
		}
	})
	flush()

	return builder.String()
}

// isBlock reports whether the element is rendered as a block. An inline
// element (such as a span or link) wrapping block content is treated as a
// block so that the content inside it keeps its structure.
func (c *markdownConverter) isBlock(s *goquery.Selection) bool {
	name := goquery.NodeName(s)
	if blockElements[name] {
		return true
	}
	if name == "img" || name == "br" {
		return false
	}
	return s.Find("p, div, ul, ol, pre, table, blockquote, h1, h2, h3, h4, h5, h6, figure").Length() > 0
}

// block renders a single block element.
func (c *markdownConverter) block(s *goquery.Selection) string {
	name := goquery.NodeName(s)
	switch name {
	case "p", "dt", "summary", "figcaption":
		text := normaliseInline(c.inline(s))
		if text == "" {
			return ""
		}
		if name == "dt" || name == "summary" {
			text = wrapInline(text, "**")
		}
		if name == "figcaption" {
			text = wrapInline(text, "*")
		}
		return text + "\n\n"
	case "h1", "h2", "h3", "h4", "h5", "h6":
		text := strings.ReplaceAll(normaliseInline(c.inline(s)), "\n", " ")
		if text == "" {
			return ""
		}
		level := int(name[1] - '0')
		return fmt.Sprintf("%s %s\n\n", strings.Repeat("#", level), text)
	case "ul":
		return c.list(s, false)
	case "ol":
		return c.list(s, true)
	case "pre":
		return c.codeBlock(s)
	case "blockquote":
		return c.blockquote(s)
	case "figure":
		return c.figure(s)
	case "table":
		return c.table(s)
	case "hr":
		return "---\n\n"
	default:
		return c.blocks(s)
	}
}

// list renders a <ul> or <ol>. Each item is rendered as blocks, so nested
// lists, paragraphs and code inside an item are kept and indented under it.
func (c *markdownConverter) list(s *goquery.Selection, ordered bool) string {
	builder := strings.Builder{}
	index := 0
	s.Children().Each(func(_ int, item *goquery.Selection) {
		if goquery.NodeName(item) != "li" {
			return
		}
		index++

		marker := "* "
		if ordered {
			marker = fmt.Sprintf("%d. ", index)
		}
		indent := strings.Repeat(" ", len(marker))

		content := strings.TrimSpace(c.blocks(item))
		// keep the list tight: paragraphs in an item become lines
		content = strings.ReplaceAll(content, "\n\n", "\n")
		lines := strings.Split(content, "\n")
		for i := 1; i < len(lines); i++ {
			if lines[i] != "" {
				lines[i] = indent + lines[i]
			}
		}
		fmt.Fprintf(&builder, "%s%s\n", marker, strings.Join(lines, "\n"))
	})
	if index == 0 {
		return ""
	}
	builder.WriteString("\n")
	return builder.String()
}

// codeBlock renders a <pre> as a fenced code block, taking the language from
// a language-* or lang-* class on the <pre> or its <code>.
func (c *markdownConverter) codeBlock(s *goquery.Selection) string {
	code := strings.Trim(s.Text(), "\n")
	if strings.TrimSpace(code) == "" {
		return ""
	}

	language := codeLanguage(s.AttrOr("class", ""))
	if language == "" {
		language = codeLanguage(s.Find("code").First().AttrOr("class", ""))
	}

	fence := "```"
	for strings.Contains(code, fence) {
		fence += "`"
	}
	return fmt.Sprintf("%s%s\n%s\n%s\n\n", fence, language, code, fence)
}

func codeLanguage(class string) string {
	match := codeLanguagePattern.FindStringSubmatch(class)
	if match == nil {
		return ""
	}
	return match[1]
}

// blockquote renders a <blockquote>, prefixing every line with "> ".
func (c *markdownConverter) blockquote(s *goquery.Selection) string {
	content := strings.TrimSpace(c.blocks(s))
	if content == "" {
		return ""
	}
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if line == "" {
			lines[i] = ">"
		} else {
			lines[i] = "> " + line
		}
	}
	return strings.Join(lines, "\n") + "\n\n"
}

// figure renders the images of a <figure>, followed by its caption in
// italics. A figure without images (such as a code card) is rendered as
// blocks.
func (c *markdownConverter) figure(s *goquery.Selection) string {
	images := s.Find("img")
	if images.Length() == 0 {
		return c.blocks(s)
	}

	builder := strings.Builder{}
	images.Each(func(_ int, img *goquery.Selection) {
		if markdown := c.image(img); markdown != "" {
			fmt.Fprintf(&builder, "%s\n\n", markdown)
		}
	})

	caption := normaliseInline(c.inline(s.Find("figcaption").First()))
	if caption != "" {
		fmt.Fprintf(&builder, "%s\n\n", wrapInline(caption, "*"))
	}
	return builder.String()
}

// table renders a <table> as a markdown table, using the first row as the
// header.
func (c *markdownConverter) table(s *goquery.Selection) string {
	var rows [][]string
	columns := 0
	s.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		if tr.Closest("table").Get(0) != s.Get(0) {
			return // a row of a nested table
		}
		var cells []string
		tr.Children().Each(func(_ int, cell *goquery.Selection) {
			name := goquery.NodeName(cell)
			if name != "td" && name != "th" {
				return
			}
			text := normaliseInline(c.inline(cell))
			text = strings.ReplaceAll(text, "\n", " ")
			text = strings.ReplaceAll(text, "|", `\|`)
			cells = append(cells, text)
		})
		if len(cells) > 0 {
			rows = append(rows, cells)
			columns = max(columns, len(cells))
		}
	})
	if len(rows) == 0 {
		return ""
	}

	builder := strings.Builder{}
	writeRow := func(cells []string) {
		for len(cells) < columns {
			cells = append(cells, "")
		}
		fmt.Fprintf(&builder, "| %s |\n", strings.Join(cells, " | "))
	}
	writeRow(rows[0])
	separator := make([]string, columns)
	for i := range separator {
		separator[i] = "---"
	}
	writeRow(separator)
	for _, row := range rows[1:] {
		writeRow(row)
	}
	builder.WriteString("\n")
	return builder.String()
}

// inline renders the content of s as inline markdown, preserving links,
// bold, italic, strikethrough, inline code, line breaks and images.
func (c *markdownConverter) inline(s *goquery.Selection) string {
	builder := strings.Builder{}
	s.Contents().Each(func(_ int, child *goquery.Selection) {
		node := child.Get(0)
		switch node.Type {
		case html.TextNode:
			builder.WriteString(collapseWhitespace(node.Data))
		case html.ElementNode:
			builder.WriteString(c.inlineElement(child))
		}
	})
	return builder.String()
}

func (c *markdownConverter) inlineElement(s *goquery.Selection) string {
	switch goquery.NodeName(s) {
	case "a":
		text := c.inline(s)
		href := strings.TrimSpace(s.AttrOr("href", ""))
		if strings.TrimSpace(text) == "" {
			return text
		}
		// fragment links (such as footnote references) and script links
		// lead nowhere outside the page
		if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "javascript:") {
			return text
		}
		return fmt.Sprintf("[%s](%s)", strings.TrimSpace(text), c.resolve(href))
	case "b", "strong":
		return wrapInline(c.inline(s), "**")
	case "i", "em":
		return wrapInline(c.inline(s), "*")
	case "del", "s", "strike":
		return wrapInline(c.inline(s), "~~")
	case "code", "kbd", "samp", "tt":
		code := collapseWhitespace(s.Text())
		if strings.TrimSpace(code) == "" {
			return code
		}
		fence := "`"
		for strings.Contains(code, fence) {
			fence += "`"
		}
		if strings.HasPrefix(code, "`") || strings.HasSuffix(code, "`") {
			code = " " + code + " "
		}
		return fence + code + fence
	case "br":
		return "\n"
	case "img":
		return c.image(s)
	default:
		return c.inline(s)
	}
}

// image renders an <img> as markdown. Lazy-loaded images keep the real
// source in data-src or data-lazy-src, with src holding a placeholder.
func (c *markdownConverter) image(img *goquery.Selection) string {
	src := ""
	for _, attr := range []string{"data-src", "data-lazy-src", "data-original", "src"} {
		value := strings.TrimSpace(img.AttrOr(attr, ""))
		if value != "" && !strings.HasPrefix(value, "data:") {
			src = value
			break
		}
	}
	if src == "" {
		return ""
	}
	alt := strings.TrimSpace(collapseWhitespace(img.AttrOr("alt", "")))
	return fmt.Sprintf("![%s](%s)", alt, c.resolve(src))
}

// resolve returns ref as an absolute URL when a base URL is known.
func (c *markdownConverter) resolve(ref string) string {
	if c.base == nil {
		return ref
	}
	parsed, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	return c.base.ResolveReference(parsed).String()
}

// collapseWhitespace replaces each run of HTML whitespace with one space,
// as a browser does when rendering text.
func collapseWhitespace(text string) string {
	return whitespacePattern.ReplaceAllString(text, " ")
}

// normaliseInline trims the spaces around the line breaks and at both ends
// of rendered inline text.
func normaliseInline(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// wrapInline wraps text with the marker, keeping leading and trailing
// whitespace outside the marker (e.g. "<em> Kaze</em>" → " *Kaze*").
func wrapInline(text, marker string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return text
	}
	leading := text[:strings.Index(text, trimmed)]
	trailing := text[len(leading)+len(trimmed):]
	return leading + marker + trimmed + marker + trailing
}
