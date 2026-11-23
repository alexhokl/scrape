package scraper

import "testing"

var _ ArticleScraper = (*MitchellHashimotoScraper)(nil)
var _ LinkScraper = (*MitchellHashimotoScraper)(nil)

func TestMitchellHashimotoScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &MitchellHashimotoScraper{},
		articleHTML: `<!DOCTYPE html><html><body><div class="max-w-5xl">
<aside><div><h1 class="nimbus uppercase tracking-wide">Mitchell Hashimoto</h1></div></aside>
<main class="flex-auto"><h1 class="mb-20 nimbus uppercase hidden md:block">Mitchell Hashimoto</h1>
<div class="flex-auto"><div class="flex flex-1 items-start"><div class="w-full md:w-9/12">
<h1 class="nimbus md:mb-11">Libghostty 1.0</h1>
<div class="-mt-[5px] prose dark:prose-dark">
<p><code>libghostty-vt</code> 1.0 is here.</p>
<hr>
<h2 id="abi-stability">ABI Stability</h2>
<ul><li>Stable C ABI</li></ul>
</div></div></div></div>
<footer><a href="mailto:m@mitchellh.com">mail</a></footer></main></div></body></html>`,
		title: "Libghostty 1.0",
		contains: []string{
			"# Libghostty 1.0\n\n`libghostty-vt` 1.0 is here.\n\n---\n\n## ABI Stability\n\n* Stable C ABI\n",
		},
		excludes: []string{"Mitchell Hashimoto", "mail"},
		indexHTML: `<html><body><main><div class="flex-auto"><div class="flex flex-col divide-y">
<a class="text-gray-900" href="/writing/libghostty-1-0"><div class="flex"><div class="flex-auto"><h1 class="nimbus">Libghostty 1.0</h1></div><div class="md:w-3/12">November 6, 2026</div></div></a>
<a class="text-gray-900" href="/writing/superlogical"><div class="flex"><div class="flex-auto"><h1 class="nimbus">Superlogical</h1></div><div class="md:w-3/12">July 29, 2026</div></div></a>
</div></div><footer><a href="https://twitter.com/mitchellh">twitter</a></footer></main></body></html>`,
		links: map[string]string{
			"Libghostty 1.0": "/writing/libghostty-1-0",
			"Superlogical":   "/writing/superlogical",
		},
	})
}
