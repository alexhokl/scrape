package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newUnseenJapanTestServer(body string) *httptest.Server {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<main id="content" class="uj-single" role="main">
		<article id="post-1" class="uj-article post-1 post">
			<figure class="uj-article__hero">
				<img src="https://example.com/hero.jpg" class="uj-article__hero-img" alt="Hero image">
				<figcaption class="uj-article__hero-caption">Picture: hero</figcaption>
			</figure>
			<header class="uj-article__header">
				<a class="uj-article__cat" href="https://example.com/category/history/">History</a>
				<h1 class="uj-article__title">Suzuki Masa: The Samurai&#8217;s Daughter</h1>
				<div class="uj-article__meta">
					<a class="uj-article__author-link" href="https://example.com/author/k/">K. Lou</a>
				</div>
			</header>
			<div class="uj-article__body entry-content">
` + body + `
			</div>
		</article>
	</main>
</body>
</html>`

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(html))
	}))
}

func TestUnseenJapanScraper_ScrapeArticle_Title(t *testing.T) {
	server := newUnseenJapanTestServer(``)
	defer server.Close()

	scraper := &UnseenJapanScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(result, "# Suzuki Masa: The Samurai’s Daughter\n\n") {
		t.Errorf("expected title at start of result, got: %q", result)
	}
}

func TestUnseenJapanScraper_ScrapeArticle_ExcludesHeroAndMeta(t *testing.T) {
	server := newUnseenJapanTestServer(`<p class="wp-block-paragraph">Body text.</p>`)
	defer server.Close()

	scraper := &UnseenJapanScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, unexpected := range []string{"hero.jpg", "Picture: hero", "K. Lou", "History"} {
		if strings.Contains(result, unexpected) {
			t.Errorf("expected %q to be excluded, got: %q", unexpected, result)
		}
	}
}

func TestUnseenJapanScraper_ScrapeArticle_Headings(t *testing.T) {
	server := newUnseenJapanTestServer(`
<h2 class="wp-block-heading">Section One</h2>
<h3 class="wp-block-heading">Subsection</h3>
<h4 class="wp-block-heading">Detail</h4>`)
	defer server.Close()

	scraper := &UnseenJapanScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, expected := range []string{"## Section One\n\n", "### Subsection\n\n", "#### Detail\n\n"} {
		if !strings.Contains(result, expected) {
			t.Errorf("expected %q in result, got: %q", expected, result)
		}
	}
}

func TestUnseenJapanScraper_ScrapeArticle_Paragraphs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text",
			input:    `<p class="wp-block-paragraph">Plain paragraph.</p>`,
			expected: "Plain paragraph.\n\n",
		},
		{
			name:     "link",
			input:    `<p class="wp-block-paragraph">Someone who, <a href="https://example.com/ginko/">much like Ogino Ginko</a>, helped.</p>`,
			expected: "Someone who, [much like Ogino Ginko](https://example.com/ginko/), helped.\n\n",
		},
		{
			name:     "emphasis with leading space",
			input:    `<p class="wp-block-paragraph">morning drama<em> Kaze, Kaoru</em> has</p>`,
			expected: "morning drama *Kaze, Kaoru* has\n\n",
		},
		{
			name:     "strong",
			input:    `<p class="wp-block-paragraph">A <strong>bold</strong> claim.</p>`,
			expected: "A **bold** claim.\n\n",
		},
		{
			name:     "code",
			input:    `<p class="wp-block-paragraph">Run <code>ls</code> now.</p>`,
			expected: "Run `ls` now.\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newUnseenJapanTestServer(tt.input)
			defer server.Close()

			scraper := &UnseenJapanScraper{}
			result, err := scraper.ScrapeArticle(server.URL)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !strings.Contains(result, tt.expected) {
				t.Errorf("expected %q in result, got: %q", tt.expected, result)
			}
		})
	}
}

func TestUnseenJapanScraper_ScrapeArticle_EmptyParagraph(t *testing.T) {
	server := newUnseenJapanTestServer(`<p class="wp-block-paragraph">   </p><p class="wp-block-paragraph">Text.</p>`)
	defer server.Close()

	scraper := &UnseenJapanScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "\n\n\n") {
		t.Errorf("expected empty paragraph to be skipped, got: %q", result)
	}
}

func TestUnseenJapanScraper_ScrapeArticle_LazyImage(t *testing.T) {
	server := newUnseenJapanTestServer(`
<div class="wp-block-image">
<figure class="aligncenter size-full"><img width="389" height="515" src="data:image/png;base64,iVBORw0KGgo=" alt="Portrait of Suzuki Masa." class="wp-image-92964 lazyload" data-src="https://media.example.com/image-17-7.jpg" /><noscript><img src="https://media.example.com/image-17-7.jpg" alt="Portrait of Suzuki Masa." /></noscript><figcaption class="wp-element-caption">One of the few surviving pictures.</figcaption></figure>
</div>`)
	defer server.Close()

	scraper := &UnseenJapanScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "![Portrait of Suzuki Masa.](https://media.example.com/image-17-7.jpg)\n\n*One of the few surviving pictures.*\n\n"
	if !strings.Contains(result, expected) {
		t.Errorf("expected %q in result, got: %q", expected, result)
	}
	if strings.Contains(result, "data:image") {
		t.Errorf("expected data: placeholder to be ignored, got: %q", result)
	}
}

func TestUnseenJapanScraper_ScrapeArticle_ImageWithoutLazyLoad(t *testing.T) {
	server := newUnseenJapanTestServer(`
<figure class="wp-block-image"><img src="https://media.example.com/plain.jpg" alt="Plain"></figure>`)
	defer server.Close()

	scraper := &UnseenJapanScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "![Plain](https://media.example.com/plain.jpg)\n\n"
	if !strings.Contains(result, expected) {
		t.Errorf("expected %q in result, got: %q", expected, result)
	}
	if strings.Contains(result, "**") || strings.Contains(result, "\n*\n") {
		t.Errorf("expected no caption, got: %q", result)
	}
}

func TestUnseenJapanScraper_ScrapeArticle_ImagePlaceholderOnly(t *testing.T) {
	server := newUnseenJapanTestServer(`
<div class="wp-block-image"><figure><img src="data:image/png;base64,iVBORw0KGgo=" alt="Missing"></figure></div>`)
	defer server.Close()

	scraper := &UnseenJapanScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "Missing") {
		t.Errorf("expected image without real URL to be skipped, got: %q", result)
	}
}

func TestUnseenJapanScraper_ScrapeArticle_SkipsNonArticleContent(t *testing.T) {
	server := newUnseenJapanTestServer(`
<p class="wp-block-paragraph">Before.</p>
<aside class="uj-cta uj-cta--newsletter">
	<h3 class="uj-cta__heading">Beyond the Japan Headlines</h3>
	<p class="uj-cta__body">Stay up to date via our free newsletter.</p>
</aside>
<script>var EXTRA_MERGE_TAGS = {};</script>
<aside class="uj-related-mini"><p class="uj-related-mini__eyebrow">Related</p></aside>
<p class="uj-member-pitch"><span class="uj-member-pitch__plain">Join Unseen Japan Insider today.</span></p>
<aside class="uj-inarticle-ad uj-no-ads"><h4>Advertisement</h4></aside>
<section class="uj-related"><h3 class="uj-related__heading">Keep reading</h3></section>
<span class="uj-no-ads"></span>
<p class="wp-block-paragraph">After.</p>`)
	defer server.Close()

	scraper := &UnseenJapanScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Before.\n\nAfter.\n\n") {
		t.Errorf("expected article paragraphs to be adjacent, got: %q", result)
	}
	for _, unexpected := range []string{
		"Beyond the Japan Headlines",
		"newsletter",
		"EXTRA_MERGE_TAGS",
		"Related",
		"Insider",
		"Advertisement",
		"Keep reading",
	} {
		if strings.Contains(result, unexpected) {
			t.Errorf("expected %q to be excluded, got: %q", unexpected, result)
		}
	}
}

func TestUnseenJapanScraper_ScrapeArticle_Lists(t *testing.T) {
	server := newUnseenJapanTestServer(`
<ul class="wp-block-list"><li>First <a href="https://example.com/">link</a></li><li>Second</li></ul>
<ol class="wp-block-list"><li>One</li><li>Two</li></ol>`)
	defer server.Close()

	scraper := &UnseenJapanScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, expected := range []string{
		"* First [link](https://example.com/)\n* Second\n\n",
		"1. One\n2. Two\n\n",
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("expected %q in result, got: %q", expected, result)
		}
	}
}

func TestUnseenJapanScraper_ScrapeArticle_Blockquote(t *testing.T) {
	server := newUnseenJapanTestServer(`<blockquote class="wp-block-quote"><p>A quote.</p></blockquote>`)
	defer server.Close()

	scraper := &UnseenJapanScraper{}
	result, err := scraper.ScrapeArticle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "> A quote.\n\n") {
		t.Errorf("expected blockquote in result, got: %q", result)
	}
}

func TestUnseenJapanScraper_ScrapeArticle_InvalidURL(t *testing.T) {
	scraper := &UnseenJapanScraper{}
	_, err := scraper.ScrapeArticle("http://invalid.invalid.invalid")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestUnseenJapanScraper_ScrapeTitle(t *testing.T) {
	server := newUnseenJapanTestServer(`<p class="wp-block-paragraph">Body.</p>`)
	defer server.Close()

	scraper := &UnseenJapanScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Suzuki Masa: The Samurai’s Daughter"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestUnseenJapanScraper_ScrapeTitle_InvalidURL(t *testing.T) {
	scraper := &UnseenJapanScraper{}
	_, err := scraper.ScrapeTitle("http://invalid.invalid.invalid")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestUnseenJapanScraper_ScrapeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trailing slash",
			input:    "https://unseen-japan.com/suzuki-masa-modern-japanese-nursing/",
			expected: "suzuki-masa-modern-japanese-nursing",
		},
		{
			name:     "no trailing slash",
			input:    "https://unseen-japan.com/suzuki-masa-modern-japanese-nursing",
			expected: "suzuki-masa-modern-japanese-nursing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper := &UnseenJapanScraper{}
			result, err := scraper.ScrapeFilename(tt.input)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestCreateArticleScraper_UnseenJapan(t *testing.T) {
	scraper, err := CreateArticleScraper("unseenjapan")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := scraper.(*UnseenJapanScraper); !ok {
		t.Errorf("expected *UnseenJapanScraper, got %T", scraper)
	}
}
