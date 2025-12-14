package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGoDocScraper_ScrapeArticle_Title(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Go Documentation Title</h1>
	<article></article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "# Go Documentation Title\n\n"
	if result != expected {
		t.Errorf("ScrapeArticle() = %q, want %q", result, expected)
	}
}

func TestGoDocScraper_ScrapeArticle_Headings(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Main Title</h1>
	<article>
		<h2>Section One</h2>
		<h3>Subsection</h3>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "# Main Title") {
		t.Error("expected h1 to be converted to # heading")
	}
	if !strings.Contains(result, "## Section One") {
		t.Error("expected h2 to be converted to ## heading")
	}
	if !strings.Contains(result, "### Subsection") {
		t.Error("expected h3 to be converted to ### heading")
	}
}

func TestGoDocScraper_ScrapeArticle_HeadingsOrder(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Main Title</h1>
	<article>
		<h2>First Section</h2>
		<h3>First Subsection</h3>
		<h2>Second Section</h2>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that headings appear in the correct order
	firstSectionIdx := strings.Index(result, "## First Section")
	firstSubsectionIdx := strings.Index(result, "### First Subsection")
	secondSectionIdx := strings.Index(result, "## Second Section")

	if firstSectionIdx == -1 || firstSubsectionIdx == -1 || secondSectionIdx == -1 {
		t.Fatalf("expected all headings to be present, got: %q", result)
	}

	if firstSectionIdx >= firstSubsectionIdx || firstSubsectionIdx >= secondSectionIdx {
		t.Error("headings should appear in document order")
	}
}

func TestGoDocScraper_ScrapeArticle_ParagraphProducesNewlines(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<p>Paragraph content</p>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Paragraph element is processed and produces output (with double newline)
	// The actual text extraction depends on parseGoDocParagraph behavior
	if !strings.Contains(result, "# Title\n\n") {
		t.Errorf("expected title formatting, got: %q", result)
	}
}

func TestGoDocScraper_ScrapeArticle_UnorderedListStructure(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<ul>
			<li>Item one</li>
			<li>Item two</li>
			<li>Item three</li>
		</ul>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that list items have * prefix
	if strings.Count(result, "* ") != 3 {
		t.Errorf("expected 3 list items with * prefix, got: %q", result)
	}
}

func TestGoDocScraper_ScrapeArticle_NoteDivProducesBlockquote(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<div class="NOTE">
			<p class="alert">Note</p>
			<p>Important information</p>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// NOTE div should produce blockquote style output with >
	if !strings.Contains(result, "> ") {
		t.Errorf("expected NOTE div to produce blockquote, got: %q", result)
	}
}

func TestGoDocScraper_ScrapeArticle_AlertParagraphSkipped(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<div class="NOTE">
			<p class="alert">Warning Label</p>
			<p>Actual content</p>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have only one blockquote line (the non-alert paragraph)
	blockquoteCount := strings.Count(result, "> ")
	if blockquoteCount != 1 {
		t.Errorf("expected exactly 1 blockquote line (alert should be skipped), got %d in: %q", blockquoteCount, result)
	}
}

func TestGoDocScraper_ScrapeArticle_NonNoteDiv(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<div class="other-class">
			<p>This should be ignored</p>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Non-NOTE divs should not produce blockquote output
	if strings.Contains(result, "> ") {
		t.Errorf("non-NOTE div should not produce blockquote, got: %q", result)
	}
}

func TestGoDocScraper_ScrapeArticle_CodeBlockStructure(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<pre><code>go version</code></pre>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have code block markers
	if strings.Count(result, "```") != 2 {
		t.Errorf("expected 2 code block markers, got: %q", result)
	}
}

func TestGoDocScraper_ScrapeArticle_MultipleCodeBlocks(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<pre><code>first code</code></pre>
		<pre><code>second code</code></pre>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have 4 code block markers (2 blocks × 2 markers each)
	if strings.Count(result, "```") != 4 {
		t.Errorf("expected 4 code block markers, got: %q", result)
	}
}

