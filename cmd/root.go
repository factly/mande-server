package cmd

import (
	"github.com/factly/mande-server/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mande-server",
	Short: "mande-server is backend for MandE application",
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
