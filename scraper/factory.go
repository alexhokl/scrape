package scraper

import (
	"fmt"
)

func CreateLinkScraper(sourceType string) (LinkScraper, error) {
	switch sourceType {
	case "guardian":
		scraper := &GuardianScraper{}
		return scraper, nil

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

	default:
		return nil, fmt.Errorf("source %s is not supported", sourceType)
	}
}
