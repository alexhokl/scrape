package scraper

import "testing"

var _ ArticleScraper = (*RickStrahlScraper)(nil)
var _ LinkScraper = (*RickStrahlScraper)(nil)

func TestRickStrahlScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &RickStrahlScraper{},
		articleHTML: `<!DOCTYPE html><html><body><div id="page-wrapper"><main class="post-content">
<article itemprop="blogPost"><header><h2 class="page-header mb-1" itemprop="headline_name"> Back to Basics: CSS @container Filtering </h2></header>
<div class="byline"><div class="share-box"><a href="https://twitter.com/intent/tweet">Share on</a></div><div class="leftbox">October 02, 2026 • 669 views</div></div>
<div class="postcontent" id="ArticleBody" itemprop="articleBody">
<div class="document-outline"><div class="document-outline-header">On this page:</div></div>
<p><img src="https://weblog.west-wind.com/imageContent/banner.jpg" alt="Container Banner"></p>
<p>I've been using <code>@media</code> queries forever.</p>
<pre><code class="language-css">@container (width > 920px) {
    .topic-outline { display: block; }
}</code></pre>
<p><a href="https://www.linqpad.net/?affiliate=x" title="LinqPad"><img src="https://weblog.west-wind.com/images/sponsors/LinqPad-DisplayAd.jpg" class="da-content-image" /></a></p>
<h2 id="resources">Resources</h2>
<ul><li><a href="https://developer.mozilla.org/container">@container MDN Documentation</a></li></ul>
<div style="margin-top: 30px;font-size: 0.8em;"><span>this post was created and published with the</span> <a href="https://markdownmonster.west-wind.com">Markdown Monster Editor</a></div>
<div style="margin: 30px 0"><h3>Other Posts you might also like</h3><ul><li><a href="/posts/2023/x">HTML Table Cell Overflow</a></li></ul></div>
</div></article></main></div></body></html>`,
		title: "Back to Basics: CSS @container Filtering",
		contains: []string{
			"![Container Banner](https://weblog.west-wind.com/imageContent/banner.jpg)\n\n",
			"I've been using `@media` queries forever.\n\n",
			"```css\n@container (width > 920px) {\n    .topic-outline { display: block; }\n}\n```\n\n",
			"## Resources\n\n* [@container MDN Documentation](https://developer.mozilla.org/container)\n",
		},
		excludes: []string{"Share on", "669 views", "On this page", "LinqPad", "Markdown Monster", "Other Posts"},
		indexHTML: `<html><body><div><div class="post-item-wrapper">
<h2 class="page-header mb-1"><a href="/posts/2026/Oct/02/Back-to-Basics-CSS-container-Filtering">Back to Basics: CSS @container Filtering</a></h2>
<div class="postcontent mt-0"><p class="margin-top margin-bottom"><a href="/posts/2026/Oct/02/Back-to-Basics-CSS-container-Filtering">Read more...</a></p></div>
</div><div class="post-item-wrapper">
<h2 class="page-header mb-1"><a href="/posts/2026/Sep/01/JavaScript-and-Timezones">Back to Basics: JavaScript and Timezones</a></h2>
</div></div><nav class="tag-cloud"><a class="tag-cloud-link" href="/posts/categories/ASP.NET">ASP.NET</a></nav></body></html>`,
		links: map[string]string{
			"Back to Basics: CSS @container Filtering": "/posts/2026/Oct/02/Back-to-Basics-CSS-container-Filtering",
			"Back to Basics: JavaScript and Timezones": "/posts/2026/Sep/01/JavaScript-and-Timezones",
		},
	})
}
