// classify project classify.go
package arena

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestCompare(t *testing.T) {
	n1, err := html.Parse(strings.NewReader(`
<div class="u-size8of12 u-xs-size12of12 u-minHeight400 u-xs-height350 u-overflowHidden js-trackedPost u-relative u-imageSpectrum" data-source="collection_home---7------0----------------" data-post-id="6304c4d25ed8" data-scroll="native">
	<a href="https://hackernoon.com/insider-reflections-on-the-ico-bubble-6304c4d25ed8?source=collection_home---7------0----------------" class="u-block u-absolute0 u-backgroundSizeCover u-backgroundOriginBorderBox u-borderLighter u-borderBox u-backgroundColorGrayLight" style="background-image: url(&quot;https://cdn-images-1.medium.com/max/2000/gradv/29/81/30/darken/25/1*T4zNZcQWg4-LPUvoBppGDQ.jpeg&quot;); background-position: 50% 50% !important;">
		<span class="u-textScreenReader">Insider Reflections on The ICO Bubble</span>
	</a>
	<div class="u-absolute u-left0 u-bottom0 u-right20Percent u-baseColor--textDark u-xs-marginLeft20 u-xs-marginRight20 u-marginLeft40 u-marginTop30 u-marginRight40 u-marginBottom30">
		<a class="" href="https://hackernoon.com/insider-reflections-on-the-ico-bubble-6304c4d25ed8?source=collection_home---7------0----------------" data-action-source="collection_home---7------0----------------" data-post-id="6304c4d25ed8">
			<h3 class="u-contentSansBold u-lineHeightTightest u-xs-fontSize24 u-paddingBottom2 u-paddingTop5 u-fontSize32">
				<div class="">Insider Reflections on The ICO&nbsp;Bubble</div>
			</h3>
			<div class="u-contentSansThin u-lineHeightBaseSans u-fontSize24 u-xs-fontSize18 u-textColorNormal u-baseColor--textNormal">
				<div class="u-fontSize18">According to Forbes, more than $2.3 billion has been raised in token sales, aka “ICO’s” so far in 2017. As an entrepreneur who has&nbsp;raised…</div>
			</div>
		</a>
		<div class="u-clearfix u-marginTop20">
			<div class="u-flexCenter">
				<div class="postMetaInline-avatar u-flex0">
					<a class="link u-baseColor--link avatar" href="https://hackernoon.com/@betashop" data-action="show-user-card" data-action-value="e1698c392e7a" data-action-type="hover" data-user-id="e1698c392e7a" data-collection-slug="hacker-daily" dir="auto">
						<div class="u-relative u-inlineBlock u-flex0">
							<img src="./Hacker Noon_files/1-k2VoaTzWs_eDzJIeJ98HoA.jpeg" class="avatar-image u-size36x36 u-xs-size32x32" alt="Go to the profile of Jason Goldberg">
								<div class="avatar-halo u-absolute u-baseColor--iconNormal u-textColorGreenNormal svgIcon" style="width: calc(100% + 10px); height: calc(100% + 10px); top:-5px; left:-5px">
									<svg viewBox="0 0 40 40"
										xmlns="http://www.w3.org/2000/svg">
										<path d="M3.44615311,11.6601601 C6.57294867,5.47967718 12.9131553,1.5 19.9642857,1.5 C27.0154162,1.5 33.3556228,5.47967718 36.4824183,11.6601601 L37.3747245,11.2087295 C34.0793076,4.69494641 27.3961457,0.5 19.9642857,0.5 C12.5324257,0.5 5.84926381,4.69494641 2.55384689,11.2087295 L3.44615311,11.6601601 Z"></path>
										<path d="M36.4824183,28.2564276 C33.3556228,34.4369105 27.0154162,38.4165876 19.9642857,38.4165876 C12.9131553,38.4165876 6.57294867,34.4369105 3.44615311,28.2564276 L2.55384689,28.7078582 C5.84926381,35.2216412 12.5324257,39.4165876 19.9642857,39.4165876 C27.3961457,39.4165876 34.0793076,35.2216412 37.3747245,28.7078582 L36.4824183,28.2564276 Z"></path>
									</svg>
								</div>
							</div>
						</a>
					</div>
					<div class="postMetaInline postMetaInline-authorLockup ui-captionStrong u-flex1 u-noWrapWithEllipsis">
						<a class="ds-link ds-link--styleSubtle link link--darken link--accent u-accentColor--textNormal u-accentColor--textDarken" href="https://hackernoon.com/@betashop" data-action="show-user-card" data-action-value="e1698c392e7a" data-action-type="hover" data-user-id="e1698c392e7a" data-collection-slug="hacker-daily" dir="auto">Jason Goldberg</a>
						<div class="ui-caption u-fontSize12 u-baseColor--textNormal u-textColorNormal js-postMetaInlineSupplemental">
							<time datetime="2017-11-12T15:39:06.921Z">Nov 12</time>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
	`))
	if err != nil {
		t.Error(err)
	}
	a1 := NewArena(*n1)

	n2, err := html.Parse(strings.NewReader(`
<div class="col u-xs-size12of12 js-trackedPost u-paddingLeft12 u-marginBottom15 u-paddingRight12 u-size4of12" data-source="collection_home---4------0----------------" data-post-id="984a72d48ebc" data-index="0" data-scroll="native">
	<div class="u-lineHeightBase postItem">
		<a href="https://hackernoon.com/all-you-need-to-know-about-css-in-js-984a72d48ebc?source=collection_home---4------0----------------" data-action="open-post" data-action-value="https://hackernoon.com/all-you-need-to-know-about-css-in-js-984a72d48ebc?source=collection_home---4------0----------------" class="u-block u-xs-height170 u-height172 u-backgroundSizeCover u-backgroundOriginBorderBox u-backgroundColorGrayLight u-borderLighter" style="background-image: url(&quot;https://cdn-images-1.medium.com/max/400/1*OH0dDfJAGv6aEGFk2zLAxw.jpeg&quot;); background-position: 50% 50% !important;">
			<span class="u-textScreenReader">All You Need To Know About CSS-in-JS</span>
		</a>
	</div>
	<div class="col u-xs-marginBottom10 u-paddingLeft0 u-paddingRight0 u-paddingTop15 u-marginBottom30">
		<a class="" href="https://hackernoon.com/all-you-need-to-know-about-css-in-js-984a72d48ebc?source=collection_home---4------0----------------" data-action-source="collection_home---4------0----------------" data-post-id="984a72d48ebc">
			<h3 class="u-contentSansBold u-lineHeightTightest u-xs-fontSize24 u-paddingBottom2 u-paddingTop5 u-fontSize32">
				<div class="u-letterSpacingTight u-lineHeightTighter u-fontSize24">All You Need To Know About CSS-in-JS</div>
			</h3>
			<div class="u-contentSansThin u-lineHeightBaseSans u-fontSize24 u-xs-fontSize18 u-textColorNormal u-baseColor--textNormal">
				<div class="u-fontSize18 u-letterSpacingTight u-lineHeightTight u-marginTop7 u-textColorNormal u-baseColor--textNormal">TL;DR: Thinking in components — No longer do you have to maintain bunch of style-sheets. CSS-in-JS abstracts the CSS model to the component…</div>
			</div>
		</a>
		<div class="u-clearfix u-marginTop20">
			<div class="u-flexCenter">
				<div class="postMetaInline-avatar u-flex0">
					<a class="link u-baseColor--link avatar" href="https://hackernoon.com/@wesharehoodies" data-action="show-user-card" data-action-value="ce572601b7e" data-action-type="hover" data-user-id="ce572601b7e" data-collection-slug="hacker-daily" dir="auto">
						<div class="u-relative u-inlineBlock u-flex0">
							<img src="./Hacker Noon_files/1-kOSpywJF8sa7Gac3vYXTtA.jpeg" class="avatar-image u-size36x36 u-xs-size32x32" alt="Go to the profile of Indrek Lasn">
								<div class="avatar-halo u-absolute u-baseColor--iconNormal u-textColorGreenNormal svgIcon" style="width: calc(100% + 10px); height: calc(100% + 10px); top:-5px; left:-5px">
									<svg viewBox="0 0 40 40"
										xmlns="http://www.w3.org/2000/svg">
										<path d="M3.44615311,11.6601601 C6.57294867,5.47967718 12.9131553,1.5 19.9642857,1.5 C27.0154162,1.5 33.3556228,5.47967718 36.4824183,11.6601601 L37.3747245,11.2087295 C34.0793076,4.69494641 27.3961457,0.5 19.9642857,0.5 C12.5324257,0.5 5.84926381,4.69494641 2.55384689,11.2087295 L3.44615311,11.6601601 Z"></path>
										<path d="M36.4824183,28.2564276 C33.3556228,34.4369105 27.0154162,38.4165876 19.9642857,38.4165876 C12.9131553,38.4165876 6.57294867,34.4369105 3.44615311,28.2564276 L2.55384689,28.7078582 C5.84926381,35.2216412 12.5324257,39.4165876 19.9642857,39.4165876 C27.3961457,39.4165876 34.0793076,35.2216412 37.3747245,28.7078582 L36.4824183,28.2564276 Z"></path>
									</svg>
								</div>
							</div>
						</a>
					</div>
					<div class="postMetaInline postMetaInline-authorLockup ui-captionStrong u-flex1 u-noWrapWithEllipsis">
						<a class="ds-link ds-link--styleSubtle link link--darken link--accent u-accentColor--textNormal u-accentColor--textDarken" href="https://hackernoon.com/@wesharehoodies" data-action="show-user-card" data-action-value="ce572601b7e" data-action-type="hover" data-user-id="ce572601b7e" data-collection-slug="hacker-daily" dir="auto">Indrek Lasn</a>
						<div class="ui-caption u-fontSize12 u-baseColor--textNormal u-textColorNormal js-postMetaInlineSupplemental">
							<time datetime="2017-11-09T22:35:08.012Z">Nov 10</time>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
	`))
	if err != nil {
		t.Error(err)
	}
	a2 := NewArena(*n2)

	n3, err := html.Parse(strings.NewReader(`
<div class="u-size8of12 u-xs-size12of12 u-minHeight400 u-xs-height350 u-overflowHidden js-trackedPost u-relative u-imageSpectrum" data-source="collection_home---7------0----------------" data-post-id="24d184e871a3">
	<a href="https://hackernoon.com/this-startup-raised-35-943-270-of-funding-via-28-kickstarter-campaigns-24d184e871a3?source=collection_home---7------0----------------" class="u-block u-absolute0 u-backgroundSizeCover u-backgroundOriginBorderBox u-borderLighter u-borderBox u-backgroundColorGrayLight" style="background-image: url(&quot;https://cdn-images-1.medium.com/max/2000/gradv/29/81/30/darken/25/1*ZjoUGyrzdaWcYop-lkfyjQ.jpeg&quot;); background-position: 50% 50% !important;">
		<span class="u-textScreenReader">This Startup Raised $35,943,270 of Funding via 28 Kickstarter Campaigns</span>
	</a>
	<div class="u-absolute u-left0 u-bottom0 u-right20Percent u-baseColor--textDark u-xs-marginLeft20 u-xs-marginRight20 u-marginLeft40 u-marginTop30 u-marginRight40 u-marginBottom30">
		<a class="" href="https://hackernoon.com/this-startup-raised-35-943-270-of-funding-via-28-kickstarter-campaigns-24d184e871a3?source=collection_home---7------0----------------" data-action-source="collection_home---7------0----------------" data-post-id="24d184e871a3">
			<h3 class="u-contentSansBold u-lineHeightTightest u-xs-fontSize24 u-paddingBottom2 u-paddingTop5 u-fontSize32">
				<div class="">This Startup Raised $35,943,270 of Funding via 28 Kickstarter Campaigns</div>
			</h3>
			<div class="u-contentSansThin u-lineHeightBaseSans u-fontSize24 u-xs-fontSize18 u-textColorNormal u-baseColor--textNormal">
				<div class="u-fontSize18">CMON is a cross Warby Parker and Parker Brothers—and is the first company to go from Kickstarter to&nbsp;IPO.</div>
			</div>
		</a>
		<div class="u-clearfix u-marginTop20">
			<div class="u-flexCenter">
				<div class="postMetaInline-avatar u-flex0">
					<a class="link u-baseColor--link avatar" href="https://hackernoon.com/@foundercollective" data-action="show-user-card" data-action-value="f49435c6fa9" data-action-type="hover" data-user-id="f49435c6fa9" data-collection-slug="hacker-daily" dir="auto">
						<div class="u-relative u-inlineBlock u-flex0">
							<img src="./Hacker Noon_files/1-q-oDqs62LJTpqS5zOa0F9g.png" class="avatar-image u-size36x36 u-xs-size32x32" alt="Go to the profile of Founder Collective">
								<div class="avatar-halo u-absolute u-baseColor--iconNormal u-textColorGreenNormal svgIcon" style="width: calc(100% + 10px); height: calc(100% + 10px); top:-5px; left:-5px">
									<svg viewBox="0 0 40 40"
										xmlns="http://www.w3.org/2000/svg">
										<path d="M3.44615311,11.6601601 C6.57294867,5.47967718 12.9131553,1.5 19.9642857,1.5 C27.0154162,1.5 33.3556228,5.47967718 36.4824183,11.6601601 L37.3747245,11.2087295 C34.0793076,4.69494641 27.3961457,0.5 19.9642857,0.5 C12.5324257,0.5 5.84926381,4.69494641 2.55384689,11.2087295 L3.44615311,11.6601601 Z"></path>
										<path d="M36.4824183,28.2564276 C33.3556228,34.4369105 27.0154162,38.4165876 19.9642857,38.4165876 C12.9131553,38.4165876 6.57294867,34.4369105 3.44615311,28.2564276 L2.55384689,28.7078582 C5.84926381,35.2216412 12.5324257,39.4165876 19.9642857,39.4165876 C27.3961457,39.4165876 34.0793076,35.2216412 37.3747245,28.7078582 L36.4824183,28.2564276 Z"></path>
									</svg>
								</div>
							</div>
						</a>
					</div>
					<div class="postMetaInline postMetaInline-authorLockup ui-captionStrong u-flex1 u-noWrapWithEllipsis">
						<a class="ds-link ds-link--styleSubtle link link--darken link--accent u-accentColor--textNormal u-accentColor--textDarken" href="https://hackernoon.com/@foundercollective" data-action="show-user-card" data-action-value="f49435c6fa9" data-action-type="hover" data-user-id="f49435c6fa9" data-collection-slug="hacker-daily" dir="auto">Founder Collective</a>
						<div class="ui-caption u-fontSize12 u-baseColor--textNormal u-textColorNormal js-postMetaInlineSupplemental">
							<time datetime="2017-11-10T13:59:03.985Z">Nov 10</time>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
	`))
	if err != nil {
		t.Error(err)
	}
	a3 := NewArena(*n3)

	res := CmpDeepRate(a1, a2, 1, 1)
	t.Log(res.Sum, res.Count)

	res = CmpDeepRate(a1, a3, 1, 1)
	t.Log(res.Sum, res.Count)

	//t.Log(a1, a2, a3)
}
