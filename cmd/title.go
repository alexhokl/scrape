package cmd

import (
	"fmt"
	"os"

	"github.com/alexhokl/scrape/scraper"
	"github.com/spf13/cobra"
)

type titleOptions struct {
	source string
	url    string
}

var titleOpts titleOptions

// titleCmd represents the article command
var titleCmd = &cobra.Command{
	Use:               "title",
	Short:             "Scrape the title of an acticle",
	PersistentPreRunE: validateTitleOptions,
	RunE:              scrapeArticleTitle,
}

func init() {
	rootCmd.AddCommand(titleCmd)

	flags := titleCmd.PersistentFlags()
	flags.StringVar(&titleOpts.source, "source", "", "Source type (e.g., guardian, etc)")
	flags.StringVarP(&titleOpts.url, "url", "u", "", "URL of the article to scrape")

	titleCmd.MarkFlagRequired("source")
	titleCmd.MarkFlagRequired("url")
}

func validateTitleOptions(_ *cobra.Command, _ []string) error {
	opts := &titleOpts

	switch opts.source {
	case "guardian":
	case "microsoft":
	case "go":
	case "tofugu":

	default:
		return fmt.Errorf("invalid source: %s", opts.source)
	}

	if opts.url == "" {
		return fmt.Errorf("url is required")
	}

	return nil
}

func scrapeArticleTitle(_ *cobra.Command, args []string) error {
	scraper, err := scraper.CreateArticleScraper(titleOpts.source)
	if err != nil {
		return fmt.Errorf("error creating scraper: %w", err)
	}
	markdownStr, err := scraper.ScrapeTitle(titleOpts.url)
	if err != nil {
		return fmt.Errorf("error scraping title of article: %w", err)
	}
	fmt.Fprintln(os.Stdout, markdownStr)

	return nil
}
