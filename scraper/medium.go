package scraper

import (
	"encoding/xml"
	"fmt"
	"net/url"
	pathpkg "path"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// MediumScraper scrapes Medium (medium.com) posts and post lists through the
// RSS feeds of authors and publications, as Medium answers plain HTTP
// clients of its pages with a bot challenge. A feed only holds the recent
// posts, and posts behind the paywall have no content in it.
type MediumScraper struct {
}

// mediumSkip selects the tracking pixel at the end of the content of a feed
// item.
var mediumSkip = []string{
	`img[src*="medium.com/_/stat"]`,
}

type mediumFeed struct {
	Items []mediumItem `xml:"channel>item"`
}

type mediumItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Content string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
}

// mediumFeedURL returns the URL of the RSS feed holding the post or listing
// the posts at pageURL.
func mediumFeedURL(pageURL string) (string, error) {
	parsed, err := url.Parse(pageURL)
	if err != nil {
		return "", fmt.Errorf("error parsing URL %s: %w", pageURL, err)
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("invalid URL: %s", pageURL)
	}
	if parsed.Path == "/feed" || strings.HasPrefix(parsed.Path, "/feed/") {
		return pageURL, nil
	}

	base := url.URL{Scheme: parsed.Scheme, Host: parsed.Host, Path: "/feed"}
	host := strings.ToLower(strings.TrimPrefix(parsed.Hostname(), "www."))
	if host == "medium.com" {
		segment, _, _ := strings.Cut(strings.Trim(parsed.Path, "/"), "/")
		if segment == "" {
			return "", fmt.Errorf("URL %s names no author or publication", pageURL)
		}
		base.Path = "/feed/" + segment
	}
	return base.String(), nil
}

// mediumPostID returns the ID at the end of the slug of a post URL, which
// stays the same whatever host or query the URL has.
func mediumPostID(postURL string) string {
	path := postURL
	if parsed, err := url.Parse(postURL); err == nil {
		path = parsed.Path
	}
	slug := pathpkg.Base(strings.TrimRight(path, "/"))
	index := strings.LastIndex(slug, "-")
	if index < 0 {
		return slug
	}
	return slug[index+1:]
}

// stripQuery returns rawURL without its query and fragment.
func stripQuery(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String()
}

func fetchMediumFeed(feedURL string) (*mediumFeed, error) {
	var body []byte
	collector := colly.NewCollector()
	collector.OnResponse(func(r *colly.Response) {
		body = r.Body
	})
	if err := collector.Visit(feedURL); err != nil {
		return nil, fmt.Errorf("error fetching feed %s: %w", feedURL, err)
	}

	var feed mediumFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("error parsing feed %s: %w", feedURL, err)
	}
	return &feed, nil
}

// findItem returns the feed item of the post at postURL.
func (s *MediumScraper) findItem(postURL string) (*mediumItem, error) {
	feedURL, err := mediumFeedURL(postURL)
	if err != nil {
		return nil, err
	}
	feed, err := fetchMediumFeed(feedURL)
	if err != nil {
		return nil, err
	}
	id := mediumPostID(postURL)
	for i := range feed.Items {
		if id != "" && mediumPostID(feed.Items[i].Link) == id {
			return &feed.Items[i], nil
		}
	}
	return nil, fmt.Errorf("post %s not found in feed %s (feeds hold recent posts only)", postURL, feedURL)
}

// ScrapeLinks scrapes the posts in the feed of the author or publication at
// url and returns a map of post title and URL
func (s *MediumScraper) ScrapeLinks(url string) (map[string]string, error) {
	feedURL, err := mediumFeedURL(url)
	if err != nil {
		return nil, err
	}
	feed, err := fetchMediumFeed(feedURL)
	if err != nil {
		return nil, err
	}
	links := make(map[string]string)
	for _, item := range feed.Items {
		title := cleanText(item.Title)
		link := strings.TrimSpace(item.Link)
		if title == "" || link == "" {
			continue
		}
		links[title] = stripQuery(link)
	}
	return links, nil
}

// ScrapeArticle scrapes the post at url from its feed and returns markdown
// in a string
func (s *MediumScraper) ScrapeArticle(url string) (string, error) {
	item, err := s.findItem(url)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(item.Content) == "" {
		return "", fmt.Errorf("post %s has no content in its feed; it may be behind the paywall", url)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(item.Content))
	if err != nil {
		return "", fmt.Errorf("error parsing content of %s: %w", url, err)
	}
	base, err := parseURL(stripQuery(item.Link))
	if err != nil {
		return "", err
	}
	content := htmlToMarkdown(doc.Find("body"), base, mediumSkip)

	title := cleanText(item.Title)
	if title == "" {
		return content, nil
	}
	return fmt.Sprintf("# %s\n\n%s", title, content), nil
}

func (s *MediumScraper) ScrapeTitle(url string) (string, error) {
	item, err := s.findItem(url)
	if err != nil {
		return "", err
	}
	title := cleanText(item.Title)
	if title == "" {
		return "", fmt.Errorf("title not found at %s", url)
	}
	return title, nil
}

func (s *MediumScraper) ScrapeFilename(url string) (string, error) {
	title, err := s.ScrapeTitle(url)
	if err != nil {
		return "", err
	}
	return generateFileNameFromTitle(title), nil
}

func parseURL(rawURL string) (*url.URL, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL %s: %w", rawURL, err)
	}
	return parsed, nil
}
