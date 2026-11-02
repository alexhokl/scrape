package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newProPublicaTestServer(body string) *httptest.Server {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<article class="wp-block-group p-grid-text-container">
		<figure>
			<img class="p-opener__art" alt="" src="https://example.com/hero.jpg?w=1149">
			<figcaption class="p-attribution screen-reader-text">
				<span class="p-attribution__credit">Hero Credit</span>
			</figcaption>
		</figure>
		<div class="wp-block-group p-opener__topic-title-dek-wrapper p-hide-below-lg">
			<p class="p-opener__series wp-block-propublica-primary-term">
				<a href="https://example.com/series/coronavirus" class="p-opener__series-link">Coronavirus</a>
			</p>
			<h1 class="p-opener__hed p-opener__hed--small wp-block-post-title">Leaked Documents Show How China’s Army of Paid Internet Trolls Helped Censor the Coronavirus</h1>
		</div>
		<div class="wp-block-group p-opener__topic-title-dek-wrapper">
			<h1 class="p-opener__hed p-hide-above-lg wp-block-post-title">Leaked Documents Show How China’s Army of Paid Internet Trolls Helped Censor the Coronavirus</h1>
			<h2 class="p-opener__dek wp-block-propublica-dek">
	As the coronavirus spread in China, the government stage-managed what appeared on the domestic internet.</h2>
		</div>
		<div class="wp-block-group p-article-meta-1">
			<div class="p-article-meta-1__byline wp-block-propublica-byline">by Byline Author</div>
			<div class="p-article-meta-1__pubdate wp-block-post-date">Dec. 19, 2020</div>
		</div>
		<div class="wp-block-group article-body">
			<aside class="wp-block-group">
				<div class="wp-block-propublica-notes--top wp-block-propublica-notes">
					<div class="wp-block-propublica-note">
						<p class="wp-block-paragraph">ProPublica is a nonprofit newsroom. Sign up to receive <a href="https://example.com/newsletters">our biggest stories</a>.</p>
					</div>
					<div class="wp-block-propublica-note--co-publish wp-block-propublica-note">
						<p class="wp-block-paragraph">This article is co-published with The New York Times.</p>
					</div>
				</div>
			</aside>
			<div class="entry-content wp-block-post-content">
` + body + `
			</div>
			<aside class="wp-block-query">
				<h2 class="wp-block-heading screen-reader-text">Corrections</h2>
			</aside>
			<aside class="wp-block-group">
				<h3 class="wp-block-post-title">Related Author Bio</h3>
			</aside>
		</div>
	</article>
</body>
</html>`

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(html))
	}))
}

func scrapeProPublicaTestArticle(t *testing.T, body string) string {
	t.Helper()

	server := newProPublicaTestServer(body)
	defer server.Close()

	scraper := &ProPublicaScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return result
}

func TestProPublicaScraper_ScrapeArticle_TitleAndDek(t *testing.T) {
	result := scrapeProPublicaTestArticle(t, ``)

	expected := "# Leaked Documents Show How China’s Army of Paid Internet Trolls Helped Censor the Coronavirus\n\n" +
		"*As the coronavirus spread in China, the government stage-managed what appeared on the domestic internet.*\n\n"
	if !strings.HasPrefix(result, expected) {
		t.Errorf("expected title and dek at start of result, got: %q", result)
	}

	if count := strings.Count(result, "# Leaked Documents"); count != 1 {
		t.Errorf("expected title once, got %d times: %q", count, result)
	}
}

func TestProPublicaScraper_ScrapeArticle_ExcludesNonBodyContent(t *testing.T) {
	result := scrapeProPublicaTestArticle(t, `<p class="wp-block-paragraph">Body text.</p>`)

	for _, unexpected := range []string{
		"hero.jpg",
		"Hero Credit",
		"Coronavirus](",
		"Byline Author",
		"Dec. 19, 2020",
		"nonprofit newsroom",
		"co-published",
		"Corrections",
		"Related Author Bio",
	} {
		if strings.Contains(result, unexpected) {
			t.Errorf("expected %q to be excluded, got: %q", unexpected, result)
		}
	}

	if !strings.Contains(result, "Body text.\n\n") {
		t.Errorf("expected body text, got: %q", result)
	}
}

func TestProPublicaScraper_ScrapeArticle_Headings(t *testing.T) {
	result := scrapeProPublicaTestArticle(t, `
<h2 class="wp-block-heading">Section One</h2>
<h3 class="wp-block-heading">A Great Sea of Placidity</h3>
<h4 class="wp-block-heading">Detail</h4>`)

	for _, expected := range []string{
		"## Section One\n\n",
		"### A Great Sea of Placidity\n\n",
		"#### Detail\n\n",
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("expected %q in result, got: %q", expected, result)
		}
	}
}

func TestProPublicaScraper_ScrapeArticle_InlineFormatting(t *testing.T) {
	result := scrapeProPublicaTestArticle(t, `
<p class="wp-block-paragraph">Li Wenliang had <a href="https://example.com/died">died of COVID-19</a>, an <em>unprecedented</em> and <strong>serious</strong> challenge.</p>`)

	expected := "Li Wenliang had [died of COVID-19](https://example.com/died), an *unprecedented* and **serious** challenge.\n\n"
	if !strings.Contains(result, expected) {
		t.Errorf("expected %q in result, got: %q", expected, result)
	}
}

func TestProPublicaScraper_ScrapeArticle_EmptyParagraphSkipped(t *testing.T) {
	result := scrapeProPublicaTestArticle(t, `<p class="wp-block-paragraph">  </p><p class="wp-block-paragraph">Text.</p>`)

	if strings.Contains(result, "\n\n\n") {
		t.Errorf("expected empty paragraph to be skipped, got: %q", result)
	}
}

func TestProPublicaScraper_ScrapeArticle_Image(t *testing.T) {
	result := scrapeProPublicaTestArticle(t, `
<figure class="wp-block-image size-propublica-position-medium"><img src="https://example.com/excerpt.jpg?w=752" alt="Excerpt" srcset="https://example.com/excerpt.jpg?w=1865 1865w"><figcaption class="wp-element-caption">A directive excerpt</figcaption></figure>
<figure class="wp-block-image"><img src="https://example.com/no-caption.jpg" alt=""></figure>`)

	for _, expected := range []string{
		"![Excerpt](https://example.com/excerpt.jpg?w=752)\n\n*A directive excerpt*\n\n",
		"![](https://example.com/no-caption.jpg)\n\n",
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("expected %q in result, got: %q", expected, result)
		}
	}
}

func TestProPublicaScraper_ScrapeArticle_LeadIn(t *testing.T) {
	result := scrapeProPublicaTestArticle(t, `
<p class="wp-block-paragraph">Before.</p>
<div class="wp-block-propublica-lead-in">
<h3 class="wp-block-heading">Special instructions were issued.</h3>
<figure class="wp-block-image"><img src="https://example.com/lead-in.jpg" alt=""></figure>
</div>
<p class="wp-block-paragraph">After.</p>`)

	expected := "Before.\n\n### Special instructions were issued.\n\n![](https://example.com/lead-in.jpg)\n\nAfter.\n\n"
	if !strings.HasSuffix(result, expected) {
		t.Errorf("expected %q at end of result, got: %q", expected, result)
	}
}

func TestProPublicaScraper_ScrapeArticle_ExcludesAdSlots(t *testing.T) {
	result := scrapeProPublicaTestArticle(t, `
<p class="wp-block-paragraph">Before.</p>
<div class="p-hide-print p-bb--size-full wp-block-propublica-ad-slot"><div><p>Advertisement</p></div></div>
<p class="wp-block-paragraph">After.</p>`)

	if strings.Contains(result, "Advertisement") {
		t.Errorf("expected ad slot to be excluded, got: %q", result)
	}
}

func TestProPublicaScraper_ScrapeArticle_Lists(t *testing.T) {
	result := scrapeProPublicaTestArticle(t, `
<ul><li>First <a href="https://example.com/a">link</a></li><li>Second</li></ul>
<ol><li>One</li><li>Two</li></ol>`)

	for _, expected := range []string{
		"* First [link](https://example.com/a)\n* Second\n\n",
		"1. One\n2. Two\n\n",
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("expected %q in result, got: %q", expected, result)
		}
	}
}

func TestProPublicaScraper_ScrapeArticle_Blockquote(t *testing.T) {
	result := scrapeProPublicaTestArticle(t, `<blockquote class="wp-block-quote"><p>Be sleek and silent.</p></blockquote>`)

	if !strings.Contains(result, "> Be sleek and silent.\n\n") {
		t.Errorf("expected blockquote in result, got: %q", result)
	}
}

func TestProPublicaScraper_ScrapeArticle_InvalidURL(t *testing.T) {
	scraper := &ProPublicaScraper{}
	_, err := scraper.ScrapeArticle("http://127.0.0.1:0/invalid")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestProPublicaScraper_ScrapeTitle(t *testing.T) {
	server := newProPublicaTestServer(``)
	defer server.Close()

	scraper := &ProPublicaScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Leaked Documents Show How China’s Army of Paid Internet Trolls Helped Censor the Coronavirus"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestProPublicaScraper_ScrapeTitle_InvalidURL(t *testing.T) {
	scraper := &ProPublicaScraper{}
	_, err := scraper.ScrapeTitle("http://127.0.0.1:0/invalid")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestProPublicaScraper_ScrapeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no trailing slash",
			input:    "https://www.propublica.org/article/leaked-documents-show-how-chinas-army-of-paid-internet-trolls-helped-censor-the-coronavirus",
			expected: "leaked-documents-show-how-chinas-army-of-paid-internet-trolls-helped-censor-the-coronavirus",
		},
		{
			name:     "trailing slash",
			input:    "https://www.propublica.org/article/some-story/",
			expected: "some-story",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper := &ProPublicaScraper{}
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
