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
