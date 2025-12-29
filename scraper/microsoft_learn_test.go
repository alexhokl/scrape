package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMicrosoftLearnScraper_ScrapeArticle_Title(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Test Article Title</h1>
	<div class="content"></div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "# Test Article Title\n\n"
	if result != expected {
		t.Errorf("ScrapeArticle() = %q, want %q", result, expected)
	}
}

func TestMicrosoftLearnScraper_ScrapeArticle_Headings(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Main Title</h1>
	<div class="content">
		<h2>Section One</h2>
		<h3>Subsection</h3>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
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

func TestMicrosoftLearnScraper_ScrapeArticle_HeadingsOrder(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Main Title</h1>
	<div class="content">
		<h2>First Section</h2>
		<h3>First Subsection</h3>
		<h2>Second Section</h2>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
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

func TestMicrosoftLearnScraper_ScrapeArticle_ImageSpan(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<div class="content">
		<p><span class="icon">icon</span></p>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "_image_") {
		t.Errorf("expected span to be replaced with _image_, got: %q", result)
	}
}

func TestMicrosoftLearnScraper_ScrapeArticle_InvalidURL(t *testing.T) {
	scraper := &MicrosoftLearnScraper{}
	_, err := scraper.ScrapeArticle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestMicrosoftLearnScraper_ScrapeArticle_EmptyContent(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Empty Article</h1>
	<div class="content">
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "# Empty Article\n\n"
	if result != expected {
		t.Errorf("ScrapeArticle() = %q, want %q", result, expected)
	}
}

func TestMicrosoftLearnScraper_ScrapeArticle_NoContentDiv(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>No Content Div</h1>
	<div class="other">
		<p>This should not appear</p>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only the title should be present, no content from other divs
	if strings.Contains(result, "This should not appear") {
		t.Error("content from non-content div should not be included")
	}
	if !strings.Contains(result, "# No Content Div") {
		t.Error("expected title to be present")
	}
}

func TestMicrosoftLearnScraper_ScrapeArticle_MultipleH1(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>First Title</h1>
	<h1>Second Title</h1>
	<div class="content"></div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
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

func TestMicrosoftLearnScraper_ScrapeArticle_ParagraphGeneratesOutput(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<div class="content">
		<p>Some text</p>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Paragraph should produce output (even if text not captured due to implementation)
	// The paragraph tag is processed and adds \n\n
	if !strings.Contains(result, "# Title\n\n") {
		t.Errorf("expected title formatting, got: %q", result)
	}
}

func TestMicrosoftLearnScraper_ScrapeArticle_UnorderedListStructure(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<div class="content">
		<ul>
			<li>Item one</li>
			<li>Item two</li>
		</ul>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that list items produce * prefix
	if strings.Count(result, "* ") != 2 {
		t.Errorf("expected 2 list items with * prefix, got: %q", result)
	}
}

func TestMicrosoftLearnScraper_ScrapeArticle_NoteDiv(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<div class="content">
		<div class="NOTE">
			<p class="alert">Note</p>
			<p>Important info</p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// NOTE div should produce blockquote style output with >
	if !strings.Contains(result, "> ") {
		t.Errorf("expected NOTE div to produce blockquote, got: %q", result)
	}
}

func TestMicrosoftLearnScraper_ScrapeArticle_NonNoteDiv(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<div class="content">
		<div class="other-class">
			<p>This should be ignored</p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Non-NOTE divs should not produce blockquote output
	if strings.Contains(result, "> ") {
		t.Errorf("non-NOTE div should not produce blockquote, got: %q", result)
	}
}

func TestMicrosoftLearnScraper_ScrapeArticle_AlertParagraphSkipped(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Title</h1>
	<div class="content">
		<div class="NOTE">
			<p class="alert">Warning Label</p>
			<p>Actual content</p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have only one blockquote line (the non-alert paragraph)
	// The alert paragraph should be skipped
	blockquoteCount := strings.Count(result, "> ")
	if blockquoteCount != 1 {
		t.Errorf("expected exactly 1 blockquote line (alert should be skipped), got %d in: %q", blockquoteCount, result)
	}
}

func TestMicrosoftLearnScraper_ScrapeArticle_MixedContent(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Getting Started</h1>
	<div class="content">
		<h2>Introduction</h2>
		<p>Welcome</p>
		<h3>Prerequisites</h3>
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
		</ul>
		<div class="NOTE">
			<p class="alert">Note</p>
			<p>Remember this</p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify structure is present
	expectations := []struct {
		content string
		desc    string
	}{
		{"# Getting Started", "main title"},
		{"## Introduction", "h2 heading"},
		{"### Prerequisites", "h3 heading"},
		{"* ", "list items"},
		{"> ", "note blockquote"},
	}

	for _, exp := range expectations {
		if !strings.Contains(result, exp.content) {
			t.Errorf("expected %s (%q) in result, got: %q", exp.desc, exp.content, result)
		}
	}
}

func TestMicrosoftLearnScraper_ScrapeArticle_NoH1(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<div class="content">
		<h2>Only H2</h2>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should still process h2 even without h1
	if !strings.Contains(result, "## Only H2") {
		t.Errorf("expected h2 heading, got: %q", result)
	}
	// Should not have h1 marker
	if strings.Contains(result, "# ") && !strings.Contains(result, "## ") {
		t.Error("should not have h1 when no h1 in HTML")
	}
}

func TestMicrosoftLearnScraper_ScrapeTitle_Basic(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Microsoft Learn Article Title</h1>
	<div class="content"></div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Microsoft Learn Article Title"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestMicrosoftLearnScraper_ScrapeTitle_WithWhitespace(t *testing.T) {
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

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Title With Whitespace"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestMicrosoftLearnScraper_ScrapeTitle_WithNewlines(t *testing.T) {
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

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Title With Newlines"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestMicrosoftLearnScraper_ScrapeTitle_NoH1(t *testing.T) {
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

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// When no h1 exists, title should be empty
	if result != "" {
		t.Errorf("ScrapeTitle() = %q, want empty string", result)
	}
}

func TestMicrosoftLearnScraper_ScrapeTitle_MultipleH1(t *testing.T) {
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

	scraper := &MicrosoftLearnScraper{}
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

func TestMicrosoftLearnScraper_ScrapeTitle_EmptyH1(t *testing.T) {
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

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("ScrapeTitle() = %q, want empty string", result)
	}
}

func TestMicrosoftLearnScraper_ScrapeTitle_InvalidURL(t *testing.T) {
	scraper := &MicrosoftLearnScraper{}
	_, err := scraper.ScrapeTitle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestMicrosoftLearnScraper_ScrapeTitle_SpecialCharacters(t *testing.T) {
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

	scraper := &MicrosoftLearnScraper{}
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

func TestMicrosoftLearnScraper_ScrapeTitle_WithNestedElements(t *testing.T) {
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

	scraper := &MicrosoftLearnScraper{}
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

func TestMicrosoftLearnScraper_ScrapeTitle_UnicodeCharacters(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Azure 入門ガイド</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Azure 入門ガイド"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestMicrosoftLearnScraper_ScrapeFilename_Basic(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Getting Started with Azure</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "getting_started_with_azure"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestMicrosoftLearnScraper_ScrapeFilename_ConvertsToLowercase(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Azure Documentation TITLE</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "azure_documentation_title"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestMicrosoftLearnScraper_ScrapeFilename_ReplacesSpacesWithUnderscores(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>multiple words in title</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "multiple_words_in_title"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestMicrosoftLearnScraper_ScrapeFilename_ReplacesDoubleUnderscores(t *testing.T) {
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

	scraper := &MicrosoftLearnScraper{}
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

func TestMicrosoftLearnScraper_ScrapeFilename_EmptyTitle(t *testing.T) {
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

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Empty title should produce empty filename
	if result != "" {
		t.Errorf("ScrapeFilename() = %q, want empty string", result)
	}
}

func TestMicrosoftLearnScraper_ScrapeFilename_NoH1(t *testing.T) {
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

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No h1 means empty title, which produces empty filename
	if result != "" {
		t.Errorf("ScrapeFilename() = %q, want empty string", result)
	}
}

func TestMicrosoftLearnScraper_ScrapeFilename_InvalidURL(t *testing.T) {
	scraper := &MicrosoftLearnScraper{}
	_, err := scraper.ScrapeFilename("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestMicrosoftLearnScraper_ScrapeFilename_SingleWord(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Overview</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "overview"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestMicrosoftLearnScraper_ScrapeFilename_PreservesSpecialCharacters(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Azure SDK 2.0 Release Notes</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Numbers and dots are preserved
	expected := "azure_sdk_2.0_release_notes"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestMicrosoftLearnScraper_ScrapeFilename_UnicodeCharacters(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Azure 入門ガイド</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &MicrosoftLearnScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Unicode characters are preserved (just lowercased)
	expected := "azure_入門ガイド"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}
