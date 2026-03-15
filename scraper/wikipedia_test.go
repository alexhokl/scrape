package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWikipediaScraper_ScrapeArticle_Title(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Merkle tree - Wikipedia</title></head>
<body>
	<h1 id="firstHeading" class="firstHeading mw-first-heading">
		<span class="mw-page-title-main">Merkle tree</span>
	</h1>
	<div id="mw-content-text" class="mw-body-content">
		<div class="mw-parser-output">
			<p>A paragraph.</p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "# Merkle tree") {
		t.Errorf("expected title in result, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_Headings(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test Article</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<div class="mw-heading mw-heading2">
				<h2 id="Overview">Overview</h2>
				<span class="mw-editsection"><a href="/w/index.php?title=Test&amp;action=edit">[edit]</a></span>
			</div>
			<p>Overview text.</p>
			<div class="mw-heading mw-heading3">
				<h3 id="Details">Details</h3>
				<span class="mw-editsection"><a href="/w/index.php?title=Test&amp;action=edit">[edit]</a></span>
			</div>
			<p>Details text.</p>
			<div class="mw-heading mw-heading4">
				<h4 id="Specifics">Specifics</h4>
				<span class="mw-editsection"><a href="/w/index.php?title=Test&amp;action=edit">[edit]</a></span>
			</div>
			<p>Specifics text.</p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "## Overview") {
		t.Errorf("expected h2 heading, got: %q", result)
	}
	if !strings.Contains(result, "### Details") {
		t.Errorf("expected h3 heading, got: %q", result)
	}
	if !strings.Contains(result, "#### Specifics") {
		t.Errorf("expected h4 heading, got: %q", result)
	}
	// Edit section links should be excluded.
	if strings.Contains(result, "[edit]") {
		t.Errorf("edit section links should be excluded, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_Paragraph(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<p>In cryptography and computer science, a hash tree or Merkle tree is a tree in which every "leaf" node is labelled with the cryptographic hash of a data block.</p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "hash tree or Merkle tree") {
		t.Errorf("expected paragraph content, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_BoldAndItalic(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<p>A <b>Merkle tree</b> is a <i>hash tree</i> structure.</p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "**Merkle tree**") {
		t.Errorf("expected bold text, got: %q", result)
	}
	if !strings.Contains(result, "*hash tree*") {
		t.Errorf("expected italic text, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_InternalLink(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<p>Used in <a href="/wiki/Cryptography" title="Cryptography">cryptography</a> systems.</p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "[cryptography](https://en.wikipedia.org/wiki/Cryptography)") {
		t.Errorf("expected internal wiki link, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_ExternalLink(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<p>See the <a class="external text" href="https://example.com/doc">documentation</a> for details.</p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "[documentation](https://example.com/doc)") {
		t.Errorf("expected external link, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_InlineCode(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<p>The function <code>hash(x)</code> computes the digest.</p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "`hash(x)`") {
		t.Errorf("expected inline code, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_ExcludesCitationSuperscripts(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<p>Merkle trees are used in distributed systems.<sup class="reference" id="cite_ref-1"><a href="#cite_note-1">[1]</a></sup></p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "[1]") {
		t.Errorf("citation superscripts should be excluded, got: %q", result)
	}
	if !strings.Contains(result, "Merkle trees are used in distributed systems.") {
		t.Errorf("expected paragraph text without citations, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_ExcludesCitationNeeded(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<p>This claim is unverified.<sup class="noprint Inline-Template Template-Fact">[<i>citation needed</i>]</sup></p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "citation needed") {
		t.Errorf("citation needed tags should be excluded, got: %q", result)
	}
	if !strings.Contains(result, "This claim is unverified.") {
		t.Errorf("expected paragraph text, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_UnorderedList(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<ul>
				<li>First item</li>
				<li>Second item</li>
				<li>Third item</li>
			</ul>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "* First item") {
		t.Errorf("expected first list item, got: %q", result)
	}
	if !strings.Contains(result, "* Second item") {
		t.Errorf("expected second list item, got: %q", result)
	}
	if !strings.Contains(result, "* Third item") {
		t.Errorf("expected third list item, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_OrderedList(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<ol>
				<li>Step one</li>
				<li>Step two</li>
			</ol>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "1. Step one") {
		t.Errorf("expected ordered list item 1, got: %q", result)
	}
	if !strings.Contains(result, "2. Step two") {
		t.Errorf("expected ordered list item 2, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_SkipsReferencesSection(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<p>Main content paragraph.</p>
			<div class="mw-heading mw-heading2">
				<h2 id="References">References</h2>
			</div>
			<div class="reflist">
				<ol class="references">
					<li id="cite_note-1">Reference text that should not appear.</li>
				</ol>
			</div>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Main content paragraph.") {
		t.Errorf("expected main content, got: %q", result)
	}
	if strings.Contains(result, "## References") {
		t.Errorf("References heading should be excluded, got: %q", result)
	}
	if strings.Contains(result, "Reference text that should not appear") {
		t.Errorf("Reference list content should be excluded, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_SkipsSeeAlsoSection(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<p>Main content.</p>
			<div class="mw-heading mw-heading2">
				<h2 id="See_also">See also</h2>
			</div>
			<ul>
				<li><a href="/wiki/Hash_function">Hash function</a></li>
			</ul>
			<div class="mw-heading mw-heading2">
				<h2 id="External_links">External links</h2>
			</div>
			<ul>
				<li><a class="external text" href="https://example.com">Example</a></li>
			</ul>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Main content.") {
		t.Errorf("expected main content, got: %q", result)
	}
	if strings.Contains(result, "## See also") {
		t.Errorf("See also heading should be excluded, got: %q", result)
	}
	if strings.Contains(result, "## External links") {
		t.Errorf("External links heading should be excluded, got: %q", result)
	}
	if strings.Contains(result, "Hash function") {
		t.Errorf("See also content should be excluded, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_ContentAfterSkippedSectionResumes(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<div class="mw-heading mw-heading2">
				<h2 id="Overview">Overview</h2>
			</div>
			<p>Overview content.</p>
			<div class="mw-heading mw-heading2">
				<h2 id="See_also">See also</h2>
			</div>
			<ul><li>Skipped item</li></ul>
			<div class="mw-heading mw-heading2">
				<h2 id="Applications">Applications</h2>
			</div>
			<p>Applications content should appear.</p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "## Overview") {
		t.Errorf("expected Overview heading, got: %q", result)
	}
	if !strings.Contains(result, "## Applications") {
		t.Errorf("expected Applications heading after skipped section, got: %q", result)
	}
	if !strings.Contains(result, "Applications content should appear.") {
		t.Errorf("expected Applications content, got: %q", result)
	}
	if strings.Contains(result, "Skipped item") {
		t.Errorf("See also content should be excluded, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_Figure(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<figure typeof="mw:File/Thumb">
				<a href="/wiki/File:Hash_Tree.svg">
					<img class="mw-file-element" src="//upload.wikimedia.org/wikipedia/commons/thumb/9/95/Hash_Tree.svg/310px-Hash_Tree.svg.png" alt="Hash tree diagram" width="310" height="230"/>
				</a>
				<figcaption>A binary hash tree</figcaption>
			</figure>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "![Hash tree diagram](https://upload.wikimedia.org/") {
		t.Errorf("expected image with https protocol, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_PreCodeBlock(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<pre>def merkle_root(leaves):
    return hash(leaves)</pre>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "```\ndef merkle_root(leaves):\n    return hash(leaves)\n```") {
		t.Errorf("expected code block, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_Blockquote(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<blockquote>A hash tree allows efficient verification of the contents.</blockquote>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "> A hash tree allows efficient verification of the contents.") {
		t.Errorf("expected blockquote, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_Table(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<table class="wikitable">
				<tr><th>Property</th><th>Value</th></tr>
				<tr><td>Depth</td><td>O(log n)</td></tr>
				<tr><td>Leaf count</td><td>n</td></tr>
			</table>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "| Property | Value |") {
		t.Errorf("expected table header, got: %q", result)
	}
	if !strings.Contains(result, "| Depth | O(log n) |") {
		t.Errorf("expected table row, got: %q", result)
	}
	if !strings.Contains(result, "| --- | --- |") {
		t.Errorf("expected table separator, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_ExcludesNavigation(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<nav>Navigation that should not appear</nav>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<p>Real content.</p>
		</div>
	</div>
	<footer>Footer that should not appear</footer>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "Navigation that should not appear") {
		t.Error("nav content should be excluded")
	}
	if strings.Contains(result, "Footer that should not appear") {
		t.Error("footer content should be excluded")
	}
	if !strings.Contains(result, "Real content.") {
		t.Error("expected article body content")
	}
}

func TestWikipediaScraper_ScrapeArticle_ListWithLinks(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<ul>
				<li>Used in <a href="/wiki/Bitcoin">Bitcoin</a> for transaction verification</li>
				<li>Applied in <a href="/wiki/Git">Git</a> version control</li>
			</ul>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "[Bitcoin](https://en.wikipedia.org/wiki/Bitcoin)") {
		t.Errorf("expected wiki link in list, got: %q", result)
	}
	if !strings.Contains(result, "[Git](https://en.wikipedia.org/wiki/Git)") {
		t.Errorf("expected wiki link in list, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_InvalidURL(t *testing.T) {
	scraper := &WikipediaScraper{}
	_, err := scraper.ScrapeArticle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestWikipediaScraper_ScrapeArticle_MixedContent(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Merkle tree - Wikipedia</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Merkle tree</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<p>A <b>Merkle tree</b> is a <a href="/wiki/Hash_function">hash</a>-based data structure.<sup class="reference" id="cite_ref-1"><a href="#cite_note-1">[1]</a></sup></p>
			<div class="mw-heading mw-heading2">
				<h2 id="Overview">Overview</h2>
				<span class="mw-editsection"><a href="/w/index.php?title=Test&amp;action=edit">[edit]</a></span>
			</div>
			<p>Merkle trees are typically used in distributed systems.</p>
			<ul>
				<li><a href="/wiki/Bitcoin">Bitcoin</a> uses Merkle trees</li>
				<li><a href="/wiki/Git">Git</a> uses a similar structure</li>
			</ul>
			<figure typeof="mw:File/Thumb">
				<a href="/wiki/File:Hash_Tree.svg">
					<img class="mw-file-element" src="//upload.wikimedia.org/wikipedia/commons/tree.png" alt="Tree diagram" width="300"/>
				</a>
			</figure>
			<div class="mw-heading mw-heading2">
				<h2 id="References">References</h2>
			</div>
			<div class="reflist">
				<ol class="references">
					<li>Should not appear</li>
				</ol>
			</div>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectations := []struct {
		content string
		desc    string
	}{
		{"# Merkle tree", "h1 title"},
		{"**Merkle tree**", "bold text"},
		{"[hash](https://en.wikipedia.org/wiki/Hash_function)", "internal link"},
		{"## Overview", "h2 heading"},
		{"Merkle trees are typically used in distributed systems.", "paragraph"},
		{"[Bitcoin](https://en.wikipedia.org/wiki/Bitcoin)", "list item link"},
		{"![Tree diagram](https://upload.wikimedia.org/", "image"},
	}

	for _, exp := range expectations {
		if !strings.Contains(result, exp.content) {
			t.Errorf("expected %s (%q) in result, got: %q", exp.desc, exp.content, result)
		}
	}

	exclusions := []struct {
		content string
		desc    string
	}{
		{"[1]", "citation superscript"},
		{"[edit]", "edit section link"},
		{"## References", "references heading"},
		{"Should not appear", "reference list content"},
	}

	for _, exc := range exclusions {
		if strings.Contains(result, exc.content) {
			t.Errorf("%s should be excluded, got: %q", exc.desc, result)
		}
	}
}

func TestWikipediaScraper_ScrapeTitle_Basic(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Merkle tree - Wikipedia</title></head>
<body>
	<h1 id="firstHeading" class="firstHeading mw-first-heading">
		<span class="mw-page-title-main">Merkle tree</span>
	</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Merkle tree"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestWikipediaScraper_ScrapeTitle_WithWhitespace(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading">  Merkle tree  </h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Merkle tree"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestWikipediaScraper_ScrapeTitle_NoHeading(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Not the right heading</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("ScrapeTitle() = %q, want empty string", result)
	}
}

func TestWikipediaScraper_ScrapeTitle_InvalidURL(t *testing.T) {
	scraper := &WikipediaScraper{}
	_, err := scraper.ScrapeTitle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestWikipediaScraper_ScrapeFilename_Basic(t *testing.T) {
	scraper := &WikipediaScraper{}

	result, err := scraper.ScrapeFilename("https://en.wikipedia.org/wiki/Merkle_tree")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Merkle_tree"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestWikipediaScraper_ScrapeFilename_WithTrailingSlash(t *testing.T) {
	scraper := &WikipediaScraper{}

	result, err := scraper.ScrapeFilename("https://en.wikipedia.org/wiki/Merkle_tree/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Merkle_tree"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestWikipediaScraper_ScrapeFilename_MultiWord(t *testing.T) {
	scraper := &WikipediaScraper{}

	result, err := scraper.ScrapeFilename("https://en.wikipedia.org/wiki/Hash_function")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Hash_function"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestWikipediaScraper_ScrapeArticle_ExcludesPortalBox(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<p>Real content.</p>
			<ul class="portalbox">
				<li>Portal link that should not appear</li>
			</ul>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "Portal link") {
		t.Errorf("portal box content should be excluded, got: %q", result)
	}
	if !strings.Contains(result, "Real content.") {
		t.Errorf("expected real content, got: %q", result)
	}
}

func TestWikipediaScraper_ScrapeArticle_ExcludesShortDescription(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1 id="firstHeading"><span class="mw-page-title-main">Test</span></h1>
	<div id="mw-content-text">
		<div class="mw-parser-output">
			<div class="shortdescription nomobile noexcerpt noprint searchaux">Tree data structure</div>
			<p>Real content about Merkle trees.</p>
		</div>
	</div>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &WikipediaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Real content about Merkle trees.") {
		t.Errorf("expected real content, got: %q", result)
	}
}
