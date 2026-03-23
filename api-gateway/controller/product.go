package controller

import (
	"api-gateway/lib/response"
	"api-gateway/lib/sl"
	"api-gateway/lib/types"
	"api-gateway/model"
	productgrpc "api-gateway/server_grpc/proto/v1/product"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type ProductController struct {
	productClientGRPC productgrpc.ProductServiceClient
	validate          *validator.Validate
}

func NewProductController(
	productClientGRPC productgrpc.ProductServiceClient,
	validate *validator.Validate) *ProductController {
	return &ProductController{
		productClientGRPC: productClientGRPC,
		validate:          validate,
	}
}

func (p *ProductController) Save(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.Save"

	var productToSave model.ProductToSave
	if err := render.DecodeJSON(r.Body, &productToSave); err != nil {
		slog.Error("failed parse request body", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid request body").SendResponse(w, r)
		return
	}

	if err := p.validate.Struct(productToSave); err != nil {
		slog.Error("failed validate request fields", sl.ErrWithOP(err, op))
		response.ValidationError(err).SendResponse(w, r)
		return
	}
	productToSend := &productgrpc.ProductToSave{
		ProductName:  productToSave.ProductName,
		ProductGroup: productToSave.ProductGroup,
		Image:        productToSave.Image,
		Stock:        productToSave.Stock,
		Price:        productToSave.Price,
	}
	saved, err := p.productClientGRPC.Save(r.Context(), productToSend)
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	resp := model.ProductResponse{
		Id:           saved.Id,
		ProductName:  saved.ProductName,
		ProductGroup: saved.ProductGroup,
		Image:        saved.Image,
		Stock:        saved.Stock,
		Price:        saved.Price,
	}
	response.Created(resp).SendResponse(w, r)
}

func (p *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.GetAll"

	stream, err := p.productClientGRPC.GetAll(r.Context(), nil)
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	var products []model.ProductResponse
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			slog.Info("stream close", sl.ErrWithOP(err, op))
			break
		}
		if err != nil {
			slog.Error("failed get response from stream", sl.ErrWithOP(err, op))
			types.HandleError(err).SendResponse(w, r)
			return
		}

		products = append(products, model.ProductResponse{
			Id:           resp.GetId(),
			ProductName:  resp.ProductName,
			ProductGroup: resp.ProductGroup,
			Image:        resp.Image,
			Stock:        resp.Stock,
			Price:        resp.Price,
		})
	}
	response.OK(products).SendResponse(w, r)
}

func (p *ProductController) GetById(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.GetWithItemsById"

	param := chi.URLParam(r, "id")
	productId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		slog.Error("failed convert str to int64", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid path variable").SendResponse(w, r)
		return
	}

	found, err := p.productClientGRPC.GetById(r.Context(), &productgrpc.RequestProductID{Id: productId})
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	productResponse := model.ProductResponse{
		Id:           found.Id,
		ProductName:  found.ProductName,
		ProductGroup: found.ProductGroup,
		Image:        found.Image,
		Stock:        found.Stock,
		Price:        found.Price,
	}
	response.OK(productResponse).SendResponse(w, r)
}

func (p *ProductController) GetByName(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.GetByName"

	productName := chi.URLParam(r, "name")
	if strings.TrimSpace(productName) == "" {
		slog.Error("Path variable is empty", slog.StringValue(productName), slog.StringValue(op))
		response.BadRequest("Invalid path variable").SendResponse(w, r)
		return
	}

	found, err := p.productClientGRPC.GetByName(r.Context(), &productgrpc.RequestProductName{Name: productName})
	if err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	productResponse := model.ProductResponse{
		Id:           found.Id,
		ProductName:  found.ProductName,
		ProductGroup: found.ProductGroup,
		Image:        found.Image,
		Stock:        found.Stock,
		Price:        found.Price,
	}
	response.OK(productResponse).SendResponse(w, r)
}

func (p *ProductController) UpdateById(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.UpdateById"

	Param := chi.URLParam(r, "id")
	productId, err := strconv.ParseInt(Param, 10, 64)
	if err != nil {
		slog.Error("failed convert str to int64", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid path variable").SendResponse(w, r)
		return
	}

	var productToUpdate model.ProductUpdate

	if err := render.DecodeJSON(r.Body, &productToUpdate); err != nil {
		slog.Error("failed parse request body", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid request body").SendResponse(w, r)
		return
	}

	requestUpdate := &productgrpc.UpdateProductByRequestID{
		ProductId: productId,
		ProductToUpdate: &productgrpc.ProductUpdate{
			ProductName:  *productToUpdate.ProductName,
			ProductGroup: *productToUpdate.ProductGroup,
			Image:        nil,
			Stock:        *productToUpdate.Stock,
			Price:        *productToUpdate.Price,
		},
	}
	if _, err := p.productClientGRPC.UpdateById(r.Context(), requestUpdate); err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.OK(nil).SendResponse(w, r)
}

func (p *ProductController) DeleteById(w http.ResponseWriter, r *http.Request) {
	const op = "control.product.DeleteById"

	param := chi.URLParam(r, "id")
	productId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		slog.Error("failed convert str to int64", sl.ErrWithOP(err, op))
		response.BadRequest("Invalid path variable").SendResponse(w, r)
		return
	}

	if _, err := p.productClientGRPC.DeleteById(r.Context(), &productgrpc.RequestProductID{Id: productId}); err != nil {
		types.HandleError(err).SendResponse(w, r)
		return
	}
	response.NoContent()
}
