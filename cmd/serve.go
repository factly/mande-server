package cmd

import (
	"log"
	"net/http"

	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/model"
	"github.com/factly/mande-server/util/razorpay"
	"github.com/factly/x/meilisearchx"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts server for mande-server.",
	Run: func(cmd *cobra.Command, args []string) {
		// db setup
		model.SetupDB()

		meilisearchx.SetupMeiliSearch("data-portal", []string{"name", "slug", "description", "title", "contact_name", "contact_email", "license", "caption", "alt_text", "plan_name", "plan_info"})

		razorpay.SetupClient()

		// register routes
		userRouter := action.RegisterUserRoutes()
		adminRouter := action.RegisterAdminRoutes()
		webhookRouter := action.RegisterWebHookRoutes()

		go func() {
			log.Fatal(http.ListenAndServe(":7720", userRouter))
		}()

		go func() {
			log.Fatal(http.ListenAndServe(":7722", webhookRouter))
		}()

		log.Fatal(http.ListenAndServe(":7721", adminRouter))

	},
}
