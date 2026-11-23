package scraper

import "testing"

var _ ArticleScraper = (*ObservationDeckScraper)(nil)
var _ LinkScraper = (*ObservationDeckScraper)(nil)

func TestObservationDeckScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &ObservationDeckScraper{},
		articleHTML: `<!DOCTYPE html><html><body class="notransition"><div class="wrapper post"><main class="page-content"><article>
<header class="header"><h1 class="header-title">Fool's Expertise</h1><div class="post-meta">September 27, 2026</div></header>
<div class="page-content">
<div class="paragraph"><p>I listened to the recent <a href="https://www.youtube.com/watch?v=x">interview</a>.</p></div>
<div class="quoteblock"><blockquote><div class="paragraph"><p>A quote.</p></div></blockquote></div>
<div class="paragraph"><p>Hinton was wrong because his prediction did not reflect <strong>what a radiologist actually does</strong>.</p></div>
</div></article></main></div></body></html>`,
		title: "Fool's Expertise",
		contains: []string{
			"# Fool's Expertise\n\nI listened to the recent [interview](https://www.youtube.com/watch?v=x).\n\n",
			"> A quote.\n\n",
			"**what a radiologist actually does**.\n",
		},
		excludes: []string{"September 27, 2026"},
		indexHTML: `<html><body class="notransition"><main>
<article class="post-item"><h4 class="post-item-title"><a href="/2026/09/27/fools-expertise/">Fool's Expertise</a></h4></article>
<article class="post-item"><h4 class="post-item-title"><a href="/2026/09/20/what-sun-got-wrong/">What Sun got wrong</a></h4></article>
</main><footer class="footer"><div class="footer_social-icons"><a href="/index.xml">RSS</a></div></footer></body></html>`,
		links: map[string]string{
			"Fool's Expertise":   "/2026/09/27/fools-expertise/",
			"What Sun got wrong": "/2026/09/20/what-sun-got-wrong/",
		},
	})
}
