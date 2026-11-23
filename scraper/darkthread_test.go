package scraper

import "testing"

var _ ArticleScraper = (*DarkthreadScraper)(nil)
var _ LinkScraper = (*DarkthreadScraper)(nil)

func TestDarkthreadScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &DarkthreadScraper{},
		articleHTML: `<!DOCTYPE html><html><body><main><div class="content"><div class="blog-content">
<article class="post container" itemscope itemtype=http://schema.org/BlogPosting itemprop=blogPost>
<header><h1 class="post-title" itemprop="name headline"><a href="/blog/sample-post/" rel=canonical>2026 重回創世紀 IV</a></h1>
<table class=post-info><tr><td class=icon><img src=/img/calendar.svg alt="date time"></td><td class=fld>2026-11-22</td></tr></table></header>
<div itemprop="articleBody mainEntityOfPage">
<p>之前提過我迷上寫程式。延伸閱讀：<a href="/blog/tetris/">俄羅斯方塊</a></p>
<p><img alt="" src="/Posts/files/2026/Fig4.webp"></p>
<pre><code class="language-csharp">var x = 1;</code></pre>
<p class="hide" data-slug="revisit-ultima-4"> 隱藏 </p>
</div>
<footer><a href="/blog/sample-post/#comments">0 comments</a></footer>
</article></div></div></main></body></html>`,
		title: "2026 重回創世紀 IV",
		contains: []string{
			"之前提過我迷上寫程式。延伸閱讀：[俄羅斯方塊](http",
			"![](http",
			"/Posts/files/2026/Fig4.webp)\n\n",
			"```csharp\nvar x = 1;\n```\n",
		},
		excludes: []string{"隱藏", "calendar.svg", "2026-11-22", "0 comments"},
		indexHTML: `<html><body><main>
<article class="post container"><header><h1 itemprop="name headline"><a href="/blog/revisit-ultima-4/" rel=canonical>2026 重回創世紀 IV</a></h1></header>
<footer><a href="/blog/revisit-ultima-4/#comments">0 comments</a></footer></article>
<article class="post container"><header><h1 itemprop="name headline"><a href="/blog/nmcli/" rel=canonical>nmcli 常用指令</a></h1></header></article>
<footer><ul class="categories"><li><a href="/blog/category/python">Python</a></li></ul></footer>
</main></body></html>`,
		links: map[string]string{
			"2026 重回創世紀 IV": "/blog/revisit-ultima-4/",
			"nmcli 常用指令":    "/blog/nmcli/",
		},
	})
}
