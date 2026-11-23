package scraper

import "testing"

var _ ArticleScraper = (*JuliaEvansScraper)(nil)
var _ LinkScraper = (*JuliaEvansScraper)(nil)

func TestJuliaEvansScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &JuliaEvansScraper{},
		articleHTML: `<!DOCTYPE html><html><body><div id="wrap">
<header><hgroup><h1><a href="/">Julia Evans</a></h1></hgroup></header>
<div id="main"><div id="content"><div><article class="hentry" role="article">
<header><h1 class="entry-title">Replacing the old battery on rechargeable bike lights</h1>
<p class="meta sans"><time class="date">September 27, 2026</time></p></header>
<main>
<p>Hello! Recently I needed bike lights.</p>
<h3 id="step-1" class="post-heading">step 1: cut it open</h3>
<img src="https://jvns.ca/images/pcb.png">
<ol><li>Use a desoldering pump</li><li>Pull them apart</li></ol>
</main>
<footer>Want a weekly digest of this blog?</footer>
</article></div></div></div></div></body></html>`,
		title: "Replacing the old battery on rechargeable bike lights",
		contains: []string{
			"Hello! Recently I needed bike lights.\n\n",
			"### step 1: cut it open\n\n",
			"![](https://jvns.ca/images/pcb.png)\n\n",
			"1. Use a desoldering pump\n2. Pull them apart\n",
		},
		excludes: []string{"Julia Evans\n", "September 27, 2026", "weekly digest"},
		indexHTML: `<html><body><div><article class="hentry">
<h3><a href="#uncategorized">Uncategorized</a></h3>
<p class="article-list with-favourites"><a href="/blog/2026/09/27/bike-lights/">Replacing the old battery</a></p>
<p class="article-list"><a href="/blog/2026/09/27/bike-lights/">Replacing the old battery</a></p>
<p class="article-list"><a href="/blog/2026/07/21/django/">More nice Django things</a></p>
</article></div></body></html>`,
		links: map[string]string{
			"Replacing the old battery": "/blog/2026/09/27/bike-lights/",
			"More nice Django things":   "/blog/2026/07/21/django/",
		},
	})
}
