package scraper

import "testing"

var _ ArticleScraper = (*PublickeyScraper)(nil)
var _ LinkScraper = (*PublickeyScraper)(nil)

func TestPublickeyScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &PublickeyScraper{},
		articleHTML: `<!DOCTYPE html><html><body><div id="container" class="clearfix">
<div id="maincol" data-pagefind-body>
<h1>マイクロソフト、「Database Hub」正式リリース</h1>
<div class="pubdate">2026年11月20日</div><br />
<div class="entrybody clearfix">
<p>マイクロソフトは<a href="https://techcommunity.microsoft.com/blog/x">新サービスを発表しました</a>。</p>
<div class="imgCenter" style="width:547px"><img src="//www.publickey1.jp/2026/database-hub-ga01.png" alt="fig" width="545" height="305"></div>
<h2>AIエージェントによる運用監視の自動化も</h2>
<p><iframe src="https://www.youtube.com/embed/x"></iframe></p>
<ul><li><a href="/blog/26/related.html">関連記事</a></li></ul>
</div>
<h4>あわせて読みたい</h4>
<div class="sharebuttons"><a href="http://www.facebook.com/share.php">fbシェア</a></div>
</div>
<div class="category subcol"><a href="/cloud/">クラウド</a></div>
</div></body></html>`,
		title: "マイクロソフト、「Database Hub」正式リリース",
		contains: []string{
			"マイクロソフトは[新サービスを発表しました](https://techcommunity.microsoft.com/blog/x)。\n\n",
			"![fig](http",
			"/2026/database-hub-ga01.png)\n\n",
			"## AIエージェントによる運用監視の自動化も\n\n",
			"* [関連記事](http",
		},
		excludes: []string{"2026年11月20日", "あわせて読みたい", "fbシェア", "クラウド", "youtube"},
		indexHTML: `<html><body><div class="newentries maincol"><ul>
<li class="clearfix"><a href="/blog/26/database_hubai.html">Database Hub正式リリース</a></li>
<li class="clearfix"><a href="/blog/26/aicopilot_managed_runtime.html">Copilot Managed Runtime</a></li>
</ul></div>
<div class="backnumber maincol"><a href="/blog/202601.html">1月</a></div></body></html>`,
		links: map[string]string{
			"Database Hub正式リリース":      "/blog/26/database_hubai.html",
			"Copilot Managed Runtime": "/blog/26/aicopilot_managed_runtime.html",
		},
	})
}
