package controller

import (
	"api-gateway/lib/response"
	"api-gateway/lib/sl"
	"api-gateway/lib/types"
	"api-gateway/model"
	deliverygrpc "api-gateway/server_grpc/proto/v1/delivery"
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DeliveryController struct {
	deliveryClientGRPC deliverygrpc.DeliveryServiceClient
	validate           *validator.Validate
}

func NewDeliveryController(
	deliveryClientGRPC deliverygrpc.DeliveryServiceClient, validate *validator.Validate) *DeliveryController {
	return &DeliveryController{deliveryClientGRPC: deliveryClientGRPC, validate: validate}
}

func (d *DeliveryController) Save(w http.ResponseWriter, r *http.Request) {
	const op = "control.delivery.Save"
	sLogger := slog.With("OP", op)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
	ID := r.Context().Value(model.UserIdKey)
	userID, ok := ID.(int64)
	if !ok {
		sLogger.Error("failed get user from context")
		response.InternalServerError().SendResponse(w, r)
		return
	}

	itemDomains := make([]*deliverygrpc.ItemDomain, 0, len(requestBody.Items))
	for i, item := range requestBody.Items {
		itemDomains[i] = &deliverygrpc.ItemDomain{
			ProductId:  item.ProductId,
			ItemPrice:  item.ItemPrice,
			Quantity:   item.Quantity,
			ItemAmount: 0,
		}
	}
	saved, err := d.deliveryClientGRPC.Save(ctx, &deliverygrpc.DeliveryItemsDomain{
		Delivery: &deliverygrpc.DeliverDomain{
			Recipient: requestBody.Recipient,
			Address:   requestBody.Address,
			CreatedAt: timestamppb.Now(),
			UserId:    userID,
		},
		Items: itemDomains,
	})
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}

	itemsResponse := make([]model.ItemResponse, 0, len(saved.Items))
	for i, item := range saved.Items {
		itemsResponse[i] = model.ItemResponse{
			DeliveryID: item.DeliveryId,
			ProductID:  item.ProductId,
			ItemPrice:  item.ItemPrice,
			Quantity:   item.Quantity,
			ItemAmount: item.ItemAmount,
		}
	}
	deliveryResponse := model.DeliveryResponse{
		ID:            saved.Delivery.Id,
		Recipient:     saved.Delivery.Recipient,
		Address:       saved.Delivery.Address,
		Status:        model.DeliveryStatus(saved.Delivery.Status.Status),
		CreatedAt:     saved.Delivery.CreatedAt.AsTime(),
		UserID:        saved.Delivery.UserId,
		DeliverAmount: saved.Delivery.DeliverAmount,
		Items:         itemsResponse,
	}
	response.Created(deliveryResponse).SendResponse(w, r)
}

func (d *DeliveryController) GetById(w http.ResponseWriter, r *http.Request) {
	const op = "control.delivery.GetWithItemsById"
	sLogger := slog.With("op", op)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sLogger.Error("failed convert str to int64", sl.Err(err))
		response.BadRequest("Invalid path variable").SendResponse(w, r)
		return
	}

	found, err := d.deliveryClientGRPC.GetById(ctx, &deliverygrpc.RequestDeliveryID{DeliveryId: id})
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	itemsResponse := make([]model.ItemResponse, 0, len(found.Items))
	for i, item := range found.Items {
		itemsResponse[i] = model.ItemResponse{
			DeliveryID: item.DeliveryId,
			ProductID:  item.ProductId,
			ItemPrice:  item.ItemPrice,
			Quantity:   item.Quantity,
			ItemAmount: item.ItemAmount,
		}
	}
	deliveryResponse := model.DeliveryResponse{
		ID:            found.Delivery.Id,
		Recipient:     found.Delivery.Recipient,
		Address:       found.Delivery.Address,
		Status:        model.DeliveryStatus(found.Delivery.Status.Status),
		CreatedAt:     found.Delivery.CreatedAt.AsTime(),
		UserID:        found.Delivery.UserId,
		DeliverAmount: found.Delivery.DeliverAmount,
		Items:         itemsResponse,
	}
	response.OK(deliveryResponse).SendResponse(w, r)
}

func (d *DeliveryController) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "control.delivery.GetWithItemsById"
	sLogger := slog.With("op", op)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := d.deliveryClientGRPC.GetAll(ctx, nil)
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	var listResponse []model.DeliveryResponse
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			sLogger.Info("stream close", sl.Err(err))
			break
		}
		if err != nil {
			sLogger.Error("failed get response from stream", sl.Err(err))
			types.HandleError(err).SendResponse(w, r)
			return
		}

		items := make([]model.ItemResponse, 0, len(resp.Items))
		for i, item := range resp.Items {
			items[i] = model.ItemResponse{
				DeliveryID: item.DeliveryId,
				ProductID:  item.ProductId,
				ItemPrice:  item.ItemPrice,
				Quantity:   item.Quantity,
				ItemAmount: item.ItemAmount,
			}
		}
		deliveryResp := model.DeliveryResponse{
			ID:            resp.Delivery.Id,
			Recipient:     resp.Delivery.Recipient,
			Address:       resp.Delivery.Address,
			Status:        model.DeliveryStatus(resp.Delivery.Status.Status),
			CreatedAt:     resp.Delivery.CreatedAt.AsTime(),
			UserID:        resp.Delivery.UserId,
			DeliverAmount: resp.Delivery.DeliverAmount,
			Items:         items,
		}
		listResponse = append(listResponse, deliveryResp)
	}
	response.OK(listResponse).SendResponse(w, r)
}

func (d *DeliveryController) UpdateById(w http.ResponseWriter, r *http.Request) {}
func (d *DeliveryController) DeleteById(w http.ResponseWriter, r *http.Request) {}
