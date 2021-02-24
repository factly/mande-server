package cmd

import (
	"github.com/factly/data-portal-server/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "data-portal-server",
	Short: "data-portal-server is backend for MandE application",
	Long:  `MandE server is developed in Go. Manage datasets available for download in multiple formats including APIs.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(config.SetupVars)
}
