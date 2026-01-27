package controller

import (
	"erp-2c/lib/response"
	"erp-2c/lib/sl"
	"erp-2c/lib/types"
	"erp-2c/model"
	"erp-2c/service/use_cases"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
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
	sLogger := slog.With("OP", op)

	var requestBody model.DeliveryToSave
	if err := render.DecodeJSON(r.Body, &requestBody); err != nil {
		sLogger.Error("failed pare request body", sl.Err(err))
		response.BadRequest("Invalid request body").SendResponse(w, r)
		return
	}

	if err := d.validate.Struct(requestBody); err != nil {
		sLogger.Error("failed validate request fields", sl.Err(err))
		response.ValidationError(err).SendResponse(w, r)
		return
	}
	DeliveryItems := requestBody.MapToDomain()
	saved, err := d.services.DeliveryService.Save(DeliveryItems)
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.Created(saved).SendResponse(w, r)
}

func (d *DeliveryController) GetById(w http.ResponseWriter, r *http.Request) {
	const op = "control.delivery.GetWithItemsById"
	sLogger := slog.With("op", op)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sLogger.Error("failed convert str to int64", sl.Err(err))
		response.BadRequest("Invalid path variable").SendResponse(w, r)
		return
	}

	found, err := d.services.DeliveryService.GetById(id)
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.OK(found).SendResponse(w, r)
}

func (d *DeliveryController) GetAll(w http.ResponseWriter, r *http.Request) {
	deliveryDomains, err := d.services.DeliveryService.GetAll()
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.OK(deliveryDomains).SendResponse(w, r)
}

func (d *DeliveryController) UpdateById(w http.ResponseWriter, r *http.Request) {}
func (d *DeliveryController) DeleteById(w http.ResponseWriter, r *http.Request) {}
