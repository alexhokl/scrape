package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFirstLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single line",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "multiple lines",
			input:    "first line\nsecond line\nthird line",
			expected: "first line",
		},
		{
			name:     "only newline",
			input:    "\n",
			expected: "",
		},
		{
			name:     "empty first line",
			input:    "\nsecond line",
			expected: "",
		},
		{
			name:     "trailing newline",
			input:    "hello world\n",
			expected: "hello world",
		},
		{
			name:     "with carriage return and newline",
			input:    "first line\r\nsecond line",
			expected: "first line\r",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := firstLine(tt.input)
			if result != tt.expected {
				t.Errorf("firstLine(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseTofuguTableCell(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple text",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "text with leading and trailing spaces",
			input:    "  hello world  ",
			expected: "hello world",
		},
		{
			name:     "text with single newline",
			input:    "hello\nworld",
			expected: "hello; world",
		},
		{
			name:     "text with double newline",
			input:    "hello\n\nworld",
			expected: "hello; world",
		},
		{
			name:     "text with multiple newlines",
			input:    "hello\n\n\nworld",
			expected: "hello; ; world",
		},
		{
			name:     "text with leading newlines",
			input:    "\n\nhello world",
			expected: "hello world",
		},
		{
			name:     "text with trailing newlines",
			input:    "hello world\n\n",
			expected: "hello world",
		},
		{
			name:     "complex case",
			input:    "  first\n\nsecond\nthird  ",
			expected: "first; second; third",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTofuguTableCell(tt.input)
			if result != tt.expected {
				t.Errorf("parseTofuguTableCell(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseTofuguBlockquote(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single line",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "text with leading and trailing spaces",
			input:    "  hello world  ",
			expected: "hello world",
		},
		{
			name:     "multiline text",
			input:    "first line\nsecond line",
			expected: "first line\n> second line",
		},
		{
			name:     "multiline text with leading spaces",
			input:    "  first line\nsecond line  ",
			expected: "first line\n> second line",
		},
		{
			name:     "three lines",
			input:    "line one\nline two\nline three",
			expected: "line one\n> line two\n> line three",
		},
		{
			name:     "with leading newline",
			input:    "\nfirst line\nsecond line",
			expected: "first line\n> second line",
		},
		{
			name:     "with trailing newline",
			input:    "first line\nsecond line\n",
			expected: "first line\n> second line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTofuguBlockquote(tt.input)
			if result != tt.expected {
				t.Errorf("parseTofuguBlockquote(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Tests for TofuguScraper.ScrapeArticle

func TestTofuguScraper_ScrapeArticle_Title(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Japanese Grammar Guide</h1>
	<article><div class="main"></div></article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "# Japanese Grammar Guide\n\n"
	if result != expected {
		t.Errorf("ScrapeArticle() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeArticle_TitleAndMeta(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Learn Hiragana</h1>
	<div class="article-header-elements">
		<ul class="meta">
			<li>By Tofugu</li>
			<li>March 2024</li>
		</ul>
	</div>
	<article><div class="main"></div></article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "# Learn Hiragana") {
		t.Error("expected title to be present")
	}
	// Meta should be included
	if !strings.Contains(result, "Tofugu") || !strings.Contains(result, "March 2024") {
		t.Errorf("expected meta information to be present, got: %q", result)
	}
}

func TestTofuguScraper_ScrapeArticle_Headings(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Main Title</h1>
	<article>
		<div class="main">
			<h2>Section One</h2>
			<h3>Subsection A</h3>
			<h4>Detail Level</h4>
			<h5>Deep Level</h5>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
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
	if !strings.Contains(result, "### Subsection A") {
		t.Error("expected h3 to be converted to ### heading")
	}
	if !strings.Contains(result, "#### Detail Level") {
		t.Error("expected h4 to be converted to #### heading")
	}
	if !strings.Contains(result, "##### Deep Level") {
		t.Error("expected h5 to be converted to ##### heading")
	}
}

func TestTofuguScraper_ScrapeArticle_Paragraphs(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title</h1>
	<article>
		<div class="main">
			<p>First paragraph content.</p>
			<p>Second paragraph content.</p>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "First paragraph content.\n\n") {
		t.Errorf("expected first paragraph with double newline, got: %q", result)
	}
	if !strings.Contains(result, "Second paragraph content.\n\n") {
		t.Errorf("expected second paragraph with double newline, got: %q", result)
	}
}

func TestTofuguScraper_ScrapeArticle_OrderedList(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title</h1>
	<article>
		<div class="main">
			<ol>
				<li>First item</li>
				<li>Second item</li>
				<li>Third item</li>
			</ol>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "1. First item") {
		t.Error("expected numbered list item 1")
	}
	if !strings.Contains(result, "2. Second item") {
		t.Error("expected numbered list item 2")
	}
	if !strings.Contains(result, "3. Third item") {
		t.Error("expected numbered list item 3")
	}
}

func TestTofuguScraper_ScrapeArticle_ExampleSentence(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title</h1>
	<article>
		<div class="main">
			<ul class="example-sentence">
				<li>日本語の文</li>
				<li>English translation</li>
			</ul>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Example") {
		t.Error("expected 'Example' label")
	}
	if !strings.Contains(result, "- Japanese:") {
		t.Error("expected Japanese label")
	}
	if !strings.Contains(result, "- English:") {
		t.Error("expected English label")
	}
	if !strings.Contains(result, "日本語の文") {
		t.Error("expected Japanese text")
	}
	if !strings.Contains(result, "English translation") {
		t.Error("expected English text")
	}
}

func TestTofuguScraper_ScrapeArticle_TableOfContents(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title</h1>
	<article>
		<div class="main">
			<ul>
				<li>Introduction</li>
				<li>Chapter One</li>
				<li>Chapter Two</li>
			</ul>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that list items have * prefix
	if !strings.Contains(result, "* Introduction") {
		t.Error("expected bullet list item for Introduction")
	}
	if !strings.Contains(result, "* Chapter One") {
		t.Error("expected bullet list item for Chapter One")
	}
	if !strings.Contains(result, "* Chapter Two") {
		t.Error("expected bullet list item for Chapter Two")
	}
}

func TestTofuguScraper_ScrapeArticle_NestedTableOfContents(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title</h1>
	<article>
		<div class="main">
			<ul>
				<li>Chapter One
					<ul>
						<li>Section 1.1
							<ul>
								<li>Subsection 1.1.1</li>
							</ul>
						</li>
					</ul>
				</li>
			</ul>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "* Chapter One") {
		t.Error("expected top-level list item")
	}
	if !strings.Contains(result, "  * Section 1.1") {
		t.Error("expected indented nested list item")
	}
	if !strings.Contains(result, "    * Subsection 1.1.1") {
		t.Error("expected double-indented nested list item")
	}
}

func TestTofuguScraper_ScrapeArticle_Table(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title</h1>
	<article>
		<div class="main">
			<table>
				<tr>
					<th>Header 1</th>
					<th>Header 2</th>
				</tr>
				<tr>
					<td>Cell 1</td>
					<td>Cell 2</td>
				</tr>
			</table>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "| Header 1 ") {
		t.Error("expected table header 1")
	}
	if !strings.Contains(result, "| Header 2 ") {
		t.Error("expected table header 2")
	}
	if !strings.Contains(result, "| --- ") {
		t.Error("expected table separator row")
	}
	if !strings.Contains(result, "| Cell 1 ") {
		t.Error("expected table cell 1")
	}
	if !strings.Contains(result, "| Cell 2 ") {
		t.Error("expected table cell 2")
	}
}

func TestTofuguScraper_ScrapeArticle_TableWithoutHeaders(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title</h1>
	<article>
		<div class="main">
			<table>
				<tr>
					<td>Cell A</td>
					<td>Cell B</td>
				</tr>
				<tr>
					<td>Cell C</td>
					<td>Cell D</td>
				</tr>
			</table>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have a dummy header row and separator
	if !strings.Contains(result, "| |") {
		t.Errorf("expected dummy header row, got: %q", result)
	}
	if !strings.Contains(result, "|---") {
		t.Errorf("expected separator row, got: %q", result)
	}
	if !strings.Contains(result, "| Cell A ") {
		t.Error("expected table cell A")
	}
}

func TestTofuguScraper_ScrapeArticle_Blockquote(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title</h1>
	<article>
		<div class="main">
			<blockquote>This is a quoted text</blockquote>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "> This is a quoted text") {
		t.Errorf("expected blockquote with > prefix, got: %q", result)
	}
}

func TestTofuguScraper_ScrapeArticle_MultilineBlockquote(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title</h1>
	<article>
		<div class="main">
			<blockquote>Line one
Line two
Line three</blockquote>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "> Line one\n> Line two\n> Line three") {
		t.Errorf("expected multiline blockquote, got: %q", result)
	}
}

func TestTofuguScraper_ScrapeArticle_InvalidURL(t *testing.T) {
	scraper := &TofuguScraper{}
	_, err := scraper.ScrapeArticle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestTofuguScraper_ScrapeArticle_EmptyContent(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Empty Article</h1>
	<article>
		<div class="main">
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "# Empty Article\n\n"
	if result != expected {
		t.Errorf("ScrapeArticle() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeArticle_NoTitle(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article>
		<div class="main">
			<h2>Section</h2>
			<p>Content</p>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should still process content even without title
	if !strings.Contains(result, "## Section") {
		t.Error("expected h2 heading")
	}
	if !strings.Contains(result, "Content") {
		t.Error("expected paragraph content")
	}
	// Should not have h1 marker
	if strings.HasPrefix(result, "# ") {
		t.Error("should not have h1 when no h1.article-title in HTML")
	}
}

func TestTofuguScraper_ScrapeArticle_HeadingsOrder(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title</h1>
	<article>
		<div class="main">
			<h2>First Section</h2>
			<p>First paragraph</p>
			<h2>Second Section</h2>
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

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check order
	titleIdx := strings.Index(result, "# Title")
	firstSectionIdx := strings.Index(result, "## First Section")
	firstParaIdx := strings.Index(result, "First paragraph")
	secondSectionIdx := strings.Index(result, "## Second Section")
	secondParaIdx := strings.Index(result, "Second paragraph")

	if titleIdx == -1 || firstSectionIdx == -1 || firstParaIdx == -1 || secondSectionIdx == -1 || secondParaIdx == -1 {
		t.Fatalf("expected all content to be present, got: %q", result)
	}

	if titleIdx >= firstSectionIdx {
		t.Error("title should appear before first section")
	}
	if firstSectionIdx >= firstParaIdx {
		t.Error("first section should appear before first paragraph")
	}
	if firstParaIdx >= secondSectionIdx {
		t.Error("first paragraph should appear before second section")
	}
	if secondSectionIdx >= secondParaIdx {
		t.Error("second section should appear before second paragraph")
	}
}

func TestTofuguScraper_ScrapeArticle_MixedContent(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Japanese Grammar Guide</h1>
	<div class="article-header-elements">
		<ul class="meta">
			<li>Grammar</li>
		</ul>
	</div>
	<article>
		<div class="main">
			<h2>Introduction</h2>
			<p>Welcome to Japanese grammar.</p>
			<h3>Prerequisites</h3>
			<ol>
				<li>Basic hiragana</li>
				<li>Basic katakana</li>
			</ol>
			<ul class="example-sentence">
				<li>これは本です</li>
				<li>This is a book</li>
			</ul>
			<table>
				<tr><th>Particle</th><th>Usage</th></tr>
				<tr><td>は</td><td>Topic marker</td></tr>
			</table>
			<blockquote>Important note about grammar</blockquote>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectations := []struct {
		content string
		desc    string
	}{
		{"# Japanese Grammar Guide", "main title"},
		{"Grammar", "meta information"},
		{"## Introduction", "h2 heading"},
		{"Welcome to Japanese grammar.", "paragraph content"},
		{"### Prerequisites", "h3 heading"},
		{"1. Basic hiragana", "ordered list item"},
		{"Example", "example label"},
		{"- Japanese:", "Japanese label"},
		{"これは本です", "Japanese example sentence"},
		{"| Particle ", "table header"},
		{"| --- ", "table separator"},
		{"| は ", "table cell"},
		{"> Important note about grammar", "blockquote"},
	}

	for _, exp := range expectations {
		if !strings.Contains(result, exp.content) {
			t.Errorf("expected %s (%q) in result, got: %q", exp.desc, exp.content, result)
		}
	}
}

func TestTofuguScraper_ScrapeArticle_OnlyDirectChildren(t *testing.T) {
	// Tests that only direct children of div.main are processed
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title</h1>
	<article>
		<div class="main">
			<div class="wrapper">
				<h2>Nested H2</h2>
			</div>
			<h2>Direct H2</h2>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
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

func TestTofuguScraper_ScrapeArticle_TableWithNewlinesInCells(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title</h1>
	<article>
		<div class="main">
			<table>
				<tr>
					<th>Header</th>
				</tr>
				<tr>
					<td>Line 1
Line 2</td>
				</tr>
			</table>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Newlines in cells should be converted to semicolons
	if !strings.Contains(result, "Line 1; Line 2") {
		t.Errorf("expected newlines in table cell to be converted to semicolons, got: %q", result)
	}
}

func TestTofuguScraper_ScrapeArticle_WhitespaceHandling(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">   Title With Spaces   </h1>
	<article>
		<div class="main">
			<p>   Paragraph with extra spaces   </p>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Title should be trimmed
	if !strings.Contains(result, "# Title With Spaces\n\n") {
		t.Errorf("expected trimmed title, got: %q", result)
	}
}

func TestTofuguScraper_ScrapeTitle_Basic(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Tofugu Article Title</h1>
	<article><div class="main"></div></article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Tofugu Article Title"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeTitle_WithWhitespace(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">   Title With Whitespace   </h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Title With Whitespace"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeTitle_WithNewlines(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
<h1 class="article-title">
Title With Newlines
</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Title With Newlines"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeTitle_NoArticleTitle(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Regular H1 Without Class</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// When no h1.article-title exists, title should be empty
	if result != "" {
		t.Errorf("ScrapeTitle() = %q, want empty string", result)
	}
}

func TestTofuguScraper_ScrapeTitle_MultipleH1ArticleTitle(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">First Title</h1>
	<h1 class="article-title">Second Title</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// When multiple h1.article-title exist, the last one overwrites previous ones
	expected := "Second Title"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeTitle_EmptyH1(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title"></h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("ScrapeTitle() = %q, want empty string", result)
	}
}

func TestTofuguScraper_ScrapeTitle_InvalidURL(t *testing.T) {
	scraper := &TofuguScraper{}
	_, err := scraper.ScrapeTitle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestTofuguScraper_ScrapeTitle_SpecialCharacters(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title with &amp; special &lt;characters&gt;</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
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

func TestTofuguScraper_ScrapeTitle_WithNestedElements(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title"><span>Nested</span> <strong>Title</strong></h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
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

func TestTofuguScraper_ScrapeTitle_JapaneseCharacters(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">日本語の文法ガイド</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "日本語の文法ガイド"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeTitle_WithTildeCharacter(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">〜から〜まで: From 〜 To 〜</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ScrapeTitle should preserve the tilde character (〜)
	// Note: ScrapeFilename removes it, but ScrapeTitle should keep it
	expected := "〜から〜まで: From 〜 To 〜"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeTitle_IgnoresRegularH1(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Regular H1</h1>
	<h1 class="article-title">Article Title</h1>
	<h1 class="other-class">Other H1</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should only capture h1.article-title, not regular h1 or h1 with other classes
	expected := "Article Title"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_Basic(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Japanese Grammar Guide</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "japanese_grammar_guide"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_ConvertsToLowercase(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">LEARN JAPANESE NOW</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "learn_japanese_now"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_RemovesTilde(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">〜から〜まで From To</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Japanese tilde 〜 should be removed
	expected := "from_to"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_ReplacesSlashWithHyphen(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">This/That Guide</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Slash should be replaced with hyphen
	expected := "this-that_guide"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_RemovesParentheses(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Guide (Complete Edition)</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Parentheses should be removed
	expected := "guide_complete_edition"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_RemovesJapaneseCharacters(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">日本語 Japanese Guide</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Japanese characters should be removed, only ASCII alphanumeric remains
	expected := "japanese_guide"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_RemovesLeadingUnderscore(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">日本語 Guide</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Leading underscore (from removed Japanese chars + space) should be trimmed
	expected := "guide"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_PreservesNumbersAndDots(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Version 2.0 Guide</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Numbers and dots should be preserved
	expected := "version_2.0_guide"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_PreservesHyphens(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Step-by-Step Guide</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Hyphens should be preserved
	expected := "step-by-step_guide"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_ComplexTitle(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">〜てから (After Doing) - Grammar Guide</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 〜 removed, てから removed (Japanese), parentheses removed, - preserved
	expected := "after_doing_-_grammar_guide"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_SingleWord(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Hiragana</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "hiragana"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

// Tests for removeNonFilenameChars helper function

func TestRemoveNonFilenameChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only alphanumeric",
			input:    "HelloWorld123",
			expected: "HelloWorld123",
		},
		{
			name:     "with underscores",
			input:    "hello_world",
			expected: "hello_world",
		},
		{
			name:     "with hyphens",
			input:    "hello-world",
			expected: "hello-world",
		},
		{
			name:     "with dots",
			input:    "file.txt",
			expected: "file.txt",
		},
		{
			name:     "removes japanese characters",
			input:    "hello日本語world",
			expected: "helloworld",
		},
		{
			name:     "removes tilde",
			input:    "hello〜world",
			expected: "helloworld",
		},
		{
			name:     "removes parentheses",
			input:    "hello(world)",
			expected: "helloworld",
		},
		{
			name:     "removes spaces",
			input:    "hello world",
			expected: "helloworld",
		},
		{
			name:     "removes special characters",
			input:    "hello@#$%world",
			expected: "helloworld",
		},
		{
			name:     "removes colons",
			input:    "title: subtitle",
			expected: "titlesubtitle",
		},
		{
			name:     "removes slashes",
			input:    "path/to/file",
			expected: "pathtofile",
		},
		{
			name:     "mixed allowed characters",
			input:    "hello_world-test.txt",
			expected: "hello_world-test.txt",
		},
		{
			name:     "only japanese characters",
			input:    "日本語のみ",
			expected: "",
		},
		{
			name:     "preserves uppercase",
			input:    "HelloWorld",
			expected: "HelloWorld",
		},
		{
			name:     "complex mix",
			input:    "〜てから (After) - Guide_v1.0",
			expected: "After-Guide_v1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeNonFilenameChars(tt.input)
			if result != tt.expected {
				t.Errorf("removeNonFilenameChars(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Tests for getBasenameFromURL helper function

func TestGetBasenameFromURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "article",
		},
		{
			name:     "simple URL",
			input:    "https://example.com/article-name",
			expected: "article-name",
		},
		{
			name:     "URL with trailing slash",
			input:    "https://example.com/article-name/",
			expected: "article-name",
		},
		{
			name:     "URL with multiple path segments",
			input:    "https://example.com/blog/2024/my-article",
			expected: "my-article",
		},
		{
			name:     "URL with query string",
			input:    "https://example.com/article?id=123",
			expected: "article?id=123",
		},
		{
			name:     "URL with file extension",
			input:    "https://example.com/document.html",
			expected: "document.html",
		},
		{
			name:     "root URL only",
			input:    "https://example.com/",
			expected: "example.com",
		},
		{
			name:     "URL without scheme",
			input:    "example.com/article",
			expected: "article",
		},
		{
			name:     "single path segment",
			input:    "article",
			expected: "article",
		},
		{
			name:     "multiple trailing slashes collapse",
			input:    "https://example.com/article//",
			expected: "article",
		},
		{
			name:     "deep nested path",
			input:    "https://example.com/a/b/c/d/e/last-segment",
			expected: "last-segment",
		},
		{
			name:     "URL with Japanese characters",
			input:    "https://example.com/日本語記事",
			expected: "日本語記事",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBasenameFromURL(tt.input)
			if result != tt.expected {
				t.Errorf("getBasenameFromURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Additional tests for TofuguScraper.ScrapeFilename edge cases

func TestTofuguScraper_ScrapeFilename_InvalidURL(t *testing.T) {
	scraper := &TofuguScraper{}
	_, err := scraper.ScrapeFilename("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestTofuguScraper_ScrapeFilename_EmptyTitleFallsBackToURL(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title"></h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL + "/my-article-path")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Empty title should fall back to URL basename
	expected := "my-article-path"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_OnlyJapaneseCharactersFallsBackToURL(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">日本語のみ</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL + "/fallback-name")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Title with only Japanese chars becomes empty, should fall back to URL basename
	expected := "fallback-name"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_MultipleSpaces(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title   With   Multiple   Spaces</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Multiple spaces become multiple underscores
	expected := "title___with___multiple___spaces"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_MultipleSlashes(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Either/Or/Both</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Multiple slashes become multiple hyphens
	expected := "either-or-both"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_NestedParentheses(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Guide (With (Nested) Parens)</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// All parentheses removed
	expected := "guide_with_nested_parens"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_MultipleTildes(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">〜から〜まで〜</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL + "/kara-made")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// All tildes and Japanese chars removed, falls back to URL
	expected := "kara-made"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_MixedCasePreservedThenLowercased(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">CamelCase And UPPERCASE</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// All characters converted to lowercase
	expected := "camelcase_and_uppercase"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_SpecialCharactersRemoved(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Title: With @Special# Characters!</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Special characters like : @ # ! are removed
	expected := "title_with_special_characters"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_NoArticleTitleFallsBackToURL(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Regular H1 Without Class</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL + "/url-fallback")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No h1.article-title means empty title, falls back to URL basename
	expected := "url-fallback"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_OnlySpecialCharactersFallsBackToURL(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">@#$%^&*!</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL + "/special-chars-fallback")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Title with only special chars becomes empty, falls back to URL basename
	expected := "special-chars-fallback"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_LeadingSpaceBecomesLeadingUnderscore(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title"> Leading Space Title</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Note: trimSpacesAndLineBreaks is called during ScrapeTitle,
	// so leading space should be trimmed before filename processing
	expected := "leading_space_title"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_TrailingSpace(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">Trailing Space Title </h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Trailing space should be trimmed during ScrapeTitle
	expected := "trailing_space_title"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_NumericOnlyTitle(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">12345</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Numeric titles should be preserved
	expected := "12345"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTofuguScraper_ScrapeFilename_WithQuestionAndAmpersand(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 class="article-title">What? Why & How</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TofuguScraper{}
	result, err := scraper.ScrapeFilename(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ? and & are removed, spaces around & become underscores (two consecutive)
	expected := "what_why__how"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}
