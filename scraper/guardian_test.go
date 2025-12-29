package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGuardianScraper_ScrapeLinks_Basic(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<a data-link-name="group-0 | card-1" aria-label="Article One" href="/article/one">Link 1</a>
	<a data-link-name="group-0 | card-2" aria-label="Article Two" href="/article/two">Link 2</a>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeLinks(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 links, got %d", len(result))
	}

	if _, ok := result["Article One"]; !ok {
		t.Error("expected 'Article One' key to be present")
	}
	if _, ok := result["Article Two"]; !ok {
		t.Error("expected 'Article Two' key to be present")
	}
}

func TestGuardianScraper_ScrapeLinks_FiltersNonGroup0(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<a data-link-name="group-0 | card-1" aria-label="Featured" href="/article/featured">Featured</a>
	<a data-link-name="group-1 | card-1" aria-label="Other" href="/article/other">Other</a>
	<a data-link-name="nav-link" aria-label="Nav" href="/nav">Nav</a>
	<a aria-label="No Data Link" href="/no-data">No Data</a>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeLinks(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 link (only group-0), got %d: %v", len(result), result)
	}

	if _, ok := result["Featured"]; !ok {
		t.Error("expected 'Featured' key to be present")
	}
}

func TestGuardianScraper_ScrapeLinks_AbsoluteURL(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<a data-link-name="group-0 | card-1" aria-label="Article" href="/world/2024/article">Link</a>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeLinks(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	url := result["Article"]
	if !strings.HasPrefix(url, server.URL) {
		t.Errorf("expected absolute URL starting with %s, got %s", server.URL, url)
	}
	if !strings.HasSuffix(url, "/world/2024/article") {
		t.Errorf("expected URL to end with /world/2024/article, got %s", url)
	}
}

func TestGuardianScraper_ScrapeLinks_EmptyPage(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<p>No links here</p>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeLinks(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected 0 links, got %d", len(result))
	}
}

func TestGuardianScraper_ScrapeLinks_InvalidURL(t *testing.T) {
	scraper := &GuardianScraper{}
	_, err := scraper.ScrapeLinks("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestGuardianScraper_ScrapeArticle_Title(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Breaking News: Important Event</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "# Breaking News: Important Event\n\n"
	if result != expected {
		t.Errorf("ScrapeArticle() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeArticle_TitleAndSubtitle(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Main Headline</h1>
	<div data-gu-name="standfirst">
		<p>This is the article standfirst subtitle</p>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "# Main Headline") {
		t.Error("expected title to be present")
	}
	if !strings.Contains(result, "## This is the article standfirst subtitle") {
		t.Errorf("expected subtitle to be present, got: %q", result)
	}
}

func TestGuardianScraper_ScrapeArticle_FullArticle(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Article Title</h1>
	<div data-gu-name="standfirst">
		<p>Article subtitle goes here</p>
	</div>
	<div class="article-body-commercial-selector">
		<p>First paragraph of the article.</p>
		<p>Second paragraph of the article.</p>
		<p>Third paragraph of the article.</p>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectations := []string{
		"# Article Title",
		"## Article subtitle goes here",
		"First paragraph of the article.",
		"Second paragraph of the article.",
		"Third paragraph of the article.",
	}

	for _, exp := range expectations {
		if !strings.Contains(result, exp) {
			t.Errorf("expected %q in result, got: %q", exp, result)
		}
	}
}

func TestGuardianScraper_ScrapeArticle_BodyParagraphs(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<div class="article-body-commercial-selector">
		<p>Paragraph one.</p>
		<p>Paragraph two.</p>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check paragraphs are separated by double newlines
	if !strings.Contains(result, "Paragraph one.\n\n") {
		t.Errorf("expected paragraph one with double newline, got: %q", result)
	}
	if !strings.Contains(result, "Paragraph two.\n\n") {
		t.Errorf("expected paragraph two with double newline, got: %q", result)
	}
}

func TestGuardianScraper_ScrapeArticle_OnlyBodyContent(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<div class="other-div">
		<p>This should not appear</p>
	</div>
	<div class="article-body-commercial-selector">
		<p>This should appear</p>
	</div>
	<div class="sidebar">
		<p>Sidebar content</p>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "This should not appear") {
		t.Error("content from other-div should not be included")
	}
	if strings.Contains(result, "Sidebar content") {
		t.Error("content from sidebar should not be included")
	}
	if !strings.Contains(result, "This should appear") {
		t.Error("content from article-body-commercial-selector should be included")
	}
}

func TestGuardianScraper_ScrapeArticle_EmptyContent(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Empty Article</h1>
	<div data-gu-name="standfirst"></div>
	<div class="article-body-commercial-selector"></div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "# Empty Article\n\n"
	if result != expected {
		t.Errorf("ScrapeArticle() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeArticle_InvalidURL(t *testing.T) {
	scraper := &GuardianScraper{}
	_, err := scraper.ScrapeArticle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestGuardianScraper_ScrapeArticle_NoTitle(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<div data-gu-name="standfirst">
		<p>Subtitle only</p>
	</div>
	<div class="article-body-commercial-selector">
		<p>Body text</p>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should not have h1 marker when no h1 in HTML
	if strings.HasPrefix(result, "# ") {
		t.Error("should not have h1 marker when no h1 in HTML")
	}
	if !strings.Contains(result, "## Subtitle only") {
		t.Error("expected subtitle to be present")
	}
	if !strings.Contains(result, "Body text") {
		t.Error("expected body text to be present")
	}
}

func TestGuardianScraper_ScrapeArticle_MultipleH1(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>First Title</h1>
	<h1>Second Title</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Both h1 elements should be captured
	if !strings.Contains(result, "# First Title") {
		t.Error("expected first h1 to be present")
	}
	if !strings.Contains(result, "# Second Title") {
		t.Error("expected second h1 to be present")
	}
}

func TestGuardianScraper_ScrapeArticle_ContentOrder(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<div data-gu-name="standfirst">
		<p>Subtitle</p>
	</div>
	<div class="article-body-commercial-selector">
		<p>Body content</p>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check order: title, subtitle, body
	titleIdx := strings.Index(result, "# Title")
	subtitleIdx := strings.Index(result, "## Subtitle")
	bodyIdx := strings.Index(result, "Body content")

	if titleIdx == -1 || subtitleIdx == -1 || bodyIdx == -1 {
		t.Fatalf("expected all content to be present, got: %q", result)
	}

	if titleIdx >= subtitleIdx {
		t.Error("title should appear before subtitle")
	}
	if subtitleIdx >= bodyIdx {
		t.Error("subtitle should appear before body")
	}
}

func TestGuardianScraper_ScrapeTitle_Basic(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Guardian Article Title</h1>
	<div class="article-body-commercial-selector"></div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Guardian Article Title"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeTitle_WithWhitespace(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>   Title With Whitespace   </h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Title With Whitespace"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeTitle_WithNewlines(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>
		Title With Newlines
	</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Title With Newlines"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeTitle_NoH1(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h2>Only H2</h2>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// When no h1 exists, title should be empty
	if result != "" {
		t.Errorf("ScrapeTitle() = %q, want empty string", result)
	}
}

func TestGuardianScraper_ScrapeTitle_MultipleH1(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>First Title</h1>
	<h1>Second Title</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// When multiple h1 exist, the last one overwrites previous ones
	expected := "Second Title"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeTitle_EmptyH1(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1></h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("ScrapeTitle() = %q, want empty string", result)
	}
}

func TestGuardianScraper_ScrapeTitle_InvalidURL(t *testing.T) {
	scraper := &GuardianScraper{}
	_, err := scraper.ScrapeTitle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestGuardianScraper_ScrapeTitle_SpecialCharacters(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title with &amp; special &lt;characters&gt;</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// HTML entities should be decoded
	expected := "Title with & special <characters>"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeTitle_WithNestedElements(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1><span>Nested</span> <strong>Title</strong></h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Text from nested elements should be extracted
	expected := "Nested Title"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeTitle_UnicodeCharacters(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>World news: événements mondiaux</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "World news: événements mondiaux"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeFilename_Basic(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Breaking News Today</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "breaking_news_today"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeFilename_ConvertsToLowercase(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>BREAKING NEWS HEADLINE</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "breaking_news_headline"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeFilename_ReplacesSpacesWithUnderscores(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>multiple words in headline</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "multiple_words_in_headline"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeFilename_ReplacesDoubleUnderscores(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title  With  Double  Spaces</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Double spaces become double underscores, then single underscores
	expected := "title_with_double_spaces"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeFilename_EmptyTitle(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1></h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Empty title should produce empty filename
	if result != "" {
		t.Errorf("ScrapeFilename() = %q, want empty string", result)
	}
}

func TestGuardianScraper_ScrapeFilename_NoH1(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h2>Only H2</h2>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No h1 means empty title, which produces empty filename
	if result != "" {
		t.Errorf("ScrapeFilename() = %q, want empty string", result)
	}
}

func TestGuardianScraper_ScrapeFilename_InvalidURL(t *testing.T) {
	scraper := &GuardianScraper{}
	_, err := scraper.ScrapeFilename("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestGuardianScraper_ScrapeFilename_SingleWord(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Politics</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "politics"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeFilename_PreservesSpecialCharacters(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>News from 2024: Latest Updates</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Numbers and colons are preserved
	expected := "news_from_2024:_latest_updates"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestGuardianScraper_ScrapeFilename_UnicodeCharacters(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>World news: événements mondiaux</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GuardianScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Unicode characters are preserved (just lowercased)
	expected := "world_news:_événements_mondiaux"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}
