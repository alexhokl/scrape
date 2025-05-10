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
		scraper := &GuardianScraper{}
		return scraper, nil

	default:
		return nil, fmt.Errorf("source %s is not supported", sourceType)
	}
}
