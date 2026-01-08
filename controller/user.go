package controller

import (
	"erp-2c/service/use_cases"
	"log/slog"
	"net/http"
)

type UserController struct {
	services *use_cases.Manager
}

func NewUserController(services *use_cases.Manager) *UserController {
	return &UserController{services: services}
}

func (c *UserController) GetById(w http.ResponseWriter, r *http.Request) {
	slog.Info("Get request GetById")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Test"))
}

func (c *UserController) GetByName(w http.ResponseWriter, r *http.Request) {
	slog.Info("Get request GetById")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Test"))
}

func (c *UserController) UpdateById(w http.ResponseWriter, r *http.Request) {
	slog.Info("PUT request UpdateById")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Test"))
}

func (c *UserController) DeleteById(w http.ResponseWriter, r *http.Request) {
	slog.Info("Delete request DeleteById")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Test"))
}
