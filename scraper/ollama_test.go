package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ollamaArticleWrapper wraps content in a minimal Ollama blog page structure.
func ollamaArticleWrapper(title, body string) string {
	return `<!DOCTYPE html>
<html>
<head><title>` + title + ` · Ollama Blog</title></head>
<body>
  <article class="mx-auto flex flex-1 max-w-2xl w-full flex-col space-y-3 px-6 py-16 md:px-0">
    <h1 class="text-4xl font-semibold tracking-tight">` + title + `</h1>
    <h2 class="text-neutral-500">February 16, 2026</h2>
    <section class="prose">` + body + `</section>
  </article>
</body>
</html>`
}

func newOllamaServer(t *testing.T, body string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return server
}

func TestOllamaScraper_ScrapeArticle_Title(t *testing.T) {
	page := ollamaArticleWrapper("Subagents and web search in Claude Code", `<p>Some content.</p>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "# Subagents and web search in Claude Code") {
		t.Errorf("expected h1 title, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_Paragraph(t *testing.T) {
	page := ollamaArticleWrapper("Test", `<p>Ollama now supports subagents and web search in Claude Code.</p>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Ollama now supports subagents and web search in Claude Code.") {
		t.Errorf("expected paragraph text, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_Headings(t *testing.T) {
	page := ollamaArticleWrapper("Test", `
<h2 id="get-started">Get started</h2>
<p>Content under h2.</p>
<h3 id="details">Details</h3>
<p>Content under h3.</p>
<h4 id="more">More</h4>
<p>Content under h4.</p>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "## Get started") {
		t.Errorf("expected h2 heading, got: %q", result)
	}
	if !strings.Contains(result, "### Details") {
		t.Errorf("expected h3 heading, got: %q", result)
	}
	if !strings.Contains(result, "#### More") {
		t.Errorf("expected h4 heading, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_CodeBlockWithLanguage(t *testing.T) {
	page := ollamaArticleWrapper("Test", `
<pre><code class="language-bash">ollama launch claude --model minimax-m2.5:cloud
</code></pre>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "```bash") {
		t.Errorf("expected fenced code block with bash language, got: %q", result)
	}
	if !strings.Contains(result, "ollama launch claude --model minimax-m2.5:cloud") {
		t.Errorf("expected code block content, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_CodeBlockNoLanguage(t *testing.T) {
	page := ollamaArticleWrapper("Test", `
<pre><code>&gt; spawn subagents to explore the auth flow

&gt; audit security issues in parallel
</code></pre>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "```\n") {
		t.Errorf("expected fenced code block without language, got: %q", result)
	}
	if !strings.Contains(result, "spawn subagents") {
		t.Errorf("expected code block content, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_InlineCode(t *testing.T) {
	page := ollamaArticleWrapper("Test", `<p>Run <code>ollama serve</code> to start the server.</p>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "`ollama serve`") {
		t.Errorf("expected inline code backtick wrapping, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_Image(t *testing.T) {
	page := ollamaArticleWrapper("Test", `
<p><img src="https://files.ollama.com/web_search_and_subagents.png" alt="Subagents and Web Search logo with Ollama" /></p>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "![Subagents and Web Search logo with Ollama](https://files.ollama.com/web_search_and_subagents.png)") {
		t.Errorf("expected markdown image, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_ImageWithEmptyAlt(t *testing.T) {
	page := ollamaArticleWrapper("Test", `<p><img src="https://files.ollama.com/diagram.png" alt="" /></p>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "![](https://files.ollama.com/diagram.png)") {
		t.Errorf("expected image with empty alt, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_UnorderedList(t *testing.T) {
	page := ollamaArticleWrapper("Test", `
<ul>
<li><code>minimax-m2.5:cloud</code></li>
<li><code>glm-5:cloud</code></li>
<li><code>kimi-k2.5:cloud</code></li>
</ul>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "* `minimax-m2.5:cloud`") {
		t.Errorf("expected list item with inline code, got: %q", result)
	}
	if !strings.Contains(result, "* `glm-5:cloud`") {
		t.Errorf("expected second list item, got: %q", result)
	}
	if !strings.Contains(result, "* `kimi-k2.5:cloud`") {
		t.Errorf("expected third list item, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_OrderedList(t *testing.T) {
	page := ollamaArticleWrapper("Test", `
<ol>
<li>Install Ollama</li>
<li>Pull a model</li>
<li>Run inference</li>
</ol>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "1. Install Ollama") {
		t.Errorf("expected ordered list item 1, got: %q", result)
	}
	if !strings.Contains(result, "2. Pull a model") {
		t.Errorf("expected ordered list item 2, got: %q", result)
	}
	if !strings.Contains(result, "3. Run inference") {
		t.Errorf("expected ordered list item 3, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_ListWithLinks(t *testing.T) {
	page := ollamaArticleWrapper("Test", `
<ul>
<li><a href="/blog/launch">ollama launch</a> for more integrations</li>
<li><a href="/blog/claude">Claude Code with Ollama</a> for basic setup</li>
</ul>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "[ollama launch](/blog/launch)") {
		t.Errorf("expected link in list item, got: %q", result)
	}
	if !strings.Contains(result, "[Claude Code with Ollama](/blog/claude)") {
		t.Errorf("expected second link in list item, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_Link(t *testing.T) {
	page := ollamaArticleWrapper("Test", `
<p>Ollama's <a href="/blog/web-search">web search</a> is now built in.</p>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "[web search](/blog/web-search)") {
		t.Errorf("expected markdown link, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_BoldAndItalic(t *testing.T) {
	page := ollamaArticleWrapper("Test", `
<p><strong>Ollama</strong> is an <em>open-source</em> tool.</p>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "**Ollama**") {
		t.Errorf("expected bold text, got: %q", result)
	}
	if !strings.Contains(result, "*open-source*") {
		t.Errorf("expected italic text, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_Blockquote(t *testing.T) {
	page := ollamaArticleWrapper("Test", `
<blockquote><p>No MCP servers or API keys required.</p></blockquote>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "> No MCP servers or API keys required.") {
		t.Errorf("expected blockquote, got: %q", result)
	}
}

func TestOllamaScraper_ScrapeArticle_ExcludesOutsideArticle(t *testing.T) {
	page := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
  <header>Navigation that should not appear</header>
  <article class="mx-auto">
    <h1>Blog Post Title</h1>
    <section class="prose">
      <p>Real content.</p>
    </section>
  </article>
  <footer>Footer that should not appear</footer>
</body>
</html>`
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "Navigation that should not appear") {
		t.Error("header content should be excluded")
	}
	if strings.Contains(result, "Footer that should not appear") {
		t.Error("footer content should be excluded")
	}
	if !strings.Contains(result, "Real content.") {
		t.Error("expected article body content")
	}
}

func TestOllamaScraper_ScrapeArticle_MixedContent(t *testing.T) {
	page := ollamaArticleWrapper("Subagents and web search in Claude Code", `
<p><img src="https://files.ollama.com/web_search_and_subagents.png" alt="Subagents and Web Search logo" /></p>
<p>Ollama now supports subagents and web search in Claude Code. No MCP servers or API keys required.</p>
<h2 id="get-started">Get started</h2>
<pre><code class="language-bash">ollama launch claude --model minimax-m2.5:cloud
</code></pre>
<p>It works with any model on Ollama's cloud.</p>
<h2 id="subagents">Subagents</h2>
<p><img src="https://files.ollama.com/subagent-1.png" alt="Screencapture of subagents." /></p>
<p>Subagents can run tasks in parallel, such as file search, code exploration, and research.</p>
<pre><code>&gt; spawn subagents to explore the auth flow
</code></pre>
<h2 id="recommended-models">Recommended cloud models</h2>
<ul>
<li><code>minimax-m2.5:cloud</code></li>
<li><code>glm-5:cloud</code></li>
</ul>
<h2 id="learn-more">Learn more</h2>
<ul>
<li><a href="/blog/launch">ollama launch</a> for more integrations</li>
<li><a href="/blog/web-search">Web search API</a> for standalone usage</li>
</ul>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectations := []struct {
		content string
		desc    string
	}{
		{"# Subagents and web search in Claude Code", "h1 title"},
		{"![Subagents and Web Search logo](https://files.ollama.com/web_search_and_subagents.png)", "first image"},
		{"Ollama now supports subagents and web search", "intro paragraph"},
		{"## Get started", "h2 heading"},
		{"```bash", "code block language"},
		{"ollama launch claude --model minimax-m2.5:cloud", "code block content"},
		{"![Screencapture of subagents.](https://files.ollama.com/subagent-1.png)", "second image"},
		{"Subagents can run tasks in parallel", "paragraph after image"},
		{"```\n", "code block without language"},
		{"spawn subagents to explore the auth flow", "bare code block content"},
		{"## Recommended cloud models", "third h2"},
		{"* `minimax-m2.5:cloud`", "list item with inline code"},
		{"## Learn more", "fourth h2"},
		{"[ollama launch](/blog/launch)", "list link"},
	}

	for _, exp := range expectations {
		if !strings.Contains(result, exp.content) {
			t.Errorf("expected %s (%q) in result, got: %q", exp.desc, exp.content, result)
		}
	}
}

func TestOllamaScraper_ScrapeArticle_InvalidURL(t *testing.T) {
	scraper := &OllamaScraper{}
	_, err := scraper.ScrapeArticle("http://invalid.localhost.test:99999/nonexistent")
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestOllamaScraper_ScrapeTitle_Basic(t *testing.T) {
	page := ollamaArticleWrapper("Subagents and web search in Claude Code", `<p>Content.</p>`)
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeTitle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Subagents and web search in Claude Code"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestOllamaScraper_ScrapeTitle_WithWhitespace(t *testing.T) {
	page := `<!DOCTYPE html>
<html><body>
  <article>
    <h1>  Title With Spaces  </h1>
    <section class="prose"><p>Content.</p></section>
  </article>
</body></html>`
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeTitle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Title With Spaces"
	if result != expected {
		t.Errorf("ScrapeTitle() = %q, want %q", result, expected)
	}
}

func TestOllamaScraper_ScrapeTitle_NoArticle(t *testing.T) {
	page := `<!DOCTYPE html>
<html><body>
  <h1>Outside Article</h1>
</body></html>`
	server := newOllamaServer(t, page)

	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeTitle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// h1 outside article should still be matched since the selector is "article h1"
	// — but here there's no article, so result should be empty.
	if result != "" {
		t.Errorf("ScrapeTitle() = %q, want empty string", result)
	}
}

func TestOllamaScraper_ScrapeTitle_InvalidURL(t *testing.T) {
	scraper := &OllamaScraper{}
	_, err := scraper.ScrapeTitle("http://invalid.localhost.test:99999/nonexistent")
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestOllamaScraper_ScrapeFilename_Basic(t *testing.T) {
	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeFilename("https://ollama.com/blog/web-search-subagents-claude-code")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "web-search-subagents-claude-code"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}

func TestOllamaScraper_ScrapeFilename_WithTrailingSlash(t *testing.T) {
	scraper := &OllamaScraper{}
	result, err := scraper.ScrapeFilename("https://ollama.com/blog/web-search-subagents-claude-code/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "web-search-subagents-claude-code"
	if result != expected {
		t.Errorf("ScrapeFilename() = %q, want %q", result, expected)
	}
}
