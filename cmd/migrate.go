package cmd

import (
	"github.com/factly/mande-server/model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Applies database migrations for mande-server.",
	Run: func(cmd *cobra.Command, args []string) {
		// db setup
		model.SetupDB()

		model.Migration()
	},
}
