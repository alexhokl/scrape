package cmd

import (
	"fmt"
	"os"

	"github.com/alexhokl/scrape/scraper"
	"github.com/spf13/cobra"
)

type filenameOptions struct {
	source string
	url    string
}

var filenameOpts filenameOptions

// filenameCmd represents the article command
var filenameCmd = &cobra.Command{
	Use:               "filename",
	Short:             "Scrape filename of an acticle",
	PersistentPreRunE: validateFilenameOptions,
	RunE:              scrapeArticleFilename,
}

func init() {
	rootCmd.AddCommand(filenameCmd)

	flags := filenameCmd.PersistentFlags()
	flags.StringVar(&filenameOpts.source, "source", "", "Source type (e.g., guardian, etc)")
	flags.StringVarP(&filenameOpts.url, "url", "u", "", "URL of the article to scrape")

	filenameCmd.MarkFlagRequired("source")
	filenameCmd.MarkFlagRequired("url")
}

func validateFilenameOptions(_ *cobra.Command, _ []string) error {
	opts := &filenameOpts

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

func scrapeArticleFilename(_ *cobra.Command, args []string) error {
	scraper, err := scraper.CreateArticleScraper(filenameOpts.source)
	if err != nil {
		return fmt.Errorf("error creating scraper: %w", err)
	}
	markdownStr, err := scraper.ScrapeFilename(filenameOpts.url)
	if err != nil {
		return fmt.Errorf("error scraping filename of article: %w", err)
	}
	fmt.Fprintln(os.Stdout, markdownStr)

	return nil
}
