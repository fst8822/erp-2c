package routers

import (
	"erp-2c/controller"
	"erp-2c/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(serviceManager *service.Manager) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	authController := controller.NewAuthController(serviceManager)
	userController := controller.NewUserController(serviceManager)
	productController := controller.NewProductController(serviceManager)

	router.Route("/api/v1", func(r chi.Router) {

		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", authController.SignUp)
			r.Post("/signin", authController.SignIn)
		})

		r.Group(func(r chi.Router) {
			//r.Use(middleware.BasicAuth("user-are", map[string]string{
			//	"admin": "admin",
			//}))

			r.Route("/user", func(r chi.Router) {

				r.Get("/{id}", userController.GetById)
				r.Get("/{name}", userController.GetByName)
				r.Put("/{id}", userController.UpdateById)
				r.Delete("/{id}", userController.DeleteById)
			})

			r.Route("/product", func(r chi.Router) {
				r.Post("/", productController.Save)
				r.Get("/", productController.GetAll)
				r.Get("/{id}", productController.GetById)
				r.Get("/{name}", productController.GetByName)
				r.Put("/{id}", productController.UpdateById)
				r.Delete("/{id}", productController.DeleteById)
			})

		})
	})

	return router
}
