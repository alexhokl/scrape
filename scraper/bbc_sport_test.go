package scraper

import "testing"

var _ ArticleScraper = (*BBCSportScraper)(nil)
var _ LinkScraper = (*BBCSportScraper)(nil)

func TestBBCSportScraper(t *testing.T) {
	runSiteTests(t, siteTestCase{
		scraper: &BBCSportScraper{},
		articleHTML: `<!DOCTYPE html><html lang="en-GB"><body><div id="main-content"><article class="ssrcss-2fe4mx-ArticleWrapper e1nh2i2l0"><div spacing="6" class="ssrcss-15twdc6-Stack e1y4nx260">
<div data-block="headline" class="ssrcss-l7emlr-ComponentWrapper"><div><header data-testid="headline"><h1 id="main-heading" type="headline" tabindex="-1" class="ssrcss-8zqj65-Heading"><span role="text">Head coach Slegers leaves Arsenal</span></h1></header></div></div>
<div data-block="image" class="ssrcss-l7emlr-ComponentWrapper"><div><div data-testid="image"><figure class="ssrcss-hc6arm-StyledFigure"><div><span><picture><img alt="Renee Slegers claps" src="https://ichef.bbci.co.uk/ace/standard/800/cpsprodpb/e62d/live/slegers.jpg" width="800" height="450"/></picture></span><span role="text" class="ssrcss-tzgxt7-StyledFigureCopyright e34k3c20"><span class="visually-hidden ssrcss-i2z2ig-VisuallyHidden">Image source, </span>Getty Images</span></div><figcaption class="ssrcss-3uli8u-StyledFigureCaption"><span class="visually-hidden ssrcss-i2z2ig-VisuallyHidden">Image caption, </span><div><p class="ssrcss-1q0x1qg-Paragraph">Renee Slegers joined Arsenal in 2023</p></div></figcaption></figure></div></div></div>
<div data-block="byline" class="ssrcss-vxwoax-ComponentWrapper"><div><div data-testid="single-byline"><span class="visually-hidden">By</span><span class="ssrcss-nsjd43-TextContributorName">Sally Freedman</span><span>BBC Sport senior journalist</span></div></div></div>
<div data-block="metadata" class="ssrcss-l7emlr-ComponentWrapper"><div data-testid="metadata"><ul role="list"><li role="listitem"><div class="visually-hidden">Published</div><time data-testid="timestamp" dateTime="2026-11-27T15:35:33.715Z">27 November 2026</time></li></ul></div></div>
<div data-block="text" class="ssrcss-vxwoax-ComponentWrapper"><div data-testid="rich-text"><div><p class="ssrcss-1q0x1qg-Paragraph"><b>Arsenal have parted company with head coach Renee Slegers.</b></p><p class="ssrcss-1q0x1qg-Paragraph">Slegers was <a href="/sport/football/articles/c39n3d3pvygo" class="ssrcss-18834jl-InlineLink">handed the role permanently.</a></p></div></div></div>
<div data-block="links" class="ssrcss-l7emlr-ComponentWrapper"><ul role="list"><li><a href="/sport/football/articles/c8752llr1ed1o"><span><p><span aria-hidden="false">Why Arsenal are behind in race for Champions League spots</span></p></span></a></li></ul></div>
<div data-block="subheadline" class="ssrcss-l7emlr-ComponentWrapper"><div data-testid="subheadline"><h2 id="Pressure-was-building-for-weeks-analysis" tabindex="-1"><span role="text">Pressure was building for weeks - analysis</span></h2></div></div>
<div data-block="text" class="ssrcss-vxwoax-ComponentWrapper"><div data-testid="rich-text"><div><p class="ssrcss-1q0x1qg-Paragraph">A lack of goals has been central to Arsenal's struggles.</p></div></div></div>
<div data-block="media" class="ssrcss-l7emlr-ComponentWrapper"><figure><div data-testid="media-player-container-landscape"><p>Watch the video player</p></div></figure></div>
<div data-block="topicList" class="ssrcss-l7emlr-ComponentWrapper"><div data-testid="topic-list"><h2 type="normal">Related topics</h2><ul role="list"><li><a href="/sport/football/womens">Women's Football</a></li></ul></div></div>
<div data-block="promoList" class="ssrcss-l7emlr-ComponentWrapper"><h2 type="normal">More on this story</h2><ul><li><div data-testid="promo"><a href="/sounds/brand/p0bzmv49"><span role="text"><p><span aria-hidden="false">Women's Football Weekly feed</span></p></span></a></div></li></ul></div>
</div></article></div></body></html>`,
		title:    "Head coach Slegers leaves Arsenal",
		filename: "head_coach_slegers_leaves_arsenal",
		contains: []string{
			"**Arsenal have parted company with head coach Renee Slegers.**\n\n",
			"[handed the role permanently.](",
			"/sport/football/articles/c39n3d3pvygo)",
			"## Pressure was building for weeks - analysis\n\n",
			"A lack of goals has been central to Arsenal's struggles.",
			"![Renee Slegers claps](https://ichef.bbci.co.uk/ace/standard/800/cpsprodpb/e62d/live/slegers.jpg)",
			"*Renee Slegers joined Arsenal in 2023*",
		},
		excludes: []string{
			"Sally Freedman", "27 November", "Published", "Why Arsenal are behind",
			"video player", "Related topics", "More on this story", "Getty Images",
			"Image caption", "Image source",
		},
		indexHTML: `<html><body><div id="main-content"><h1 id="main-heading" type="title">Football</h1><ul role="list">
<li><div data-testid="promo" type="article"><div spacing="2"><a href="/sport/football/articles/ck0ev91zlwgpo"><span role="text"><p><span aria-hidden="false">How a football team rebuilt after plane crash</span></p></span></a><p>How a Brazilian football team rose again.</p></div><ul role="list"><li role="listitem"><div>Attribution</div><a href="/sport/football"><span type="attribution">Football</span></a></li></ul></div></li>
<li><div data-testid="promo" type="article"><div spacing="2"><a href="/sport/football/articles/cmy49d5kpw78o"><span role="text"><p><span aria-hidden="false">Head coach Slegers leaves Arsenal</span></p></span></a></div></div></li>
<li><div data-testid="promo" type="live"><div spacing="2"><a href="/sport/football/live/c0l1v3"><span role="text"><p><span aria-hidden="false">LIVE: Premier League news</span></p></span></a></div></div></li>
<li><div data-testid="promo"><div spacing="2"><a href="/sounds/brand/p0bzmv49"><span role="text"><p><span aria-hidden="false">Women's Football Weekly feed</span></p></span></a></div></div></li>
</ul></div></body></html>`,
		links: map[string]string{
			"How a football team rebuilt after plane crash": "/sport/football/articles/ck0ev91zlwgpo",
			"Head coach Slegers leaves Arsenal":             "/sport/football/articles/cmy49d5kpw78o",
		},
	})
}
