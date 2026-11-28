package scraper

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// siteConfig describes where the parts of an article, and the links of an
// index page, are on a site. Scrapers of sites with a regular layout delegate
// to it instead of walking the HTML themselves.
type siteConfig struct {
	// title selects the article title.
	title string
	// body selects the article content; the first match is used.
	body string
	// skip selects elements inside the body to leave out, such as the
	// title, metadata, share buttons and sign-up boxes.
	skip []string
	// links selects the links to articles on an index page. A selected
	// element that is not an <a> contributes the first <a> inside it.
	links string
	// linkText optionally selects the element inside each link holding its
	// text, for links that also wrap an excerpt or a date.
	linkText string
	// userAgent optionally replaces colly's default User-Agent, for sites
	// that reject it.
	userAgent string
}

// visit fetches the page at pageURL and calls fn with its document and URL.
func (s siteConfig) visit(pageURL string, fn func(doc *goquery.Selection, base *url.URL)) error {
	collector := colly.NewCollector()
	if s.userAgent != "" {
		collector.UserAgent = s.userAgent
	}
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		fn(e.DOM, e.Request.URL)
	})
	return collector.Visit(pageURL)
}

// article scrapes the page at pageURL and returns the title as a heading
// followed by the content as markdown.
func (s siteConfig) article(pageURL string) (string, error) {
	var title, content string
	found := false

	err := s.visit(pageURL, func(doc *goquery.Selection, base *url.URL) {
		title = s.findTitle(doc)
		body := doc.Find(s.body).First()
		if body.Length() == 0 {
			return
		}
		found = true
		content = htmlToMarkdown(body, base, s.skip)
	})
	if err != nil {
		return "", err
	}
	if !found {
		return "", fmt.Errorf("article content not found at %s", pageURL)
	}

	if title == "" {
		return content, nil
	}
	return fmt.Sprintf("# %s\n\n%s", title, content), nil
}

// scrapeTitle scrapes the article title of the page at pageURL.
func (s siteConfig) scrapeTitle(pageURL string) (string, error) {
	var title string
	err := s.visit(pageURL, func(doc *goquery.Selection, _ *url.URL) {
		title = s.findTitle(doc)
	})
	if err != nil {
		return "", err
	}
	if title == "" {
		return "", fmt.Errorf("title not found at %s", pageURL)
	}
	return title, nil
}

// findTitle returns the text of the title element, falling back to the
// og:title metadata.
func (s siteConfig) findTitle(doc *goquery.Selection) string {
	title := cleanText(doc.Find(s.title).First().Text())
	if title != "" {
		return title
	}
	return cleanText(doc.Find(`meta[property="og:title"]`).First().AttrOr("content", ""))
}

// scrapeLinks scrapes the article links of the index page at pageURL and
// returns a map of link text and absolute URL.
func (s siteConfig) scrapeLinks(pageURL string) (map[string]string, error) {
	links := make(map[string]string)
	seen := make(map[string]bool)

	err := s.visit(pageURL, func(doc *goquery.Selection, base *url.URL) {
		doc.Find(s.links).Each(func(_ int, sel *goquery.Selection) {
			anchor := sel
			if goquery.NodeName(sel) != "a" {
				anchor = sel.Find("a[href]").First()
			}
			href := strings.TrimSpace(anchor.AttrOr("href", ""))
			if href == "" {
				return
			}
			ref, err := url.Parse(href)
			if err != nil {
				return
			}
			ref.Fragment = ""
			absolute := base.ResolveReference(ref).String()

			textSel := sel
			if s.linkText != "" {
				if inner := sel.Find(s.linkText).First(); inner.Length() > 0 {
					textSel = inner
				}
			}
			text := cleanText(textSel.Text())
			if text == "" || seen[absolute] {
				return
			}
			seen[absolute] = true
			links[text] = absolute
		})
	})
	if err != nil {
		return nil, err
	}
	return links, nil
}

// filename returns the last path segment of pageURL without a .html
// extension, for use as a file name.
func (s siteConfig) filename(pageURL string) string {
	path := pageURL
	if parsed, err := url.Parse(pageURL); err == nil && parsed.Path != "" {
		path = parsed.Path
	}
	name := getBasenameFromURL(path)
	name = strings.TrimSuffix(name, ".html")
	name = strings.TrimSuffix(name, ".htm")
	if name == "" {
		return "article"
	}
	return name
}

// cleanText collapses whitespace and trims text taken from an element.
func cleanText(text string) string {
	return strings.TrimSpace(collapseWhitespace(text))
}
