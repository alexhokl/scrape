package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// siteScraper is a scraper implementing both interfaces, as the scrapers
// built on siteConfig do.
type siteScraper interface {
	ArticleScraper
	LinkScraper
}

// siteTestCase describes the fixtures and expectations of a siteConfig-based
// scraper.
type siteTestCase struct {
	scraper siteScraper
	// articleHTML is served at /blog/sample-post/
	articleHTML string
	title       string
	// contains lists markdown expected in the article
	contains []string
	// excludes lists text expected to be left out of the article
	excludes []string
	// indexHTML is served at /
	indexHTML string
	// links maps the expected link text to its path on the test server
	links map[string]string
	// filename is the expected file name; empty means "sample-post"
	filename string
}

const siteTestArticlePath = "/blog/sample-post/"

func newSiteTestServer(articleHTML, indexHTML string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		switch r.URL.Path {
		case siteTestArticlePath:
			_, _ = w.Write([]byte(articleHTML))
		case "/":
			_, _ = w.Write([]byte(indexHTML))
		default:
			http.NotFound(w, r)
		}
	}))
}

func runSiteTests(t *testing.T, tc siteTestCase) {
	t.Helper()
	server := newSiteTestServer(tc.articleHTML, tc.indexHTML)
	defer server.Close()
	articleURL := server.URL + siteTestArticlePath

	t.Run("ScrapeArticle", func(t *testing.T) {
		result, err := tc.scraper.ScrapeArticle(articleURL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.HasPrefix(result, "# "+tc.title+"\n\n") {
			t.Errorf("expected title heading %q at start, got: %q", tc.title, result)
		}
		for _, expected := range tc.contains {
			if !strings.Contains(result, expected) {
				t.Errorf("expected %q in result, got: %q", expected, result)
			}
		}
		for _, unexpected := range tc.excludes {
			if strings.Contains(result, unexpected) {
				t.Errorf("expected %q to be excluded, got: %q", unexpected, result)
			}
		}
	})

	t.Run("ScrapeTitle", func(t *testing.T) {
		result, err := tc.scraper.ScrapeTitle(articleURL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != tc.title {
			t.Errorf("expected %q, got %q", tc.title, result)
		}
	})

	t.Run("ScrapeFilename", func(t *testing.T) {
		result, err := tc.scraper.ScrapeFilename(articleURL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := tc.filename
		if expected == "" {
			expected = "sample-post"
		}
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("ScrapeLinks", func(t *testing.T) {
		result, err := tc.scraper.ScrapeLinks(server.URL + "/")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != len(tc.links) {
			t.Errorf("expected %d links, got %d: %v", len(tc.links), len(result), result)
		}
		for text, path := range tc.links {
			if result[text] != server.URL+path {
				t.Errorf("expected link %q to be %q, got %q", text, server.URL+path, result[text])
			}
		}
	})

	t.Run("ScrapeArticle_NotFound", func(t *testing.T) {
		_, err := tc.scraper.ScrapeArticle(server.URL + "/missing/")
		if err == nil {
			t.Error("expected an error for a missing page")
		}
	})
}

func TestSiteConfig_Article_MissingBody(t *testing.T) {
	server := newSiteTestServer(`<html><body><h1>Title</h1><p>No body here.</p></body></html>`, ``)
	defer server.Close()

	site := siteConfig{title: "h1", body: "div.content"}
	_, err := site.article(server.URL + siteTestArticlePath)
	if err == nil {
		t.Fatal("expected an error when the body is not found")
	}
	if !strings.Contains(err.Error(), "article content not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSiteConfig_Title_FallsBackToOpenGraph(t *testing.T) {
	server := newSiteTestServer(`<html><head><meta property="og:title" content="  OG   Title "></head><body><div class="content"><p>Text.</p></div></body></html>`, ``)
	defer server.Close()

	site := siteConfig{title: "h1.missing", body: "div.content"}
	title, err := site.scrapeTitle(server.URL + siteTestArticlePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if title != "OG Title" {
		t.Errorf("expected %q, got %q", "OG Title", title)
	}
}

func TestSiteConfig_Title_NotFound(t *testing.T) {
	server := newSiteTestServer(`<html><body><p>Nothing.</p></body></html>`, ``)
	defer server.Close()

	site := siteConfig{title: "h1", body: "body"}
	_, err := site.scrapeTitle(server.URL + siteTestArticlePath)
	if err == nil {
		t.Fatal("expected an error when the title is not found")
	}
}

func TestSiteConfig_Links_DeduplicatesAndUsesLinkText(t *testing.T) {
	index := `<html><body>
<a class="card" href="/one/#comments"><h3>One</h3><p>Excerpt</p></a>
<a class="card" href="/one/"><h3>One again</h3></a>
<a class="card" href="https://other.example/two"><h3>Two</h3></a>
<a class="card" href="/empty/"><h3> </h3></a>
<a class="card"><h3>No href</h3></a>
</body></html>`
	server := newSiteTestServer(``, index)
	defer server.Close()

	site := siteConfig{links: "a.card", linkText: "h3"}
	links, err := site.scrapeLinks(server.URL + "/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := map[string]string{
		"One": server.URL + "/one/",
		"Two": "https://other.example/two",
	}
	if len(links) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, links)
	}
	for text, url := range expected {
		if links[text] != url {
			t.Errorf("expected %q to be %q, got %q", text, url, links[text])
		}
	}
}

func TestSiteConfig_Links_ElementWrappingAnchor(t *testing.T) {
	server := newSiteTestServer(``, `<html><body><h2 class="t"><a href="/post/">Post</a></h2></body></html>`)
	defer server.Close()

	site := siteConfig{links: "h2.t"}
	links, err := site.scrapeLinks(server.URL + "/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if links["Post"] != server.URL+"/post/" {
		t.Errorf("unexpected links: %v", links)
	}
}

func TestSiteConfig_Filename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"trailing slash", "https://example.com/blog/my-post/", "my-post"},
		{"no trailing slash", "https://example.com/posts/2026/Oct/02/My-Post", "My-Post"},
		{"html extension", "https://example.com/2026/10/07/my-post.html", "my-post"},
		{"query string", "https://example.com/blog/my-post?source=rss", "my-post"},
		{"root", "https://example.com/", "article"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := siteConfig{}.filename(tt.input)
			if result != tt.expected {
				t.Errorf("filename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSiteConfig_UserAgent(t *testing.T) {
	var got string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.UserAgent()
		_, _ = w.Write([]byte(`<html><body><h1>Title</h1></body></html>`))
	}))
	defer server.Close()

	site := siteConfig{title: "h1", userAgent: "Mozilla/5.0 (compatible; scrape)"}
	if _, err := site.scrapeTitle(server.URL); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != site.userAgent {
		t.Errorf("expected User-Agent %q, got %q", site.userAgent, got)
	}
}
