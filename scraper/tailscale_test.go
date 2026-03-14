package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTailscaleScraper_ScrapeArticle_Title(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Tailscale Services</h1>
		<div class="ts-prose"></div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "# Tailscale Services") {
		t.Errorf("expected title in result, got: %q", result)
	}
}

func TestTailscaleScraper_ScrapeArticle_Headings(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Main Title</h1>
		<div class="ts-prose">
			<h2><a href="#use-cases"><span id="inner-text">Use cases</span></a></h2>
			<h3><a href="#sub"><span id="inner-text">Subsection</span></a></h3>
			<h4><a href="#sub2"><span id="inner-text">Sub-subsection</span></a></h4>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "## Use cases") {
		t.Errorf("expected h2 to be converted to ## heading, got: %q", result)
	}
	if !strings.Contains(result, "### Subsection") {
		t.Errorf("expected h3 to be converted to ### heading, got: %q", result)
	}
	if !strings.Contains(result, "#### Sub-subsection") {
		t.Errorf("expected h4 to be converted to #### heading, got: %q", result)
	}
}

func TestTailscaleScraper_ScrapeArticle_HeadingsWithoutInnerText(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Main Title</h1>
		<div class="ts-prose">
			<h2>Plain Heading</h2>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "## Plain Heading") {
		t.Errorf("expected h2 fallback text, got: %q", result)
	}
}

func TestTailscaleScraper_ScrapeArticle_Paragraph(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Title</h1>
		<div class="ts-prose">
			<p>This is a paragraph of text.</p>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "This is a paragraph of text.") {
		t.Errorf("expected paragraph content, got: %q", result)
	}
}

func TestTailscaleScraper_ScrapeArticle_UnorderedList(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Title</h1>
		<div class="ts-prose">
			<ul>
				<li>First item</li>
				<li>Second item</li>
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

	scraper := &TailscaleScraper{}
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
}

func TestTailscaleScraper_ScrapeArticle_OrderedListWithCodeBlock(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Title</h1>
		<div class="ts-prose">
			<ol>
				<li>
					<p>Run the tailscale serve command to expose the service:</p>
					<div class="group relative overflow-hidden">
						<pre class="refractor language-shell"><code class="language-shell">tailscale serve --service=svc:web-server --https=443 127.0.0.1:8080</code></pre>
					</div>
				</li>
				<li>
					<p>Verify the service is running.</p>
				</li>
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

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "1. Run the tailscale serve command to expose the service:") {
		t.Errorf("expected first list item text, got: %q", result)
	}
	if !strings.Contains(result, "```shell") {
		t.Errorf("expected code block language fence, got: %q", result)
	}
	if !strings.Contains(result, "tailscale serve --service=svc:web-server --https=443 127.0.0.1:8080") {
		t.Errorf("expected code block content, got: %q", result)
	}
	if !strings.Contains(result, "2. Verify the service is running.") {
		t.Errorf("expected second list item text, got: %q", result)
	}
}

func TestTailscaleScraper_ScrapeArticle_OrderedList(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Title</h1>
		<div class="ts-prose">
			<ol>
				<li>Step one</li>
				<li>Step two</li>
				<li>Step three</li>
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

	scraper := &TailscaleScraper{}
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
	if !strings.Contains(result, "3. Step three") {
		t.Errorf("expected ordered list item 3, got: %q", result)
	}
}

func TestTailscaleScraper_ScrapeArticle_CodeBlock(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Title</h1>
		<div class="ts-prose">
			<div class="group relative overflow-hidden">
				<pre class="refractor language-shell"><code class="language-shell">tailscale up</code></pre>
			</div>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "```shell") {
		t.Errorf("expected code block with shell language, got: %q", result)
	}
	if !strings.Contains(result, "tailscale up") {
		t.Errorf("expected code block content, got: %q", result)
	}
	if !strings.Contains(result, "```") {
		t.Errorf("expected closing code fence, got: %q", result)
	}
}

func TestTailscaleScraper_ScrapeArticle_NoteDiv(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Title</h1>
		<div class="ts-prose">
			<div class="note">
				<p>This is an important note.</p>
			</div>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "> This is an important note.") {
		t.Errorf("expected note content as blockquote, got: %q", result)
	}
}

func TestTailscaleScraper_ScrapeArticle_NoteDivWithCodeBlock(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Title</h1>
		<div class="ts-prose">
			<div class="note">
				<p>Make sure to start the resource on a tailnet device.</p>
				<p>For example, you might start a basic web server:</p>
				<div class="group relative overflow-hidden">
					<pre class="refractor language-shell"><code class="language-shell"># Install globally
npm install -g http-server

# Then start a basic web server
http-server -p 8080</code></pre>
				</div>
				<p>This example starts a basic web server on port 8080.</p>
			</div>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "> Make sure to start the resource on a tailnet device.") {
		t.Errorf("expected first note paragraph as blockquote, got: %q", result)
	}
	if !strings.Contains(result, "```shell") {
		t.Errorf("expected code block language fence, got: %q", result)
	}
	if !strings.Contains(result, "# Install globally") {
		t.Errorf("expected code block content, got: %q", result)
	}
	if !strings.Contains(result, "http-server -p 8080") {
		t.Errorf("expected code block content, got: %q", result)
	}
	if !strings.Contains(result, "> This example starts a basic web server on port 8080.") {
		t.Errorf("expected last note paragraph as blockquote, got: %q", result)
	}
}

