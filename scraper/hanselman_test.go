package scraper

import "testing"

var _ ArticleScraper = (*HanselmanScraper)(nil)
var _ LinkScraper = (*HanselmanScraper)(nil)

func TestHanselmanScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &HanselmanScraper{},
		articleHTML: `<!DOCTYPE html><html><body class="line-darkbrown">
<section id="topbar" class="containerOuter"><section class="containerInner"><h1>Scott Hanselman</h1></section></section>
<main><article>
<section class="containerOuter blogHeader"><section class="containerInnerNarrow"><h2 class="blogTitle"><a href="https://www.hanselman.com/blog/sample-post">Migrating a Synology NAS to a UniFi UNAS Pro 8</a></h2></section></section>
<section class="containerOuter blogBodyContainer"><section class="containerInner blogBody">
<div class="ad-item ad-sponsoreditem"> Sponsored By<br/><div id="lqdmad"></div></div>
<p>I’ve had a <a href="/blog/synology">Synology NAS</a> for a very long time.</p>
<h4>Robocopy</h4>
<pre><code>robocopy \\nas\share D:\share /MIR /MT:32</code></pre>
<div class="about-scott"><div class="bioBox"><div class="bioBoxInner">Scott Hanselman is a former professor. <a href="http://hanselman.com/about">About</a></div></div></div>
<div class="ads-box"><div class="ad-item"><strong>Hosting By</strong></div></div>
<div class="item-ads"><div class="lqm_ad"></div></div>
<div class="item-footer"><div id="comments-start"></div><a class="dbc-comment-on-post-link" href="https://www.hanselman.com/blog/sample-post/comments">Comment on this post [4]</a>
<div class="item-footer-tools"> Share on <a href="https://bsky.app/intent/compose">Bluesky</a></div></div>
</section></section></article></main></body></html>`,
		title: "Migrating a Synology NAS to a UniFi UNAS Pro 8",
		contains: []string{
			"# Migrating a Synology NAS to a UniFi UNAS Pro 8\n\nI’ve had a [Synology NAS](http",
			"/blog/synology) for a very long time.\n\n#### Robocopy\n\n",
			"```\nrobocopy \\\\nas\\share D:\\share /MIR /MT:32\n```\n",
		},
		excludes: []string{"Scott Hanselman", "Sponsored By", "Hosting By", "Comment on this post", "Share on"},
		indexHTML: `<html><body><main>
<article><section class="containerOuter blogHeader"><section class="containerInnerNarrow"><h2 class="blogTitle"><a href="/blog/debugging-my-new-network">Debugging my new network</a></h2></section></section>
<section class="blogBody"><p><a href="/blog/synology">Synology NAS</a></p></section></article>
<article><section class="containerOuter blogHeader"><section class="containerInnerNarrow"><h2 class="blogTitle"><a href="/blog/migrating-a-synology-nas">Migrating a Synology NAS</a></h2></section></section></article>
</main></body></html>`,
		links: map[string]string{
			"Debugging my new network": "/blog/debugging-my-new-network",
			"Migrating a Synology NAS": "/blog/migrating-a-synology-nas",
		},
	})
}
