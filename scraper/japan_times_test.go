package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newJapanTimesTestServer(bodyClass string, body string) *httptest.Server {
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<div class="jt-top-article-details paddingDefault" data-widget-id="582">
		<div class="row article-print-logo">
			<img src="/theme_japantimes/images/logo-with-slogan.svg" alt="The Japan Times" />
		</div>
		<div class="row">
			<div class="col-md-8">
				<div class="article-section">
					<a href="/community/" title="COMMUNITY">COMMUNITY</a>
				</div>
				<h1 class="title-article">As JET Programme turns 40, alumni and ministers stress its importance in age of AI</h1>
				<h3 class="subtitle-article">Four decades of exchange</h3>
			</div>
		</div>
	</div>
	<div class="jt-article-details">
		<div class="article-image-gallery">
			<article>
				<a title="Peter Cave joined the JET Programme&amp;#039;s inaugural cohort in 1987." class="fancybox" href="https://example.com/uploads/573331.jpeg"
					data-fancybox="gallery" data-caption="Peter Cave joined the JET Programme&#039;s inaugural cohort in 1987.  | COURTESY OF PETER CAVE ">
					<div class="layout-ratio">
						<picture>
							<img src="https://cdn.example.com/?url=https://example.com/uploads/573331.jpeg&q=100" alt="Peter Cave joined the JET Programme&amp;#039;s inaugural cohort in 1987.">
						</picture>
					</div>
				</a>
				<div class="imgCaption"><span class="imgCaption-text">Peter Cave joined the JET Programme's inaugural cohort in 1987.</span></div>
			</article>
			<article>
				<a class="fancybox" href="https://example.com/uploads/second.jpeg" data-caption="Second image"></a>
			</article>
		</div>
		<div class="article-meta">
			<div class="publish-date">Nov 11, 2026</div>
		</div>
		<div id="jtarticle">
			<div class="` + bodyClass + `">
` + body + `
			</div>
		</div>
		<div class="article-info">
			<div class="share-bottom no-print">Share</div>
		</div>
		<div class="subscribe no-print">
			<p>In a time of both misinformation and too much information, quality journalism is more crucial than ever.</p>
		</div>
	</div>
	<div class="jt-related-articles no-print">
		<div class="article"><div class="article-title"><a href="/news/related/">Related story</a></div></div>
	</div>
	<div class="jt-sidebar-articles no-print">
		<div class="article-byline">By <a href="/author/1/someone/">Sidebar Author</a></div>
	</div>
</body>
</html>`

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(html))
	}))
}

func scrapeJapanTimesTestArticle(t *testing.T, bodyClass string, body string) string {
	t.Helper()

	server := newJapanTimesTestServer(bodyClass, body)
	defer server.Close()

	scraper := &JapanTimesScraper{}
	result, err := scraper.ScrapeArticle(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return result
}

func TestJapanTimesScraper_ScrapeArticle_TitleSubtitleAndImage(t *testing.T) {
	result := scrapeJapanTimesTestArticle(t, "article-body", `<p>Body text.</p>`)

	expected := "# As JET Programme turns 40, alumni and ministers stress its importance in age of AI\n\n" +
		"*Four decades of exchange*\n\n" +
		"![](https://example.com/uploads/573331.jpeg)\n\n" +
		"*Peter Cave joined the JET Programme's inaugural cohort in 1987. | COURTESY OF PETER CAVE*\n\n" +
		"Body text.\n\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestJapanTimesScraper_ScrapeArticle_ExcludesNonBodyContent(t *testing.T) {
	result := scrapeJapanTimesTestArticle(t, "article-body", `<p>Body text.</p>`)

	for _, unexpected := range []string{
		"COMMUNITY",
		"second.jpeg",
		"Nov 11, 2026",
		"Share",
		"misinformation",
		"Related story",
		"Sidebar Author",
		"subscribers only",
	} {
		if strings.Contains(result, unexpected) {
			t.Errorf("expected %q to be excluded, got: %q", unexpected, result)
		}
	}
}

func TestJapanTimesScraper_ScrapeArticle_ExcludesAdSlots(t *testing.T) {
	result := scrapeJapanTimesTestArticle(t, "article-body", `
<p>Before.</p><p>Middle.<div class='jt_content_ad DisplayAd' id='div-gpt-ad-1'><script>googletag.cmd.push(function() { googletag.display('div-gpt-ad-1'); });</script></div></p><p>After.</p>`)

	if strings.Contains(result, "googletag") {
		t.Errorf("expected ad slot to be excluded, got: %q", result)
	}

	expected := "Before.\n\nMiddle.\n\nAfter.\n\n"
	if !strings.HasSuffix(result, expected) {
		t.Errorf("expected %q at end of result, got: %q", expected, result)
	}
}

func TestJapanTimesScraper_ScrapeArticle_Paywall(t *testing.T) {
	result := scrapeJapanTimesTestArticle(t, "article-body blurred-text", `
<p>Preview paragraph.</p>
<div class="blurred-text-gradient"></div>`)

	expected := "Preview paragraph.\n\n*The rest of this article is available to subscribers only.*\n\n"
	if !strings.HasSuffix(result, expected) {
		t.Errorf("expected %q at end of result, got: %q", expected, result)
	}
}

func TestJapanTimesScraper_ScrapeArticle_Headings(t *testing.T) {
	result := scrapeJapanTimesTestArticle(t, "article-body", `
<h2>Section One</h2>
<h3>Subsection</h3>
<h4>Detail</h4>`)

	for _, expected := range []string{
		"## Section One\n\n",
		"### Subsection\n\n",
		"#### Detail\n\n",
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("expected %q in result, got: %q", expected, result)
		}
	}
}

func TestJapanTimesScraper_ScrapeArticle_InlineFormatting(t *testing.T) {
	result := scrapeJapanTimesTestArticle(t, "article-body", `
<p><span class="dateline">Tokyo – </span>The <a href="https://example.com/jet">JET Programme</a> is an <em>irreplaceable</em> and <strong>national</strong> asset.</p>`)

	expected := "Tokyo – The [JET Programme](https://example.com/jet) is an *irreplaceable* and **national** asset.\n\n"
	if !strings.Contains(result, expected) {
		t.Errorf("expected %q in result, got: %q", expected, result)
	}
}

func TestJapanTimesScraper_ScrapeArticle_Lists(t *testing.T) {
	result := scrapeJapanTimesTestArticle(t, "article-body", `
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

func TestJapanTimesScraper_ScrapeArticle_Blockquote(t *testing.T) {
	result := scrapeJapanTimesTestArticle(t, "article-body", `<blockquote><p>A quoted remark.</p></blockquote>`)

	if !strings.Contains(result, "> A quoted remark.\n\n") {
		t.Errorf("expected blockquote in result, got: %q", result)
	}
}

func TestJapanTimesScraper_ScrapeArticle_InvalidURL(t *testing.T) {
	scraper := &JapanTimesScraper{}
	_, err := scraper.ScrapeArticle("http://127.0.0.1:0/invalid")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestJapanTimesScraper_ScrapeTitle(t *testing.T) {
	server := newJapanTimesTestServer("article-body", ``)
	defer server.Close()

	scraper := &JapanTimesScraper{}
	result, err := scraper.ScrapeTitle(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "As JET Programme turns 40, alumni and ministers stress its importance in age of AI"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestJapanTimesScraper_ScrapeTitle_InvalidURL(t *testing.T) {
	scraper := &JapanTimesScraper{}
	_, err := scraper.ScrapeTitle("http://127.0.0.1:0/invalid")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestJapanTimesScraper_ScrapeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trailing slash",
			input:    "https://www.japantimes.co.jp/community/2026/11/11/jet-programme-40-years/",
			expected: "jet-programme-40-years",
		},
		{
			name:     "no trailing slash",
			input:    "https://www.japantimes.co.jp/business/2026/11/14/japan-philippines-oil-reserve-deal",
			expected: "japan-philippines-oil-reserve-deal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper := &JapanTimesScraper{}
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
