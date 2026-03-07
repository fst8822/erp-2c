package controller

import (
	"erp-2c/lib/response"
	"erp-2c/lib/sl"
	"erp-2c/lib/types"
	"erp-2c/model"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type deliveryServiceInt interface {
	Save(delivery model.DeliveryItemsDomain) (*model.DeliveryItemsDomain, error)
	GetById(deliveryId int64) (*model.DeliveryItemsDomain, error)
	GetAll() (*model.DeliveryItemListDomain, error)
	GetByStatus(status model.DeliveryStatus) (*model.DeliveryItemListDomain, error)
	UpdateById(deliveryId int64, update model.UpdateStatus) error
	DeleteById(deliveryId int64) error
}

type DeliveryController struct {
	deliveryService deliveryServiceInt
	validate        *validator.Validate
}

func NewDeliveryController(
	deliveryService deliveryServiceInt, validate *validator.Validate) *DeliveryController {
	return &DeliveryController{deliveryService: deliveryService, validate: validate}
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
	ID := r.Context().Value("userIdKey")
	userID, ok := ID.(int64)
	if !ok {
		sLogger.Error("failed failed get user from context")
		response.InternalServerError().SendResponse(w, r)
		return
	}
	DeliveryItems := requestBody.MapToDomain(userID)
	saved, err := d.deliveryService.Save(DeliveryItems)
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

	found, err := d.deliveryService.GetById(id)
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.OK(found).SendResponse(w, r)
}

func (d *DeliveryController) GetAll(w http.ResponseWriter, r *http.Request) {
	deliveryDomains, err := d.deliveryService.GetAll()
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.OK(deliveryDomains).SendResponse(w, r)
}

func (d *DeliveryController) UpdateById(w http.ResponseWriter, r *http.Request) {}
func (d *DeliveryController) DeleteById(w http.ResponseWriter, r *http.Request) {}
