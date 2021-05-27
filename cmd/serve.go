package cmd

import (
	"log"
	"net/http"

	"github.com/dlmiddlecote/sqlstats"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/model"
	"github.com/factly/mande-server/util/razorpay"
	"github.com/factly/x/meilisearchx"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		meilisearchx.SetupMeiliSearch("mande", []string{"name", "slug", "description", "title", "contact_name", "contact_email", "license", "caption", "alt_text", "plan_name", "plan_info"})

		razorpay.SetupClient()

		// register routes
		userRouter := action.RegisterUserRoutes()
		adminRouter := action.RegisterAdminRoutes()
		webhookRouter := action.RegisterWebHookRoutes()

		go func() {
			promRouter := chi.NewRouter()

			sqlDB, _ := model.DB.DB()
			collector := sqlstats.NewStatsCollector(viper.GetString("database_name"), sqlDB)

			prometheus.MustRegister(collector)

			promRouter.Mount("/metrics", promhttp.Handler())
			log.Fatal(http.ListenAndServe(":8001", promRouter))
		}()

		go func() {
			log.Fatal(http.ListenAndServe(":8002", userRouter))
		}()

		go func() {
			log.Fatal(http.ListenAndServe(":8003", webhookRouter))
		}()

		log.Fatal(http.ListenAndServe(":8000", adminRouter))

	},
}
