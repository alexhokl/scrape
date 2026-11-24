package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var _ ArticleScraper = (*MediumScraper)(nil)
var _ LinkScraper = (*MediumScraper)(nil)

// mediumFeedXML is trimmed from https://skamille.medium.com/feed
const mediumFeedXML = `<?xml version="1.0" encoding="UTF-8"?><rss xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
    <channel>
        <title><![CDATA[Stories by Camille Fournier on Medium]]></title>
        <link>https://medium.com/@skamille?source=rss-12421a4f856------2</link>
        <image>
            <url>https://cdn-images-1.medium.com/fit/c/150/150/0*example.jpeg</url>
            <title>Stories by Camille Fournier on Medium</title>
            <link>https://medium.com/@skamille?source=rss-12421a4f856------2</link>
        </image>
        <item>
            <title><![CDATA[A sidebar on Code Review]]></title>
            <link>https://skamille.medium.com/a-sidebar-on-code-review-e75438fdfa5d?source=rss-12421a4f856------2</link>
            <guid isPermaLink="false">https://medium.com/p/e75438fdfa5d</guid>
            <content:encoded><![CDATA[<p>Code review notes.</p>]]></content:encoded>
        </item>
        <item>
            <title><![CDATA[The Manager’s Path in the Age of AI]]></title>
            <link>https://skamille.medium.com/the-managers-path-in-the-age-of-ai-279cb6611d66?source=rss-12421a4f856------2</link>
            <guid isPermaLink="false">https://medium.com/p/279cb6611d66</guid>
            <category><![CDATA[management]]></category>
            <dc:creator><![CDATA[Camille Fournier]]></dc:creator>
            <pubDate>Wed, 05 Nov 2025 16:01:02 GMT</pubDate>
            <content:encoded><![CDATA[<p>This week, I gave a talk at <a href="https://leaddev.com/leaddev-new-york/">LDX3 New York</a>. Below is not exactly a transcript of what I said.</p><figure><img alt="" src="https://cdn-images-1.medium.com/max/1024/1*32LP8iirLzUIZrJJrQS7lA.png" /></figure><h3>Framing: A lot of people want management dead</h3><p>The talk started with 3 quotes from tech CEOs.</p><blockquote>“There is no need for a permanent middle management layer.” — <a href="https://block.xyz/inside/from-hierarchy-to-intelligence">Jack Dorsey</a></blockquote><h4>Technical Skills</h4><ol><li>First skill</li><li>Second skill</li></ol><img src="https://medium.com/_/stat?event=post.clientViewed&referrerSource=full_rss&postId=279cb6611d66" width="1" height="1" alt="">]]></content:encoded>
        </item>
        <item>
            <title><![CDATA[Members Only Post]]></title>
            <link>https://skamille.medium.com/members-only-post-0123456789ab?source=rss-12421a4f856------2</link>
            <content:encoded><![CDATA[]]></content:encoded>
        </item>
    </channel>
</rss>`

func newMediumTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/feed" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
		_, _ = w.Write([]byte(mediumFeedXML))
	}))
	t.Cleanup(server.Close)
	return server
}

func TestMediumFeedURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"author post", "https://medium.com/@skamille/some-post-279cb6611d66", "https://medium.com/feed/@skamille"},
		{"author profile", "https://medium.com/@skamille", "https://medium.com/feed/@skamille"},
		{"publication post", "https://medium.com/better-programming/a-post-0123456789ab?source=x", "https://medium.com/feed/better-programming"},
		{"subdomain post", "https://skamille.medium.com/the-managers-path-in-the-age-of-ai-279cb6611d66", "https://skamille.medium.com/feed"},
		{"subdomain profile", "https://skamille.medium.com/", "https://skamille.medium.com/feed"},
		{"custom domain", "https://blog.example.com/a-post-0123456789ab", "https://blog.example.com/feed"},
		{"feed URL", "https://medium.com/feed/@skamille", "https://medium.com/feed/@skamille"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mediumFeedURL(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestMediumFeedURL_Invalid(t *testing.T) {
	for _, input := range []string{"https://medium.com/", "not a url"} {
		if _, err := mediumFeedURL(input); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}

func TestMediumPostID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"plain", "https://skamille.medium.com/the-managers-path-in-the-age-of-ai-279cb6611d66", "279cb6611d66"},
		{"query", "https://skamille.medium.com/the-managers-path-in-the-age-of-ai-279cb6611d66?source=rss-12421a4f856------2", "279cb6611d66"},
		{"trailing slash", "https://medium.com/@skamille/a-post-e75438fdfa5d/", "e75438fdfa5d"},
		{"no dash", "https://medium.com/p/e75438fdfa5d", "e75438fdfa5d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := mediumPostID(tt.input); result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestMediumScraper_ScrapeArticle(t *testing.T) {
	server := newMediumTestServer(t)
	s := &MediumScraper{}

	result, err := s.ScrapeArticle(server.URL + "/the-managers-path-in-the-age-of-ai-279cb6611d66")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(result, "# The Manager’s Path in the Age of AI\n\n") {
		t.Errorf("expected title heading at start, got: %q", result)
	}
	for _, expected := range []string{
		"[LDX3 New York](https://leaddev.com/leaddev-new-york/)",
		"https://cdn-images-1.medium.com/max/1024/1*32LP8iirLzUIZrJJrQS7lA.png",
		"Framing: A lot of people want management dead",
		"There is no need for a permanent middle management layer.",
		"Technical Skills",
		"Second skill",
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("expected %q in result, got: %q", expected, result)
		}
	}
	if strings.Contains(result, "_/stat") {
		t.Errorf("expected tracking pixel to be excluded, got: %q", result)
	}
}

func TestMediumScraper_ScrapeArticle_NotInFeed(t *testing.T) {
	server := newMediumTestServer(t)
	s := &MediumScraper{}

	_, err := s.ScrapeArticle(server.URL + "/an-old-post-ffffffffffff")
	if err == nil || !strings.Contains(err.Error(), "not found in feed") {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestMediumScraper_ScrapeArticle_NoContent(t *testing.T) {
	server := newMediumTestServer(t)
	s := &MediumScraper{}

	_, err := s.ScrapeArticle(server.URL + "/members-only-post-0123456789ab")
	if err == nil || !strings.Contains(err.Error(), "paywall") {
		t.Errorf("expected paywall error, got: %v", err)
	}
}

func TestMediumScraper_ScrapeTitleAndFilename(t *testing.T) {
	server := newMediumTestServer(t)
	s := &MediumScraper{}
	postURL := server.URL + "/the-managers-path-in-the-age-of-ai-279cb6611d66"

	title, err := s.ScrapeTitle(postURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if title != "The Manager’s Path in the Age of AI" {
		t.Errorf("unexpected title: %q", title)
	}

	filename, err := s.ScrapeFilename(postURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filename != "the_manager’s_path_in_the_age_of_ai" {
		t.Errorf("unexpected filename: %q", filename)
	}
}

func TestMediumScraper_ScrapeLinks(t *testing.T) {
	server := newMediumTestServer(t)
	s := &MediumScraper{}

	links, err := s.ScrapeLinks(server.URL + "/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := map[string]string{
		"A sidebar on Code Review":            "https://skamille.medium.com/a-sidebar-on-code-review-e75438fdfa5d",
		"The Manager’s Path in the Age of AI": "https://skamille.medium.com/the-managers-path-in-the-age-of-ai-279cb6611d66",
		"Members Only Post":                   "https://skamille.medium.com/members-only-post-0123456789ab",
	}
	if len(links) != len(expected) {
		t.Errorf("expected %d links, got %d: %v", len(expected), len(links), links)
	}
	for text, link := range expected {
		if links[text] != link {
			t.Errorf("expected link %q for %q, got %q", link, text, links[text])
		}
	}
}
