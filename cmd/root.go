// Package cmd
package cmd

import (
	"os"

	"github.com/alexhokl/helper/cli"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "scrape",
	Short:        "A CLI tool scrape links and articles",
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scrape.yaml)")
}

func initConfig() {
	cli.ConfigureViper(cfgFile, "scrape", false, "")
}
