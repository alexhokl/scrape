package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGrafanaScraper_ScrapeArticle_Title(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>OpenTelemetry with Prometheus</h1>
		<div class="rich-text"></div>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "# OpenTelemetry with Prometheus") {
		t.Errorf("expected title in result, got: %q", result)
	}
}

func TestGrafanaScraper_ScrapeArticle_Headings(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>Main Title</h1>
		<div class="rich-text">
			<h2 class="scroll-mt-20 mt-8 group">
				<a aria-label="Link to section" href="#section-one"></a>
				<span>Section One</span>
			</h2>
			<h3 class="scroll-mt-20 mt-8 group">
				<a aria-label="Link to section" href="#sub-section"></a>
				<span>Sub Section</span>
			</h3>
			<h4 class="scroll-mt-20 mt-8 group">
				<a aria-label="Link to section" href="#sub-sub"></a>
				<span>Sub Sub Section</span>
			</h4>
		</div>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "## Section One") {
		t.Errorf("expected h2 to be converted to ## heading, got: %q", result)
	}
	if !strings.Contains(result, "### Sub Section") {
		t.Errorf("expected h3 to be converted to ### heading, got: %q", result)
	}
	if !strings.Contains(result, "#### Sub Sub Section") {
		t.Errorf("expected h4 to be converted to #### heading, got: %q", result)
	}
}

func TestGrafanaScraper_ScrapeArticle_HeadingFallbackText(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>Main Title</h1>
		<div class="rich-text">
			<h2>Plain Heading</h2>
		</div>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "## Plain Heading") {
		t.Errorf("expected h2 fallback text, got: %q", result)
	}
}

func TestGrafanaScraper_ScrapeArticle_Paragraph(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>Title</h1>
		<div class="rich-text">
			<p class="mt-4 leading-relaxed">This is a paragraph of text.</p>
		</div>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "This is a paragraph of text.") {
		t.Errorf("expected paragraph content, got: %q", result)
	}
}

func TestGrafanaScraper_ScrapeArticle_UnorderedList(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>Title</h1>
		<div class="rich-text">
			<ul>
				<li><div>First item</div></li>
				<li><div>Second item</div></li>
			</ul>
		</div>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
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

func TestGrafanaScraper_ScrapeArticle_OrderedList(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>Title</h1>
		<div class="rich-text">
			<ol>
				<li><div>Step one</div></li>
				<li><div>Step two</div></li>
				<li><div>Step three</div></li>
			</ol>
		</div>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
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

func TestGrafanaScraper_ScrapeArticle_CodeBlock(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>Title</h1>
		<div class="rich-text">
			<div class="relative">
				<button>Copy</button>
				<pre class="bg-gray-900 p-4 rounded-lg"><code>otlp:
  promote_resource_attributes:
    - service.name</code></pre>
			</div>
		</div>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "```") {
		t.Errorf("expected fenced code block, got: %q", result)
	}
	if !strings.Contains(result, "otlp:") {
		t.Errorf("expected code block content, got: %q", result)
	}
	if !strings.Contains(result, "service.name") {
		t.Errorf("expected code block content, got: %q", result)
	}
}

func TestGrafanaScraper_ScrapeArticle_InlineCode(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>Title</h1>
		<div class="rich-text">
			<p class="mt-4 leading-relaxed">Use the <code class="bg-gray-200 px-2 py-1 rounded text-sm font-mono">service.name</code> attribute.</p>
		</div>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "`service.name`") {
		t.Errorf("expected inline code backtick wrapping, got: %q", result)
	}
}

func TestGrafanaScraper_ScrapeArticle_InlineLink(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>Title</h1>
		<div class="rich-text">
			<p class="mt-4 leading-relaxed">See the <a href="https://grafana.com/docs" target="_blank">Grafana docs</a> for more.</p>
		</div>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "[Grafana docs](https://grafana.com/docs)") {
		t.Errorf("expected markdown link, got: %q", result)
	}
}

func TestGrafanaScraper_ScrapeArticle_ExcludesContentOutsideMain(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<header>Header content that should not appear</header>
	<main>
		<h1>Title</h1>
		<div class="rich-text">
			<p class="mt-4 leading-relaxed">Real content</p>
		</div>
	</main>
	<footer>Footer content that should not appear</footer>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "Header content that should not appear") {
		t.Error("header content should not be included")
	}
	if strings.Contains(result, "Footer content that should not appear") {
		t.Error("footer content should not be included")
	}
	if !strings.Contains(result, "Real content") {
		t.Error("expected article body content to be present")
	}
}

