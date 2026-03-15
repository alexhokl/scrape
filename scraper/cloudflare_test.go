package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCloudflareScraper_ScrapeArticle_Title(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1 pt1 pt3-l mb1">Introducing Markdown for Agents</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1"></div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "# Introducing Markdown for Agents") {
		t.Errorf("expected title in result, got: %q", result)
	}
}

func TestCloudflareScraper_ScrapeArticle_Headings(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Main Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<div class="flex anchor relative">
						<h2 id="section-one">Section One</h2>
						<a href="#section-one" aria-hidden="true"><svg></svg></a>
					</div>
					<div class="flex anchor relative">
						<h3 id="subsection">Subsection</h3>
						<a href="#subsection" aria-hidden="true"><svg></svg></a>
					</div>
					<div class="flex anchor relative">
						<h4 id="detail">Detail</h4>
						<a href="#detail" aria-hidden="true"><svg></svg></a>
					</div>
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "## Section One") {
		t.Errorf("expected h2 to be converted to ## heading, got: %q", result)
	}
	if !strings.Contains(result, "### Subsection") {
		t.Errorf("expected h3 to be converted to ### heading, got: %q", result)
	}
	if !strings.Contains(result, "#### Detail") {
		t.Errorf("expected h4 to be converted to #### heading, got: %q", result)
	}
}

func TestCloudflareScraper_ScrapeArticle_Paragraph(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<p>This is a paragraph of text about Cloudflare.</p>
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "This is a paragraph of text about Cloudflare.") {
		t.Errorf("expected paragraph content, got: %q", result)
	}
}

func TestCloudflareScraper_ScrapeArticle_InlineCode(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<p>Use the <code>text/markdown</code> content type.</p>
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "`text/markdown`") {
		t.Errorf("expected inline code backtick wrapping, got: %q", result)
	}
}

func TestCloudflareScraper_ScrapeArticle_Link(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<p>Learn more about <a href="https://en.wikipedia.org/wiki/Markdown"><u>Markdown</u></a> format.</p>
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "[Markdown](https://en.wikipedia.org/wiki/Markdown)") {
		t.Errorf("expected markdown link, got: %q", result)
	}
}

func TestCloudflareScraper_ScrapeArticle_BoldAndItalic(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<p><b>That's a 80% reduction</b> in token usage, the <i>lingua franca</i> of AI.</p>
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "**That's a 80% reduction**") {
		t.Errorf("expected bold text, got: %q", result)
	}
	if !strings.Contains(result, "*lingua franca*") {
		t.Errorf("expected italic text, got: %q", result)
	}
}

func TestCloudflareScraper_ScrapeArticle_CodeBlock(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<pre class="language-javascript"><code class="language-javascript">const r = await fetch(url);</code></pre>
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "```javascript") {
		t.Errorf("expected code block with javascript language, got: %q", result)
	}
	if !strings.Contains(result, "const r = await fetch(url);") {
		t.Errorf("expected code block content, got: %q", result)
	}
}

func TestCloudflareScraper_ScrapeArticle_CodeBlockNoLanguage(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<pre><code>plain code block</code></pre>
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "```\nplain code block\n```") {
		t.Errorf("expected code block without language, got: %q", result)
	}
}

func TestCloudflareScraper_ScrapeArticle_UnorderedList(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<ul>
						<li><p>First item</p></li>
						<li><p>Second item</p></li>
					</ul>
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
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

func TestCloudflareScraper_ScrapeArticle_OrderedList(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<ol>
						<li>Step one</li>
						<li>Step two</li>
						<li>Step three</li>
					</ol>
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
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

func TestCloudflareScraper_ScrapeArticle_UnorderedListWithLinks(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<ul>
						<li><p>Workers AI <a href="https://developers.cloudflare.com/workers-ai/"><u>AI.toMarkdown()</u></a> supports documents.</p></li>
					</ul>
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "[AI.toMarkdown()](https://developers.cloudflare.com/workers-ai/)") {
		t.Errorf("expected link in list item, got: %q", result)
	}
}

func TestCloudflareScraper_ScrapeArticle_Figure(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<figure class="kg-card kg-image-card">
						<img src="https://cf-assets.example.com/image.png" alt="Diagram" class="kg-image" width="800" height="400" loading="lazy"/>
					</figure>
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "![Diagram](https://cf-assets.example.com/image.png)") {
		t.Errorf("expected markdown image, got: %q", result)
	}
}

