// Package scraper contains utility functions for HTML scraping.
package scraper

import (
	"strings"
)

func removeExtraSpaces(rawText string) string {
	strings := []rune(rawText)
	cleaned := make([]rune, 0, len(strings))
	spaceFound := false
	for _, r := range strings {
		if r == '\n' || r == '\r' {
			continue
		}
		if r == ' ' {
			if !spaceFound {
				cleaned = append(cleaned, r)
				spaceFound = true
			}
		} else {
			cleaned = append(cleaned, r)
			spaceFound = false
		}
	}

	return string(cleaned)
}

func trimSpacesAndLineBreaks(input string) string {
	return strings.Trim(
		strings.Trim(
			strings.Trim(input, " "),
			"\n",
		),
		" ",
	)
}

func generateFileNameFromTitle(title string) string {
	return strings.ToLower(
		strings.ReplaceAll(
			strings.ReplaceAll(title, " ", "_"),
			"__",
			"_",
		),
	)
}