func TestTailscaleScraper_ScrapeArticle_InlineCode(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Title</h1>
		<div class="ts-prose">
			<p>Connect on port <code>8080</code> to get started.</p>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "port `8080`") {
		t.Errorf("expected inline code backtick wrapping, got: %q", result)
	}
}

func TestTailscaleScraper_ScrapeArticle_InlineCodeInNoteDiv(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Title</h1>
		<div class="ts-prose">
			<div class="note">
				<p>Make sure to start the resource on a tailnet device.</p>
				<p>For example, you might start a basic web server using Node.js and the <a href="https://www.npmjs.com/package/http-server"><code>http-server</code> package</a> that listens on port <code>8080</code>:</p>
				<div class="group relative overflow-hidden">
					<pre class="refractor language-shell"><code class="language-shell"># Install globally
npm install -g http-server

# Then start a basic web server
http-server -p 8080</code></pre>
				</div>
			</div>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "port `8080`") {
		t.Errorf("expected inline code backtick wrapping for 8080, got: %q", result)
	}
	if !strings.Contains(result, "`http-server` package") {
		t.Errorf("expected inline code backtick wrapping for http-server, got: %q", result)
	}
}

func TestTailscaleScraper_ScrapeArticle_ExcludesContentOutsideArticle(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<nav>Navigation content that should not appear</nav>
	<article id="main-content">
		<h1>Title</h1>
		<div class="ts-prose">
			<p>Real content</p>
		</div>
	</article>
	<aside>Sidebar content that should not appear</aside>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "Navigation content that should not appear") {
		t.Error("nav content should not be included")
	}
	if strings.Contains(result, "Sidebar content that should not appear") {
		t.Error("aside content should not be included")
	}
	if !strings.Contains(result, "Real content") {
		t.Error("expected article body content to be present")
	}
}

func TestTailscaleScraper_ScrapeArticle_InvalidURL(t *testing.T) {
	scraper := &TailscaleScraper{}
	_, err := scraper.ScrapeArticle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestTailscaleScraper_ScrapeTitle_Basic(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Tailscale Services</h1>
		<div class="ts-prose"></div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Tailscale Services"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestTailscaleScraper_ScrapeTitle_WithWhitespace(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>  Title With Whitespace  </h1>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Title With Whitespace"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestTailscaleScraper_ScrapeTitle_NoArticle(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Outside Article</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// h1 outside article#main-content should not be captured
	if result != "" {
		t.Errorf("ScrapeTitle() = %q, want empty string", result)
	}
}

func TestTailscaleScraper_ScrapeTitle_InvalidURL(t *testing.T) {
	scraper := &TailscaleScraper{}
	_, err := scraper.ScrapeTitle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestTailscaleScraper_ScrapeFilename_UsesURLBasename(t *testing.T) {
	scraper := &TailscaleScraper{}

	result, err := scraper.ScrapeFilename("https://tailscale.com/docs/features/tailscale-services")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "tailscale-services"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTailscaleScraper_ScrapeFilename_TrailingSlash(t *testing.T) {
	scraper := &TailscaleScraper{}

	result, err := scraper.ScrapeFilename("https://tailscale.com/docs/features/tailscale-services/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "tailscale-services"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestTailscaleScraper_ScrapeArticle_MixedContent(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article id="main-content">
		<h1>Getting Started</h1>
		<div class="ts-prose">
			<h2><a href="#intro"><span id="inner-text">Introduction</span></a></h2>
			<p>Welcome to Tailscale.</p>
			<ul>
				<li>Easy to set up</li>
				<li>Secure by default</li>
			</ul>
			<div class="note">
				<p>Requires an account.</p>
			</div>
			<div class="group relative overflow-hidden">
				<pre class="refractor language-shell"><code class="language-shell">curl -fsSL https://tailscale.com/install.sh | sh</code></pre>
			</div>
		</div>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &TailscaleScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectations := []struct {
		content string
		desc    string
	}{
		{"# Getting Started", "h1 title"},
		{"## Introduction", "h2 heading"},
		{"Welcome to Tailscale.", "paragraph"},
		{"* Easy to set up", "list item 1"},
		{"* Secure by default", "list item 2"},
		{"> Requires an account.", "note blockquote"},
		{"```shell", "code block language"},
		{"curl -fsSL https://tailscale.com/install.sh | sh", "code block content"},
	}

	for _, exp := range expectations {
		if !strings.Contains(result, exp.content) {
			t.Errorf("expected %s (%q) in result, got: %q", exp.desc, exp.content, result)
		}
	}
}
