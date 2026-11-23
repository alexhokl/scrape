package scraper

import "testing"

var _ ArticleScraper = (*BellingcatScraper)(nil)
var _ LinkScraper = (*BellingcatScraper)(nil)

func TestBellingcatScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &BellingcatScraper{},
		articleHTML: `<!DOCTYPE html><html class="no-js"><body class="single-post"><main class="site__content"><div class="singular"><div class="container"><div class="singular__content">
<div class="singular__content__text">
<div class="singular__content__text__title"><h1>How Online Maps Lay Bare the US Military Bases Targeted by Iran</h1></div>
<div class="singular__content__meta">Sebastian Vandermeersch</div>
<div class="singular__content__text__content">
<p class="wp-block-paragraph">In September, CBS News obtained <a href="https://www.cbsnews.com/x">new images</a>.</p>
<figure class="wp-block-image size-large"><div class="media"><img src="https://www.bellingcat.com/app/uploads/2026/11/Map.jpg" alt="" class="wp-image-54241" loading="lazy"></div><figcaption class="wp-element-caption">Map of the bases.</figcaption></figure>
<h2 class="wp-block-heading">Camp Buehring</h2>
<blockquote class="twitter-tweet"><p lang="en" dir="ltr">Exclusive images released</p></blockquote>
<blockquote class="instagram-media" data-instgrm-permalink="https://www.instagram.com/p/x/"><div><a href="https://www.instagram.com/p/x/">View this post on Instagram</a></div></blockquote>
</div></div></div></div></div></main></body></html>`,
		title: "How Online Maps Lay Bare the US Military Bases Targeted by Iran",
		contains: []string{
			"In September, CBS News obtained [new images](https://www.cbsnews.com/x).\n\n",
			"![](https://www.bellingcat.com/app/uploads/2026/11/Map.jpg)\n\n*Map of the bases.*\n\n",
			"## Camp Buehring\n\n> Exclusive images released\n",
		},
		excludes: []string{"Sebastian Vandermeersch", "Instagram"},
		indexHTML: `<html><body><main class="site__content">
<div class="featured-posts container"><article class="featured-post"><div class="featured-post__content"><a class="featured-post__title" href="/news/2026/11/05/online-maps/">How Online Maps Lay Bare the US Military Bases</a></div>
<div class="featured-post__meta tags"><ul><li><a class="h6 underline-slide" href="/tag/middle-east">Middle East</a></li></ul></div></article></div>
<section class="section--investigations container"><div class="teasers--default"><article class="teaser--alt">
<a class="teaser__image--default image__wrapper" href="/news/2026/10/06/dafs-india/"><img src="x.jpg" alt=""></a>
<div class="teaser__text"><a class="teaser__title--default" href="/news/2026/10/06/dafs-india/">Groups Associated With India’s Far-Right</a></div></article></div></section>
</main></body></html>`,
		links: map[string]string{
			"How Online Maps Lay Bare the US Military Bases": "/news/2026/11/05/online-maps/",
			"Groups Associated With India’s Far-Right":       "/news/2026/10/06/dafs-india/",
		},
	})
}
