package scraper

import "testing"

var _ ArticleScraper = (*ProtonScraper)(nil)
var _ LinkScraper = (*ProtonScraper)(nil)

func TestProtonScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &ProtonScraper{},
		articleHTML: `<!DOCTYPE html><html><body class="bg-white"><main><article>
<div class="px-4"><div class="flow-root max-w-8xl mx-auto">
<header class="space-y-6 mb-6"><h1 class="font-serif text-4xl">DeGoogle your files with Easy Switch</h1></header>
<div class="lg:gap-14 lg:grid"><div>
<figure class="-mx-4 mb-6"><img src="https://res.cloudinary.com/hero.png" alt="Hero"></figure>
<footer class="space-y-6 mb-6">By Proton Team</footer>
<div class="post-content">
<p class="wp-block-paragraph"><strong>Don’t see Easy Switch yet?</strong></p>
<p class="wp-block-paragraph">Moved with <a href="/easyswitch">Easy Switch</a> and <a href="https://drive.proton.me/" target="_blank">Proton Drive<span class="sr-only">(new window)</span></a>.</p>
<div class="text-center"><a class="btn" href="https://drive.proton.me/">Use Easy Switch</a></div>
<p class="has-text-align-center wp-block-paragraph">Not on Drive yet? <a href="https://account.proton.me/">Create a free account</a>.</p>
<h2 class="wp-block-heading">Start an Easy Switch import</h2>
<ol class="wp-block-list"><li>Select <strong>Easy Switch</strong>.</li><li>Start the import.</li></ol>
<div class="wp-block-image"><figure class="aligncenter size-full"><img src="https://res.cloudinary.com/step.png" alt="Step"></figure></div>
</div>
<footer class="mb-20 mt-6">Share this article</footer>
</div></div></div></div></article></main></body></html>`,
		title: "DeGoogle your files with Easy Switch",
		contains: []string{
			"**Don’t see Easy Switch yet?**\n\n",
			"[Proton Drive](https://drive.proton.me/).\n\n",
			"## Start an Easy Switch import\n\n",
			"1. Select **Easy Switch**.\n2. Start the import.\n\n",
			"![Step](https://res.cloudinary.com/step.png)",
		},
		excludes: []string{"(new window)", "Use Easy Switch", "Not on Drive yet?", "hero.png", "By Proton Team", "Share this article"},
		indexHTML: `<html><body><div class="container-grid">
<div class="blog-item group"><div><div class="blog-item-title mt-3"><a class="no-underline" href="/blog/gmail-ads">Google can now use your data for ads</a></div></div></div>
<div class="blog-item group"><div><div class="blog-item-title mt-3"><a class="no-underline" href="/blog/is-google-drive-secure">Is Google Drive secure?</a></div></div></div>
</div><a href="/blog/all">All articles</a></body></html>`,
		links: map[string]string{
			"Google can now use your data for ads": "/blog/gmail-ads",
			"Is Google Drive secure?":              "/blog/is-google-drive-secure",
		},
	})
}
