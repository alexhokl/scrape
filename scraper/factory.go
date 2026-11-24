package scraper

import (
	"fmt"
)

func CreateLinkScraper(sourceType string) (LinkScraper, error) {
	switch sourceType {
	case "guardian":
		scraper := &GuardianScraper{}
		return scraper, nil
	case "publickey":
		return &PublickeyScraper{}, nil
	case "proton":
		return &ProtonScraper{}, nil
	case "darkthread":
		return &DarkthreadScraper{}, nil
	case "rakyll":
		return &RakyllScraper{}, nil
	case "signal":
		return &SignalScraper{}, nil
	case "charm":
		return &CharmScraper{}, nil
	case "filippo":
		return &FilippoScraper{}, nil
	case "juliaevans":
		return &JuliaEvansScraper{}, nil
	case "kleppmann":
		return &KleppmannScraper{}, nil
	case "mitchellh":
		return &MitchellHashimotoScraper{}, nil
	case "jessfraz":
		return &JessFrazelleScraper{}, nil
	case "rickstrahl":
		return &RickStrahlScraper{}, nil
	case "hanselman":
		return &HanselmanScraper{}, nil
	case "devblogs":
		return &DevBlogsScraper{}, nil
	case "observationdeck":
		return &ObservationDeckScraper{}, nil
	case "bellingcat":
		return &BellingcatScraper{}, nil
	case "medium":
		return &MediumScraper{}, nil

	default:
		return nil, fmt.Errorf("source %s is not supported", sourceType)
	}
}

func CreateArticleScraper(sourceType string) (ArticleScraper, error) {
	switch sourceType {
	case "guardian":
		return &GuardianScraper{}, nil
	case "microsoft":
		return &MicrosoftLearnScraper{}, nil
	case "go":
		return &GoDocScraper{}, nil
	case "tofugu":
		return &TofuguScraper{}, nil
	case "newyorktimes":
		return &NewYorkTimesScraper{}, nil
	case "tailscale":
		return &TailscaleScraper{}, nil
	case "cloudflare":
		return &CloudflareScraper{}, nil
	case "wikipedia":
		return &WikipediaScraper{}, nil
	case "ollama":
		return &OllamaScraper{}, nil
	case "grafana":
		return &GrafanaScraper{}, nil
	case "unseenjapan":
		return &UnseenJapanScraper{}, nil
	case "propublica":
		return &ProPublicaScraper{}, nil
	case "japantimes":
		return &JapanTimesScraper{}, nil
	case "publickey":
		return &PublickeyScraper{}, nil
	case "proton":
		return &ProtonScraper{}, nil
	case "darkthread":
		return &DarkthreadScraper{}, nil
	case "rakyll":
		return &RakyllScraper{}, nil
	case "signal":
		return &SignalScraper{}, nil
	case "charm":
		return &CharmScraper{}, nil
	case "filippo":
		return &FilippoScraper{}, nil
	case "juliaevans":
		return &JuliaEvansScraper{}, nil
	case "kleppmann":
		return &KleppmannScraper{}, nil
	case "mitchellh":
		return &MitchellHashimotoScraper{}, nil
	case "jessfraz":
		return &JessFrazelleScraper{}, nil
	case "rickstrahl":
		return &RickStrahlScraper{}, nil
	case "hanselman":
		return &HanselmanScraper{}, nil
	case "devblogs":
		return &DevBlogsScraper{}, nil
	case "observationdeck":
		return &ObservationDeckScraper{}, nil
	case "bellingcat":
		return &BellingcatScraper{}, nil
	case "medium":
		return &MediumScraper{}, nil

	default:
		return nil, fmt.Errorf("source %s is not supported", sourceType)
	}
}