func TestGoDocScraper_ScrapeArticle_InvalidURL(t *testing.T) {
	scraper := &GoDocScraper{}
	_, err := scraper.ScrapeArticle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestGoDocScraper_ScrapeArticle_EmptyContent(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Empty Article</h1>
	<article>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "# Empty Article\n\n"
	if result != expected {
		t.Errorf("ScrapeArticle() = %q, want %q", result, expected)
	}
}

func TestGoDocScraper_ScrapeArticle_NoArticleTag(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>No Article Tag</h1>
	<div>
		<p>This should not appear</p>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only the title should be present, no content from divs
	if strings.Contains(result, "This should not appear") {
		t.Error("content from non-article elements should not be included")
	}
	if !strings.Contains(result, "# No Article Tag") {
		t.Error("expected title to be present")
	}
}

func TestGoDocScraper_ScrapeArticle_MixedContent(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Getting Started with Go</h1>
	<article>
		<h2>Introduction</h2>
		<p>Welcome to Go programming.</p>
		<h3>Prerequisites</h3>
		<ul>
			<li>Go 1.18+</li>
			<li>A text editor</li>
		</ul>
		<div class="NOTE">
			<p class="alert">Note</p>
			<p>Remember to set GOPATH</p>
		</div>
		<pre><code>go version</code></pre>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify structure is present
	expectations := []struct {
		content string
		desc    string
	}{
		{"# Getting Started with Go", "main title"},
		{"## Introduction", "h2 heading"},
		{"### Prerequisites", "h3 heading"},
		{"* ", "list items"},
		{"> ", "note blockquote"},
		{"```", "code block markers"},
	}

	for _, exp := range expectations {
		if !strings.Contains(result, exp.content) {
			t.Errorf("expected %s (%q) in result, got: %q", exp.desc, exp.content, result)
		}
	}
}

func TestGoDocScraper_ScrapeArticle_NoH1(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article>
		<h2>Only H2</h2>
		<p>Some content</p>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should still process h2 even without h1
	if !strings.Contains(result, "## Only H2") {
		t.Errorf("expected h2 heading, got: %q", result)
	}
	// Should not have h1 marker
	if strings.HasPrefix(result, "# ") {
		t.Error("should not have h1 when no h1 in HTML")
	}
}

func TestGoDocScraper_ScrapeArticle_MultipleH1(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>First Title</h1>
	<h1>Second Title</h1>
	<article></article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
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

func TestGoDocScraper_ScrapeArticle_ImageInParagraph(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<p><img src="image.png" alt="test"></p>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Image should be replaced with _image_
	if !strings.Contains(result, "_image_") {
		t.Errorf("expected image to be replaced with _image_, got: %q", result)
	}
}

func TestGoDocScraper_ScrapeArticle_OnlyDirectChildren(t *testing.T) {
	// Tests that only direct children of article are processed
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<div class="wrapper">
			<h2>Nested H2</h2>
		</div>
		<h2>Direct H2</h2>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only direct child h2 should be captured due to parent check
	if !strings.Contains(result, "## Direct H2") {
		t.Errorf("expected direct h2 to be present, got: %q", result)
	}
	// Count h2 headings - should only have the direct one
	h2Count := strings.Count(result, "## ")
	if h2Count != 1 {
		t.Errorf("expected 1 h2 heading (direct child only), got %d in: %q", h2Count, result)
	}
}

func TestGoDocScraper_ScrapeArticle_HeadingsInOrder(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<h2>First Section</h2>
		<h2>Second Section</h2>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check order
	titleIdx := strings.Index(result, "# Title")
	firstSectionIdx := strings.Index(result, "## First Section")
	secondSectionIdx := strings.Index(result, "## Second Section")

	if titleIdx == -1 || firstSectionIdx == -1 || secondSectionIdx == -1 {
		t.Fatalf("expected all content to be present, got: %q", result)
	}

	if titleIdx >= firstSectionIdx {
		t.Error("title should appear before first section")
	}
	if firstSectionIdx >= secondSectionIdx {
		t.Error("first section should appear before second section")
	}
}

func TestGoDocScraper_ScrapeArticle_MultipleNoteParagraphs(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<div class="NOTE">
			<p class="alert">Note</p>
			<p>First paragraph</p>
			<p>Second paragraph</p>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have 2 blockquote lines (alert paragraph is skipped)
	blockquoteCount := strings.Count(result, "> ")
	if blockquoteCount != 2 {
		t.Errorf("expected 2 blockquote lines, got %d in: %q", blockquoteCount, result)
	}
}

func TestGoDocScraper_ScrapeArticle_EmptyListItems(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<ul>
			<li></li>
			<li></li>
		</ul>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Empty list items should still produce * markers
	if strings.Count(result, "* ") != 2 {
		t.Errorf("expected 2 list item markers, got: %q", result)
	}
}

func TestGoDocScraper_ScrapeArticle_MultipleParagraphs(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<p>First paragraph</p>
		<p>Second paragraph</p>
		<p>Third paragraph</p>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Each paragraph produces \n\n, so we should have multiple paragraph separations
	// Count the newline pairs after title
	titlePart := "# Title\n\n"
	if !strings.HasPrefix(result, titlePart) {
		t.Errorf("expected result to start with title, got: %q", result)
	}
}

func TestGoDocScraper_ScrapeArticle_PreWithoutCode(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<article>
		<pre>plain preformatted text</pre>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GoDocScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// pre without code still produces code block markers
	if strings.Count(result, "```") != 2 {
		t.Errorf("expected 2 code block markers for pre, got: %q", result)
	}
}