func TestCloudflareScraper_ScrapeArticle_Blockquote(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<blockquote>This is a quoted statement.</blockquote>
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "> This is a quoted statement.") {
		t.Errorf("expected blockquote content, got: %q", result)
	}
}

func TestCloudflareScraper_ScrapeArticle_ExcludesContentOutsideArticle(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<nav>Navigation content that should not appear</nav>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<p>Real content</p>
				</div>
			</section>
		</article>
	</main>
	<footer>Footer content that should not appear</footer>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "Navigation content that should not appear") {
		t.Error("nav content should not be included")
	}
	if strings.Contains(result, "Footer content that should not appear") {
		t.Error("footer content should not be included")
	}
	if !strings.Contains(result, "Real content") {
		t.Error("expected article body content to be present")
	}
}

func TestCloudflareScraper_ScrapeArticle_SkipsBoilerplateFooter(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Title</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<p>Real article content.</p>
				</div>
			</section>
			<section class="post-full-content flex flex-row flex-wrap mw7 center mb4">
				<div class="post-content lh-copy w-100 gray1 bt b--gray8 pt4">
					Cloudflare's connectivity cloud protects entire corporate networks.
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Real article content.") {
		t.Errorf("expected real article content, got: %q", result)
	}
	if strings.Contains(result, "connectivity cloud") {
		t.Errorf("boilerplate footer should not be included, got: %q", result)
	}
}

func TestCloudflareScraper_ScrapeArticle_InvalidURL(t *testing.T) {
	scraper := &CloudflareScraper{}
	_, err := scraper.ScrapeArticle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestCloudflareScraper_ScrapeTitle_Basic(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1 pt1 pt3-l mb1">Introducing Markdown for Agents</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1"></div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Introducing Markdown for Agents"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestCloudflareScraper_ScrapeTitle_WithWhitespace(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">  Title With Whitespace  </h1>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Title With Whitespace"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestCloudflareScraper_ScrapeTitle_NoArticle(t *testing.T) {
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

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("ScrapeTitle() = %q, want empty string", result)
	}
}

func TestCloudflareScraper_ScrapeTitle_InvalidURL(t *testing.T) {
	scraper := &CloudflareScraper{}
	_, err := scraper.ScrapeTitle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestCloudflareScraper_ScrapeFilename_UsesURLBasename(t *testing.T) {
	scraper := &CloudflareScraper{}

	result, err := scraper.ScrapeFilename("https://blog.cloudflare.com/markdown-for-agents/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "markdown-for-agents"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestCloudflareScraper_ScrapeFilename_NoTrailingSlash(t *testing.T) {
	scraper := &CloudflareScraper{}

	result, err := scraper.ScrapeFilename("https://blog.cloudflare.com/markdown-for-agents")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "markdown-for-agents"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestCloudflareScraper_ScrapeArticle_MixedContent(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="post">
		<article class="post-full mw-100 ph3 ph0-l fs-20px">
			<h1 class="f6 f7-l fw4 gray1">Getting Started with Cloudflare</h1>
			<section class="post-full-content">
				<div class="post-content lh-copy gray1">
					<div class="flex anchor relative">
						<h2 id="introduction">Introduction</h2>
						<a href="#introduction" aria-hidden="true"><svg></svg></a>
					</div>
					<p>Welcome to <b>Cloudflare</b>.</p>
					<ul>
						<li><p>Easy to set up</p></li>
						<li><p>Secure by default</p></li>
					</ul>
					<pre class="language-shell"><code class="language-shell">curl https://api.cloudflare.com</code></pre>
					<figure class="kg-card kg-image-card">
						<img src="https://cf-assets.example.com/diagram.png" alt="Architecture" class="kg-image"/>
					</figure>
				</div>
			</section>
		</article>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &CloudflareScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectations := []struct {
		content string
		desc    string
	}{
		{"# Getting Started with Cloudflare", "h1 title"},
		{"## Introduction", "h2 heading"},
		{"**Cloudflare**", "bold text"},
		{"* Easy to set up", "list item 1"},
		{"* Secure by default", "list item 2"},
		{"```shell", "code block language"},
		{"curl https://api.cloudflare.com", "code block content"},
		{"![Architecture](https://cf-assets.example.com/diagram.png)", "image"},
	}

	for _, exp := range expectations {
		if !strings.Contains(result, exp.content) {
			t.Errorf("expected %s (%q) in result, got: %q", exp.desc, exp.content, result)
		}
	}
}
