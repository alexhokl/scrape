package scraper

import "testing"

var _ ArticleScraper = (*KleppmannScraper)(nil)
var _ LinkScraper = (*KleppmannScraper)(nil)

func TestKleppmannScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &KleppmannScraper{},
		articleHTML: `<!DOCTYPE html><html><body class="wordpress"><div id="page"><div class="wrapper"><div id="content">
<h1>Announcing the Centre for Cryptographic Trust Infrastructure (CCTI)</h1>
<p style="font-size: 80%">Published by Martin Kleppmann on 07 Oct 2026.</p>
<p>We are establishing the <strong>Centre for Cryptographic Trust Infrastructure (CCTI)</strong>.</p>
<h2 id="introducing-ccti">Introducing CCTI</h2>
<p>Stay tuned!</p>
<div class="mailing-list-signup"><p>If you found this post useful, please support me on Patreon.</p></div>
</div></div></div></body></html>`,
		title: "Announcing the Centre for Cryptographic Trust Infrastructure (CCTI)",
		contains: []string{
			"# Announcing the Centre for Cryptographic Trust Infrastructure (CCTI)\n\nWe are establishing the **Centre for Cryptographic Trust Infrastructure (CCTI)**.\n\n",
			"## Introducing CCTI\n\nStay tuned!\n",
		},
		excludes: []string{"Published by", "Patreon"},
		indexHTML: `<html><body><div id="content"><ul>
<li><a href="/2026/10/07/ccti.html">Announcing the Centre for Cryptographic Trust Infrastructure (CCTI)</a></li>
<li><a href="/2025/12/08/ai-formal-verification.html">Prediction: AI will make formal verification go mainstream</a></li>
</ul></div><div class="widget"><ul><li class="page_item"><a href="/contact.html">About/Contact</a></li></ul></div></body></html>`,
		links: map[string]string{
			"Announcing the Centre for Cryptographic Trust Infrastructure (CCTI)": "/2026/10/07/ccti.html",
			"Prediction: AI will make formal verification go mainstream":          "/2025/12/08/ai-formal-verification.html",
		},
	})
}
