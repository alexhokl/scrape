package scraper

import "testing"

var _ ArticleScraper = (*CharmScraper)(nil)
var _ LinkScraper = (*CharmScraper)(nil)

func TestCharmScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &CharmScraper{},
		// the site leaves <p> and <li> unclosed
		articleHTML: `<!DOCTYPE html><html><body class="blog-post"><div id="root"><div class="content blog"><article>
<div class="title"><h1>Wishlist Endpoint Discovery</h1></div>
<div class="post">
<p>We just added endpoint discovery to <a href="https://github.com/charmbracelet/wishlist">Wishlist</a>.</p><a class="cart wishlist" href="/wishlist"><div class="info"><h3 class="puffy">Wishlist</h3><p>Your SSH directory.</p></div></a>
<p>It becomes a <em>TUI</em> for your <code>~/.ssh/config</code>.
<h2>Endpoint Discovery</h2>
<ul>
<li>DNS SRV records
<li>Tailscale
</ul>
<pre tabindex=0 style="color:#c4c4c4"><code><span style=display:flex;><span>wishlist --tailscale.net=charm</span></span></code></pre>
</div></article></div></div></body></html>`,
		title: "Wishlist Endpoint Discovery",
		contains: []string{
			"We just added endpoint discovery to [Wishlist](https://github.com/charmbracelet/wishlist).\n\n",
			"It becomes a *TUI* for your `~/.ssh/config`.\n\n",
			"## Endpoint Discovery\n\n",
			"* DNS SRV records\n* Tailscale\n\n",
			"```\nwishlist --tailscale.net=charm\n```\n",
		},
		excludes: []string{"Your SSH directory."},
		indexHTML: `<html><body><div class="content blog">
<article><p class="date">Nov 1, 2026</p><h1><a href="/blog/crush-and-docker-mcp/">Crush + Docker MCP</a></h1></article>
<article><p class="date">Oct 1, 2026</p><h1><a href="/blog/v2/">v2</a></h1></article>
</div><ul class="social"><li class="github"><a href="https://github.com/charmbracelet">GitHub</a></li></ul></body></html>`,
		links: map[string]string{
			"Crush + Docker MCP": "/blog/crush-and-docker-mcp/",
			"v2":                 "/blog/v2/",
		},
	})
}
