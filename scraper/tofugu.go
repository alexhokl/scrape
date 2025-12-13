package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type TofuguScraper struct {
}

// ScrapeArticle scrapes the article content from the specified URL
// and returns markdown in a string
func (g *TofuguScraper) ScrapeArticle(url string) (string, error) {
	c := colly.NewCollector()

	var markdown string

	// title
	c.OnHTML("h1.article-title", func(e *colly.HTMLElement) {
		markdown += parseTofuguTitle(e)
	})

	// minor title
	c.OnHTML("div.article-header-elements ul.meta", func(e *colly.HTMLElement) {
		markdown += fmt.Sprintf("%s\n\n", trimSpacesAndLineBreaks(e.Text))
	})

	// article body
	c.OnHTML("article div.main", func(e *colly.HTMLElement) {
		markdown += parseTofuguArticle(e)
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func parseTofuguTitle(e *colly.HTMLElement) string {
	builder := strings.Builder{}
	builder.WriteString("# ")
	builder.WriteString(trimSpacesAndLineBreaks(e.Text))
	builder.WriteString("\n\n")
	return builder.String()
}

func parseTofuguArticle(e *colly.HTMLElement) string {
	builder := strings.Builder{}

	// builder.WriteString(e.ChildText("div.short-explanation"))
	// builder.WriteString("\n\n")

	// iterate each child element
	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		if child.DOM.Parent().IsSelection(e.DOM) {
			switch child.Name {
			case "h2":
				builder.WriteString("## ")
				builder.WriteString(removeExtraSpaces(child.Text))
				builder.WriteString("\n\n")
			case "h3":
				builder.WriteString("### ")
				builder.WriteString(removeExtraSpaces(child.Text))
				builder.WriteString("\n\n")
			case "h4":
				builder.WriteString("#### ")
				builder.WriteString(removeExtraSpaces(child.Text))
				builder.WriteString("\n\n")
			case "h5":
				builder.WriteString("##### ")
				builder.WriteString(removeExtraSpaces(child.Text))
				builder.WriteString("\n\n")
			case "p":
				builder.WriteString(removeExtraSpaces(child.Text))
				builder.WriteString("\n\n")
			case "ul":
				if child.DOM.HasClass("example-sentence") {
					builder.WriteString(parseExampleList(child))
					break
				}

				// assume it is table of contents
				builder.WriteString(parseTableOfContents(child))

			case "ol":
				child.ForEach("li", func(index int, li *colly.HTMLElement) {
					builder.WriteString(fmt.Sprintf("%d. ", index+1))
					builder.WriteString(trimSpacesAndLineBreaks(li.Text))
					builder.WriteString("\n")
				})
				builder.WriteString("\n")
			case "table":
				builder.WriteString(parseTofuguTable(child))
			case "blockquote":
				builder.WriteString("> ")
				builder.WriteString(parseTofuguBlockquote(child.Text))
				builder.WriteString("\n\n")
			}
		}
	})

	return builder.String()
}

func parseTableOfContents(e *colly.HTMLElement) string {
	builder := strings.Builder{}
	e.ForEach("li", func(_ int, li *colly.HTMLElement) {
		if len(li.DOM.ParentsFiltered("ul").Nodes) > 1 {
			// skip it is a child as it should have been processed in the code
			// below already
			return
		}

		builder.WriteString(parseTofuguListItem(li, 0))
		li.ForEach("ul li", func(_ int, subli *colly.HTMLElement) {
			builder.WriteString(parseTofuguListItem(subli, 1))
			subli.ForEach("ul li", func(_ int, subsubli *colly.HTMLElement) {
				builder.WriteString(parseTofuguListItem(subsubli, 2))
			})
		})
	})
	builder.WriteString("\n")

	return builder.String()
}

func parseTofuguListItem(e *colly.HTMLElement, level int) string {
	builder := strings.Builder{}
	for range level {
		builder.WriteString("  ")
	}
	builder.WriteString("* ")
	builder.WriteString(firstLine(e.Text))
	builder.WriteString("\n")
	return builder.String()
}

func parseExampleList(e *colly.HTMLElement) string {
	builder := strings.Builder{}
	builder.WriteString("Example\n\n")
	e.ForEach("li", func(index int, li *colly.HTMLElement) {
		switch index {
		case 0:
			builder.WriteString("- Japanese:\n")
		case 1:
			builder.WriteString("- English:\n")
		}
		sentence := fmt.Sprintf("  * %s\n", trimSpacesAndLineBreaks(li.Text))
		builder.WriteString(sentence)
		builder.WriteString("\n")
	})
	builder.WriteString("\n")

	return builder.String()
}

func firstLine(s string) string {
	lines := strings.SplitN(s, "\n", 2)
	return lines[0]
}

func parseTofuguTable(table *colly.HTMLElement) string {
	builder := strings.Builder{}
	table.ForEach("tr", func(rowIndex int, tr *colly.HTMLElement) {
		tr.ForEach("th", func(_ int, th *colly.HTMLElement) {
			builder.WriteString("| ")
			builder.WriteString(parseTofuguTableCell(th.Text))
			builder.WriteString(" ")
		})
		if builder.Len() > 0 && rowIndex == 0 {
			// render separator row
			builder.WriteString("|\n")
			tr.ForEach("th", func(_ int, th *colly.HTMLElement) {
				builder.WriteString("| --- ")
			})
		}
		if builder.Len() == 0 && rowIndex == 0 {
			// render dummy header row
			tr.ForEach("td", func(_ int, td *colly.HTMLElement) {
				builder.WriteString("| ")
			})
			builder.WriteString("|\n")
			// render separator row
			tr.ForEach("td", func(_ int, th *colly.HTMLElement) {
				builder.WriteString("|---")
			})
			builder.WriteString("|\n")
		}
		tr.ForEach("td", func(_ int, td *colly.HTMLElement) {
			builder.WriteString("| ")
			builder.WriteString(parseTofuguTableCell(td.Text))
			builder.WriteString(" ")
		})
		builder.WriteString("|\n")
	})
	builder.WriteString("\n")

	return builder.String()
}

func parseTofuguTableCell(text string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			trimSpacesAndLineBreaks(text),
			"\n\n",
			"\n",
		),
		"\n",
		"; ",
	)
}

func parseTofuguBlockquote(text string) string {
	return strings.ReplaceAll(trimSpacesAndLineBreaks(text), "\n", "\n> ")
}
