package controller

import (
	"erp-2c/service/use_cases"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type DeliveryController struct {
	services *use_cases.Manager
	validate *validator.Validate
}

func NewDeliveryController(
	services *use_cases.Manager, validate *validator.Validate) *DeliveryController {
	return &DeliveryController{services: services, validate: validate}
}

func (d *DeliveryController) Save(w http.ResponseWriter, r *http.Request)       {}
func (d *DeliveryController) GetById(w http.ResponseWriter, r *http.Request)    {}
func (d *DeliveryController) GetAll(w http.ResponseWriter, r *http.Request)     {}
func (d *DeliveryController) UpdateById(w http.ResponseWriter, r *http.Request) {}
