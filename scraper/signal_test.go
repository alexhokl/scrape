package scraper

import "testing"

var _ ArticleScraper = (*SignalScraper)(nil)
var _ LinkScraper = (*SignalScraper)(nil)

func TestSignalScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &SignalScraper{},
		articleHTML: `<!DOCTYPE html><html><body id="signal"><section id="post" class="section post"><article>
<div class="container is-centered"><div class="blog-post-header columns is-centered"><div class="column is-7">
<h1 class="has-text-centered">Signal is Serious About Protecting Your Privacy</h1>
<p class="body2">Meredith Whittaker on 20 Nov 2026</p>
</div></div></div>
<div class="container is-centered"><div class="blog-post-content columns is-centered"><div class="column is-7">
<p><img src="/blog/images/header.png" alt="Header"></p>
<p>Privacy is our mission.<sup id="fnref:1"><a href="#fn:1" class="footnote">1</a></sup></p>
<ol><li>Open <strong>Settings</strong></li></ol>
<hr>
<div class="footnotes"><ol><li id="fn:1"><p>A footnote.</p></li></ol></div>
<div class="social-sharing has-text-centered"><p><a class="btn btn-twitter" href="https://twitter.com/intent/tweet">Tweet</a></p>
<h3>Don't have Signal? <a href="/download/">Give it a try!</a></h3></div>
</div></div></div></article></section></body></html>`,
		title: "Signal is Serious About Protecting Your Privacy",
		contains: []string{
			"![Header](http",
			"/blog/images/header.png)\n\n",
			"Privacy is our mission.1\n\n",
			"1. Open **Settings**\n",
			"---\n\n1. A footnote.\n",
		},
		excludes: []string{"Meredith Whittaker", "Tweet", "Give it a try"},
		indexHTML: `<html><body><div class="columns"><div class="column is-8 is-offset-2">
<h3 class="has-text-centered"><a href="/blog/protecting-your-privacy/">Signal is Serious About Protecting Your Privacy</a></h3>
<p><a href="/blog/protecting-your-privacy/">Read more...</a></p>
</div></div>
<div class="columns"><div class="column is-8 is-offset-2">
<h3 class="has-text-centered"><a href="/blog/proudly-nonprofit/">Proudly nonprofit</a></h3>
</div></div>
<div class="column"><ul><li><a href="/donate/">Donate</a></li></ul></div></body></html>`,
		links: map[string]string{
			"Signal is Serious About Protecting Your Privacy": "/blog/protecting-your-privacy/",
			"Proudly nonprofit": "/blog/proudly-nonprofit/",
		},
	})
}
