package scraper

import "testing"

var _ ArticleScraper = (*RakyllScraper)(nil)
var _ LinkScraper = (*RakyllScraper)(nil)

func TestRakyllScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &RakyllScraper{},
		articleHTML: `<!DOCTYPE html><html><body>
<p><a href="https://rakyll.org/">Home</a> <a href="https://rakyll.org/archive/">Articles</a></p>
<article>
<h1>Intellectual click</h1>
<p class="meta">November 21, 2026</p>
<p>There is a kind of connection I don’t have a good name for.</p>
<p>I have felt this with <em>only two</em> people in my life.</p>
</article></body></html>`,
		title: "Intellectual click",
		contains: []string{
			"# Intellectual click\n\nThere is a kind of connection I don’t have a good name for.\n\n",
			"I have felt this with *only two* people in my life.\n",
		},
		excludes: []string{"November 21, 2026", "Home", "Articles"},
		indexHTML: `<html><body><p><a href="/">Home</a></p><ul class="posts">
<li><span>Nov 21, 2026</span> <a href="/intellectual-click/">Intellectual click</a></li>
<li><span>Oct 1, 2026</span> <a href="/chiptune/">Catchy 8-bit music</a></li>
</ul></body></html>`,
		links: map[string]string{
			"Intellectual click": "/intellectual-click/",
			"Catchy 8-bit music": "/chiptune/",
		},
	})
}
