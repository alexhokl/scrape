package cmd

import (
	"fmt"
	"os"

	"github.com/alexhokl/scrape/scraper"
	"github.com/spf13/cobra"
)

type linksOptions struct {
	source string
	url    string
}

var linksOpts linksOptions

// linksCmd represents the links command
var linksCmd = &cobra.Command{
	Use:               "links",
	Short:             "scrape links",
	PersistentPreRunE: validateLinksOptions,
	RunE:              scrapeLinks,
}

func init() {
	rootCmd.AddCommand(linksCmd)

	flags := linksCmd.PersistentFlags()
	flags.StringVarP(&linksOpts.url, "url", "u", "", "URL of the links to scrape")
	flags.StringVar(&linksOpts.source, "source", "", "Source type (e.g., guardian, etc)")

	linksCmd.MarkFlagRequired("source")
	linksCmd.MarkFlagRequired("url")
}

func validateLinksOptions(_ *cobra.Command, _ []string) error {
	opts := &linksOpts

	switch opts.source {
	case "guardian":
		break

	default:
		return fmt.Errorf("invalid source: %s", opts.source)
	}

	if opts.url == "" {
		return fmt.Errorf("url is required")
	}

	return nil
}

func scrapeLinks(_ *cobra.Command, args []string) error {
	scraper, err := scraper.CreateLinkScraper(linksOpts.source)
	if err != nil {
		return fmt.Errorf("error creating scraper: %w", err)
	}
	links, err := scraper.ScrapeLinks(linksOpts.url)
	if err != nil {
		return fmt.Errorf("error scraping links: %w", err)
	}

	for linkText, articleURL := range links {
		fmt.Fprintf(os.Stdout, "[%v](%v)\n", linkText, articleURL)
	}

	return nil
}