func TestGrafanaScraper_ScrapeArticle_InvalidURL(t *testing.T) {
	scraper := &GrafanaScraper{}
	_, err := scraper.ScrapeArticle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestGrafanaScraper_ScrapeTitle_Basic(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>OpenTelemetry with Prometheus</h1>
		<div class="rich-text"></div>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "OpenTelemetry with Prometheus"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestGrafanaScraper_ScrapeTitle_WithWhitespace(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>  Title With Whitespace  </h1>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Title With Whitespace"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestGrafanaScraper_ScrapeTitle_NoMain(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>Outside Main</h1>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// h1 outside <main> should not be captured
	if result != "" {
		t.Errorf("ScrapeTitle() = %q, want empty string", result)
	}
}

func TestGrafanaScraper_ScrapeTitle_InvalidURL(t *testing.T) {
	scraper := &GrafanaScraper{}
	_, err := scraper.ScrapeTitle("http://invalid.localhost.test:99999/nonexistent")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestGrafanaScraper_ScrapeFilename_UsesURLBasename(t *testing.T) {
	scraper := &GrafanaScraper{}

	result, err := scraper.ScrapeFilename("https://grafana.com/blog/opentelemetry-with-prometheus-better-integration-through-resource-attribute-promotion/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "opentelemetry-with-prometheus-better-integration-through-resource-attribute-promotion"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestGrafanaScraper_ScrapeFilename_NoTrailingSlash(t *testing.T) {
	scraper := &GrafanaScraper{}

	result, err := scraper.ScrapeFilename("https://grafana.com/blog/my-post")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "my-post"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestGrafanaScraper_ScrapeArticle_Image_NextjsProxy(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>Title</h1>
		<div class="rich-text">
			<p class="mt-4 leading-relaxed">Check out the dashboard:</p>
			<img alt="Lightweight APM for OpenTelemetry dashboard" loading="lazy" width="1000" height="1000" decoding="async" data-nimg="1" class="overflow-hidden shadow-md rounded-2xl" style="color:transparent" src="/mw/_next/image/?url=https%3A%2F%2Fs3.amazonaws.com%2Fa-us.storyblok.com%2Ff%2F1022730%2F58702db7b1%2Flightweight-amp.png&amp;w=3840&amp;q=75"/>
		</div>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "![Lightweight APM for OpenTelemetry dashboard](https://s3.amazonaws.com/a-us.storyblok.com/f/1022730/58702db7b1/lightweight-amp.png)") {
		t.Errorf("expected markdown image with decoded S3 URL, got: %q", result)
	}
}

func TestGrafanaScraper_ScrapeArticle_Image_DirectURL(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>Title</h1>
		<div class="rich-text">
			<img alt="A diagram" src="https://example.com/image.png"/>
		</div>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "![A diagram](https://example.com/image.png)") {
		t.Errorf("expected markdown image with direct URL, got: %q", result)
	}
}

func TestGrafanaImageURL_NextjsProxy(t *testing.T) {
	src := "/mw/_next/image/?url=https%3A%2F%2Fs3.amazonaws.com%2Fa-us.storyblok.com%2Ff%2F1022730%2F58702db7b1%2Flightweight-amp.png&w=3840&q=75"
	result := parseGrafanaImageURL(src)
	expected := "https://s3.amazonaws.com/a-us.storyblok.com/f/1022730/58702db7b1/lightweight-amp.png"
	if result != expected {
		t.Errorf("parseGrafanaImageURL() = %q, want %q", result, expected)
	}
}

func TestGrafanaImageURL_DirectURL(t *testing.T) {
	src := "https://example.com/image.png"
	result := parseGrafanaImageURL(src)
	if result != src {
		t.Errorf("parseGrafanaImageURL() = %q, want %q", result, src)
	}
}

func TestGrafanaScraper_ScrapeArticle_MixedContent(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main>
		<h1>Getting Started with Grafana</h1>
		<div class="rich-text">
			<h2 class="group">
				<a aria-label="Link to section" href="#intro"></a>
				<span>Introduction</span>
			</h2>
			<p class="mt-4 leading-relaxed">Welcome to <a href="https://grafana.com">Grafana</a>.</p>
			<ul>
				<li><div>Easy to use</div></li>
				<li><div>Open source</div></li>
			</ul>
			<div class="relative">
				<button>Copy</button>
				<pre><code>docker run -p 3000:3000 grafana/grafana</code></pre>
			</div>
			<p class="mt-4 leading-relaxed">Use <code class="bg-gray-200 px-2 py-1 rounded text-sm font-mono">service.name</code> for identification.</p>
		</div>
	</main>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := &GrafanaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectations := []struct {
		content string
		desc    string
	}{
		{"# Getting Started with Grafana", "h1 title"},
		{"## Introduction", "h2 heading"},
		{"[Grafana](https://grafana.com)", "inline link"},
		{"* Easy to use", "list item 1"},
		{"* Open source", "list item 2"},
		{"```", "code block fence"},
		{"docker run -p 3000:3000 grafana/grafana", "code block content"},
		{"`service.name`", "inline code"},
	}

	for _, exp := range expectations {
		if !strings.Contains(result, exp.content) {
			t.Errorf("expected %s (%q) in result, got: %q", exp.desc, exp.content, result)
		}
	}
}
