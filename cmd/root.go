package cmd

import (
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
	_ = rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scrape.yaml)")
}

func initConfig() {
	cli.ConfigureViper(cfgFile, "scrape", false, "")
}
