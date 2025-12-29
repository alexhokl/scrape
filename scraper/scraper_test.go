package scraper

import (
	"errors"
	"testing"
)

// mockLinkScraper is a test implementation of LinkScraper
type mockLinkScraper struct {
	links map[string]string
	err   error
}

func (m *mockLinkScraper) ScrapeLinks(url string) (map[string]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.links, nil
}

// mockArticleScraper is a test implementation of ArticleScraper
type mockArticleScraper struct {
	content string
	title   string
	err     error
}

func (m *mockArticleScraper) ScrapeArticle(url string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.content, nil
}

func (m *mockArticleScraper) ScrapeTitle(url string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.title, nil
}

func (m *mockArticleScraper) ScrapeFilename(url string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.title, nil
}

// Compile-time interface implementation checks
var _ LinkScraper = (*mockLinkScraper)(nil)
var _ ArticleScraper = (*mockArticleScraper)(nil)

func TestLinkScraper_ReturnsLinks(t *testing.T) {
	expectedLinks := map[string]string{
		"Example Article": "https://example.com/article1",
		"Another Article": "https://example.com/article2",
		"Third Article":   "https://example.com/article3",
	}

	scraper := &mockLinkScraper{links: expectedLinks}

	links, err := scraper.ScrapeLinks("https://example.com")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(links) != len(expectedLinks) {
		t.Errorf("expected %d links, got %d", len(expectedLinks), len(links))
	}

	for text, url := range expectedLinks {
		if links[text] != url {
			t.Errorf("expected link %q -> %q, got %q", text, url, links[text])
		}
	}
}

func TestLinkScraper_ReturnsEmptyMap(t *testing.T) {
	scraper := &mockLinkScraper{links: map[string]string{}}

	links, err := scraper.ScrapeLinks("https://example.com/empty")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if links == nil {
		t.Error("expected non-nil map, got nil")
	}

	if len(links) != 0 {
		t.Errorf("expected empty map, got %d items", len(links))
	}
}

func TestLinkScraper_ReturnsError(t *testing.T) {
	expectedErr := errors.New("failed to scrape links: network error")
	scraper := &mockLinkScraper{err: expectedErr}

	links, err := scraper.ScrapeLinks("https://example.com/error")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if links != nil {
		t.Errorf("expected nil links on error, got %v", links)
	}
}

func TestArticleScraper_ReturnsContent(t *testing.T) {
	expectedContent := "# Article Title\n\nThis is the article content in markdown format."
	expectedTitle := "Article Title"
	scraper := &mockArticleScraper{content: expectedContent, title: expectedTitle}

	content, err := scraper.ScrapeArticle("https://example.com/article")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if content != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, content)
	}
}

func TestArticleScraper_ReturnsEmptyContent(t *testing.T) {
	scraper := &mockArticleScraper{content: ""}

	content, err := scraper.ScrapeArticle("https://example.com/empty")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if content != "" {
		t.Errorf("expected empty content, got %q", content)
	}
}

func TestArticleScraper_ReturnsError(t *testing.T) {
	expectedErr := errors.New("failed to scrape article: page not found")
	scraper := &mockArticleScraper{err: expectedErr}

	content, err := scraper.ScrapeArticle("https://example.com/notfound")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if content != "" {
		t.Errorf("expected empty content on error, got %q", content)
	}
}

func TestLinkScraper_InterfaceAcceptsAnyURL(t *testing.T) {
	testCases := []string{
		"https://example.com",
		"http://localhost:8080/path",
		"https://example.com/path?query=value",
		"",
	}

	scraper := &mockLinkScraper{links: map[string]string{}}

	for _, url := range testCases {
		t.Run(url, func(t *testing.T) {
			_, err := scraper.ScrapeLinks(url)
			if err != nil {
				t.Errorf("unexpected error for URL %q: %v", url, err)
			}
		})
	}
}

func TestArticleScraper_InterfaceAcceptsAnyURL(t *testing.T) {
	testCases := []string{
		"https://example.com/article",
		"http://localhost:8080/post/123",
		"https://example.com/path?id=1",
		"",
	}

	scraper := &mockArticleScraper{content: "test content"}

	for _, url := range testCases {
		t.Run(url, func(t *testing.T) {
			_, err := scraper.ScrapeArticle(url)
			if err != nil {
				t.Errorf("unexpected error for URL %q: %v", url, err)
			}
		})
	}
}
