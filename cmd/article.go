package cmd

import (
	"fmt"
	"os"

	"github.com/alexhokl/scrape/scraper"
	"github.com/spf13/cobra"
)

type articleOptions struct {
	format string
	source string
	url    string
}

var articleOpts articleOptions

// articleCmd represents the article command
var articleCmd = &cobra.Command{
	Use:               "article",
	Short:             "Scrape an acticle",
	PersistentPreRunE: validateArticleOptions,
	RunE:              scrapeArticle,
}

func init() {
	rootCmd.AddCommand(articleCmd)

	flags := articleCmd.PersistentFlags()
	flags.StringVar(&articleOpts.format, "format", "markdown", "Output format (markdown)")
	flags.StringVar(&articleOpts.source, "source", "", "Source type (e.g., guardian, etc)")
	flags.StringVarP(&articleOpts.url, "url", "u", "", "URL of the article to scrape")

	articleCmd.MarkFlagRequired("source")
	articleCmd.MarkFlagRequired("url")
}

func validateArticleOptions(_ *cobra.Command, _ []string) error {
	opts := &articleOpts

	switch opts.format {
	case "markdown":
		break

	default:
		return fmt.Errorf("invalid format: %s", opts.format)
	}

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

func scrapeArticle(_ *cobra.Command, args []string) error {
	scraper, err := scraper.CreateArticleScraper(articleOpts.source)
	if err != nil {
		return fmt.Errorf("error creating scraper: %w", err)
	}
	markdownStr, err := scraper.ScrapeArticle(articleOpts.url)
	if err != nil {
		return fmt.Errorf("error scraping article: %w", err)
	}
	fmt.Fprintln(os.Stdout, markdownStr)

	return nil
}
