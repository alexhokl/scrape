package scraper

import "testing"

var _ ArticleScraper = (*FilippoScraper)(nil)
var _ LinkScraper = (*FilippoScraper)(nil)

func TestFilippoScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &FilippoScraper{},
		articleHTML: `<!DOCTYPE html><html><body>
<header><a href="https://filippo.io">Filippo</a><div class="subscribe"><p><a href="https://words.filippo.io/rss/">Feed</a></p></div></header>
<main><article><time datetime="2026-10-16">16 Oct 2026</time>
<h1>The NIST Elliptic Curves Seeds Bounty Has Been Claimed</h1>
<section><!-- buttondown-editor-mode: plaintext -->
<p>The output of <a href="https://man.archlinux.org/man/drand48.3.en">drand48</a> seeded with the time.</p>
<pre><code>seed = drand48()</code></pre>
<h2 id="the-picture">The picture</h2>
<p><img src="https://assets.buttondown.email/images/dice.jpeg" alt="Dice"></p>
</section></article></main></body></html>`,
		title: "The NIST Elliptic Curves Seeds Bounty Has Been Claimed",
		contains: []string{
			"The output of [drand48](https://man.archlinux.org/man/drand48.3.en) seeded with the time.\n\n",
			"```\nseed = drand48()\n```\n\n",
			"## The picture\n\n",
			"![Dice](https://assets.buttondown.email/images/dice.jpeg)\n",
		},
		excludes: []string{"16 Oct 2026", "Feed", "buttondown-editor-mode"},
		indexHTML: `<html><body><header><div class="subscribe"><p><a href="/rss/">Feed</a></p></div></header><main>
<article><h1><a href="/seeds-solved/">The NIST Elliptic Curves Seeds Bounty Has Been Claimed</a></h1></article>
<article><h1><a href="/right-answers/">(Post-Quantum) Cryptographic Right Answers</a></h1></article>
</main></body></html>`,
		links: map[string]string{
			"The NIST Elliptic Curves Seeds Bounty Has Been Claimed": "/seeds-solved/",
			"(Post-Quantum) Cryptographic Right Answers":             "/right-answers/",
		},
	})
}
