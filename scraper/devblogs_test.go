package scraper

import "testing"

var _ ArticleScraper = (*DevBlogsScraper)(nil)
var _ LinkScraper = (*DevBlogsScraper)(nil)

func TestDevBlogsScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &DevBlogsScraper{},
		articleHTML: `<!DOCTYPE html><html><body><main id="main" class="devblogs-main"><div class="wrapper"><div class="container-evo mt-56">
<div id="single-wrapper" class="container-three-column-post">
<div class="pe-xl-198"><header class="ms-md-111 mb-40"><h1>Performance Improvements in .NET 11</h1><span>Stephen Toub</span></header></div>
<article id="post-60786" class="middle-column pe-xl-198">
<div class="entry-content sharepostcontent">
<p>Whew! .NET 11 is indeed one louder.</p>
<blockquote><p><strong>Note:</strong> benchmarks are on .NET 11 RC.</p></blockquote>
<pre><code class="language-csharp">IEnumerable&lt;int&gt; Values() =&gt; [1, 2, 3];</code></pre>
<pre class="prettyprint">dotnet run -c Release</pre>
</div></article>
<aside class="bh-sidebar"><a href="https://aka.ms/dotnet/twitter">Twitter</a></aside>
</div></div></div></main></body></html>`,
		title: "Performance Improvements in .NET 11",
		contains: []string{
			"# Performance Improvements in .NET 11\n\nWhew! .NET 11 is indeed one louder.\n\n",
			"> **Note:** benchmarks are on .NET 11 RC.\n\n",
			"```csharp\nIEnumerable<int> Values() => [1, 2, 3];\n```\n\n",
			"```\ndotnet run -c Release\n```\n",
		},
		excludes: []string{"Stephen Toub", "Twitter"},
		indexHTML: `<html><body>
<div class="pinned-wrap d-flex flex-column"><div class="w-100 pinned-right d-flex"><div><h3 class="fs-24 mb-20"><a class="single-click" href="/dotnet/announcing-dotnet-11/">Announcing .NET 11</a></h3></div></div></div>
<div class="masonry-card box bg-white"><div class="card-body"><div class="card-top"><h3 class="fs-24 fw-600 mb-16"><a class="single-click excerpt-title" href="/dotnet/nuget-staged-publishing/">Introducing staged publishing on NuGet.org</a></h3></div></div></div>
<span class="card-text listdisplay"><ul class="ftr-archives"><li><a href="/dotnet/2026/11/">November 2026</a></li></ul></span>
</body></html>`,
		links: map[string]string{
			"Announcing .NET 11":                         "/dotnet/announcing-dotnet-11/",
			"Introducing staged publishing on NuGet.org": "/dotnet/nuget-staged-publishing/",
		},
	})
}
