package scraper

import "testing"

var _ ArticleScraper = (*JessFrazelleScraper)(nil)
var _ LinkScraper = (*JessFrazelleScraper)(nil)

func TestJessFrazelleScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &JessFrazelleScraper{},
		articleHTML: `<!DOCTYPE html><html class="no-js"><body><div class="container single"><article>
<h1>Home Hardware, Part 2: Width Matters</h1>
<div class="time">Wednesday, September 9, 2026 · 9 min read</div>
<aside class="toc"><h4>Table of Contents</h4><nav id="TableOfContents"><ul><li><a href="#special">Special</a></li></ul></nav></aside>
<p>We’re back!<sup class="footnote-ref" id="fnref:calc"><a href="#fn:calc">1</a></sup></p>
<h2 id="special">Special, apparently</h2>
<table><thead><tr><th>JLC stackup</th><th>Width</th></tr></thead><tbody><tr><td>2116</td><td>0.1778 mm</td></tr></tbody></table>
<div class="footnotes"><hr><ol><li id="fn:calc">More on the calculator.
<a class="footnote-return" href="#fnref:calc"><sup>[return]</sup></a></li></ol></div>
</article></div></body></html>`,
		title: "Home Hardware, Part 2: Width Matters",
		contains: []string{
			"# Home Hardware, Part 2: Width Matters\n\nWe’re back!1\n\n",
			"## Special, apparently\n\n",
			"| JLC stackup | Width |\n| --- | --- |\n| 2116 | 0.1778 mm |\n",
			"1. More on the calculator.\n",
		},
		excludes: []string{"9 min read", "Table of Contents", "[return]"},
		indexHTML: `<html><body><h1>Posts</h1>
<article><h3><a href="/post/home-hardware-part-two/">Home Hardware, Part 2: Width Matters</a></h3></article>
<article><h3><a href="/post/books/">Books I Recommend</a></h3></article>
<article class="publication"><h3><a href="https://queue.acm.org/detail.cfm?id=1">ACM Queue</a></h3></article>
</body></html>`,
		links: map[string]string{
			"Home Hardware, Part 2: Width Matters": "/post/home-hardware-part-two/",
			"Books I Recommend":                    "/post/books/",
		},
	})
}
