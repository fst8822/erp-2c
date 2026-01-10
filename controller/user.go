package controller

import (
	"erp-2c/dto/response"
	"erp-2c/service/use_cases"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserController struct {
	services *use_cases.Manager
}

func NewUserController(services *use_cases.Manager) *UserController {
	return &UserController{services: services}
}

func (c *UserController) GetById(w http.ResponseWriter, r *http.Request) {
	slog.Info("Get request GetById")

	idParam := chi.URLParam(r, "id")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		resp := response.BadRequest("failed to parse path variable", idParam)
		render.Status(r, resp.Code)
		render.JSON(w, r, response.BadRequest("failed parse", idParam))
		return
	}

	found, _ := c.services.UserService.GetById(userId)
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, response.OK("successful", found))

}

func (c *UserController) UpdateById(w http.ResponseWriter, r *http.Request) {
	slog.Info("PUT request UpdateById")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Test2"))
}

func (c *UserController) DeleteById(w http.ResponseWriter, r *http.Request) {
	slog.Info("Delete request DeleteById")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Test3"))
}
