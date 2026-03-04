package routers

import (
	"erp-2c/controller"
	"erp-2c/lib/observability/app_metrics"
	"erp-2c/security"
	"erp-2c/service/use_cases"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(serviceManager *use_cases.Manager) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	validate := validator.New()

	authController := controller.NewAuthController(serviceManager, validate)
	userController := controller.NewUserController(serviceManager, validate)
	productController := controller.NewProductController(serviceManager, validate)
	deliveryController := controller.NewDeliveryController(serviceManager, validate)
	notifyController := controller.NewNotifyController(serviceManager)

	router.Handle("/metrics", promhttp.Handler())
	router.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(app_metrics.MetricsMiddleware)
			r.Route("/auth", func(r chi.Router) {
				r.Post("/signup", authController.SignUp)
				r.Post("/signin", authController.SignIn)
			})
		})
		r.Group(func(r chi.Router) {
			r.Use(app_metrics.MetricsMiddleware)
			r.Use(security.JwtMiddleware)

			r.Route("/user", func(r chi.Router) {
				r.Get("/{id}", userController.GetById)
				r.Put("/{id}", userController.Save)
			})

			r.Route("/product", func(r chi.Router) {
				r.Post("/", productController.Save)
				r.Get("/", productController.GetAll)
				r.Get("/{id}", productController.GetById)
				r.Put("/{id}", productController.UpdateById)
				r.Delete("/{id}", productController.DeleteById)
			})

			r.Route("/delivery", func(r chi.Router) {
				r.Post("/", deliveryController.Save)
				r.Get("/", deliveryController.GetAll)
				r.Get("/{id}", deliveryController.GetById)
				r.Put("/{id}", deliveryController.UpdateById)
				r.Delete("/{id}", deliveryController.DeleteById)
			})

			r.Route("/ws", func(r chi.Router) {
				r.Get("/subscribe", notifyController.UpgradeConnection)
			})
		})
	})

	return router
}
