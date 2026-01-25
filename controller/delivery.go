package controller

import (
	"erp-2c/lib/response"
	"erp-2c/lib/sl"
	"erp-2c/lib/types"
	"erp-2c/model"
	"erp-2c/service/use_cases"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type DeliveryController struct {
	services *use_cases.Manager
	validate *validator.Validate
}

func NewDeliveryController(
	services *use_cases.Manager, validate *validator.Validate) *DeliveryController {
	return &DeliveryController{services: services, validate: validate}
}

func (d *DeliveryController) Save(w http.ResponseWriter, r *http.Request) {
	const op = "control.delivery.Save"

	var requestBody model.DeliveryToSave
	if err := render.DecodeJSON(r.Body, &requestBody); err != nil {
		slog.Error("failed pare request body", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid request body").SendResponse(w, r)
		return
	}

	if err := d.validate.Struct(requestBody); err != nil {
		slog.Error("failed validate request fields", sl.ErrWithOP(err, op))
		response.ValidationError(err).SendResponse(w, r)
		return
	}
	saved, err := d.services.DeliveryService.Save(requestBody)
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.Created(saved)

}
func (d *DeliveryController) GetById(w http.ResponseWriter, r *http.Request)    {}
func (d *DeliveryController) GetAll(w http.ResponseWriter, r *http.Request)     {}
func (d *DeliveryController) UpdateById(w http.ResponseWriter, r *http.Request) {}
func (d *DeliveryController) DeleteById(w http.ResponseWriter, r *http.Request) {}
