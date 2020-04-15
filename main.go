package main

import (
	"net/http"

	"github.com/factly/data-portal-api/models"

	"github.com/factly/data-portal-api/actions"

	"github.com/go-chi/chi"
)

func registerRoutes() http.Handler {
	r := chi.NewRouter()

	r.Route("/currencies", func(r chi.Router) {
		r.Post("/", actions.CreateCurrency)
		r.Route("/{currencyId}", func(r chi.Router) {
			r.Get("/", actions.GetCurrency)
			r.Delete("/", actions.DeleteCurrency)
			r.Put("/", actions.UpdateCurrency)
		})
	})
	r.Route("/users", func(r chi.Router) {
		r.Post("/", actions.CreateUser)
		r.Route("/{userId}", func(r chi.Router) {
			r.Get("/", actions.GetUser)
			r.Delete("/", actions.DeleteUser)
			r.Put("/", actions.UpdateUser)
		})
	})
	r.Route("/plans", func(r chi.Router) {
		r.Post("/", actions.CreatePlan)
		r.Route("/{planId}", func(r chi.Router) {
			r.Get("/", actions.GetPlan)
			r.Delete("/", actions.DeletePlan)
			r.Put("/", actions.UpdatePlan)
		})

	})
	r.Route("/memberships", func(r chi.Router) {
		r.Post("/", actions.CreateMembership)
		r.Route("/{membershipId}", func(r chi.Router) {
			r.Get("/", actions.GetMembership)
			r.Delete("/", actions.DeleteMembership)
			r.Put("/", actions.UpdateMembership)
		})

	})
	r.Route("/payments", func(r chi.Router) {
		r.Post("/", actions.CreatePayment)
		r.Route("/{paymentId}", func(r chi.Router) {
			r.Get("/", actions.GetPayment)
			r.Delete("/", actions.DeletePayment)
			r.Put("/", actions.UpdatePayment)
		})
	})
	r.Route("/products", func(r chi.Router) {
		r.Post("/", actions.CreateProduct)
		r.Route("/{productId}", func(r chi.Router) {
			r.Get("/", actions.GetProduct)
			r.Delete("/", actions.DeleteProduct)
			r.Put("/", actions.UpdateProduct)
			r.Post("/type", actions.CreateProductType)
			r.Put("/type", actions.UpdateProductType)
			r.Delete("/type", actions.DeleteProductType)
			r.Post("/status", actions.CreateStatus)
			r.Put("/status", actions.UpdateStatus)
			r.Delete("/status", actions.DeleteStatus)
		})
	})
	r.Route("/tags", func(r chi.Router) {
		r.Post("/", actions.CreateTag)
		r.Route("/{tagId}", func(r chi.Router) {
			r.Get("/", actions.GetTag)
			r.Delete("/", actions.DeleteTag)
			r.Put("/", actions.UpdateTag)
		})

	})
	return r
}

func main() {
	// db setup
	models.SetupDB()

	// create tables
	models.DB.AutoMigrate(
		&models.Currency{},
		&models.Payment{},
		&models.Membership{},
		&models.Plan{},
		&models.User{},
		&models.Product{},
		&models.ProductType{},
		&models.Status{},
		&models.Tag{},
	)
	// Adding foreignKey
	models.DB.Model(&models.Payment{}).AddForeignKey("currency_id", "currencies(id)", "RESTRICT", "RESTRICT")
	models.DB.Model(&models.Membership{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	models.DB.Model(&models.Membership{}).AddForeignKey("plan_id", "plans(id)", "RESTRICT", "RESTRICT")
	models.DB.Model(&models.Membership{}).AddForeignKey("payment_id", "payments(id)", "RESTRICT", "RESTRICT")
	models.DB.Model(&models.Product{}).AddForeignKey("currency_id", "currencies(id)", "RESTRICT", "RESTRICT")
	models.DB.Model(&models.Product{}).AddForeignKey("status_id", "statuses(id)", "RESTRICT", "RESTRICT")
	models.DB.Model(&models.Product{}).AddForeignKey("product_type_id", "product_types(id)", "RESTRICT", "RESTRICT")
	r := registerRoutes()

	http.ListenAndServe(":3000", r)
}
