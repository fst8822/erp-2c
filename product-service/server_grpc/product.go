package server_grpc

import (
	"context"
	"product-service/lib/sl"
	"product-service/lib/types"
	"product-service/model"
	servergrpc "product-service/server_grpc/proto/v1/product"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type productServiceInt interface {
	Save(productToSave model.ProductToSave) (*model.ProductDomain, error)
	GetById(productId int64) (*model.ProductDomain, error)
	GetByName(productName string) (*model.ProductDomain, error)
	GetAll() ([]model.ProductDomain, error)
	GetExistIds(productIds []int64) ([]int64, error)
	UpdateById(productId int64, productToUpdate model.ProductUpdate) error
	DeleteById(productId int64) error
}

type ProductGRPCServer struct {
	servergrpc.UnimplementedProductServiceServer
	productServer productServiceInt
}

func NewProductGRPCServer(productServer productServiceInt) *ProductGRPCServer {
	return &ProductGRPCServer{productServer: productServer}
}

// Save todo need check ctx
func (p *ProductGRPCServer) Save(
	ctx context.Context, req *servergrpc.ProductToSave) (*servergrpc.ProductDomain, error) {
	const op = " delivery-service.server_grpc.service.product.Save"

	productToSave := model.ProductToSave{
		ProductName:  req.GetProductName(),
		ProductGroup: req.GetProductGroup(),
		Image:        req.GetImage(),
		Stock:        req.GetStock(),
		Price:        req.GetPrice(),
	}
	saved, err := p.productServer.Save(productToSave)
	if err != nil {
		slog.Error("failed save product",
			slog.String("product name", productToSave.ProductName), sl.ErrWithOP(err, op))
		return nil, types.HandleError(err)
	}

	return &servergrpc.ProductDomain{
		Id:           saved.Id,
		ProductName:  saved.ProductName,
		ProductGroup: saved.ProductGroup,
		Image:        saved.Image,
		Stock:        saved.Stock,
		Price:        saved.Price,
	}, nil
}

func (p *ProductGRPCServer) GetById(
	ctx context.Context, req *servergrpc.RequestProductID) (*servergrpc.ProductDomain, error) {
	const op = " delivery-service.server_grpc.service.product.GetById"

	found, err := p.productServer.GetById(req.Id)
	if err != nil {
		slog.Error("failed to find product", sl.ErrWithOP(err, op))
		return nil, types.HandleError(err)
	}
	return &servergrpc.ProductDomain{
		Id:           found.Id,
		ProductName:  found.ProductName,
		ProductGroup: found.ProductGroup,
		Image:        found.Image,
		Stock:        found.Stock,
		Price:        found.Price,
	}, nil
}

func (p *ProductGRPCServer) GetByName(
	ctx context.Context, req *servergrpc.RequestProductName) (*servergrpc.ProductDomain, error) {
	const op = " delivery-service.server_grpc.service.product.GetByName"

	found, err := p.productServer.GetByName(req.Name)
	if err != nil {
		slog.Error("failed to find product", sl.ErrWithOP(err, op))
		return nil, types.HandleError(err)
	}
	return &servergrpc.ProductDomain{
		Id:           found.Id,
		ProductName:  found.ProductName,
		ProductGroup: found.ProductGroup,
		Image:        found.Image,
		Stock:        found.Stock,
		Price:        found.Price,
	}, nil
}

func (p *ProductGRPCServer) GetAll(
	empty *emptypb.Empty, stream grpc.ServerStreamingServer[servergrpc.ProductDomain]) error {
	const op = " delivery-service.server_grpc.service.product.GetAll"

	founds, err := p.productServer.GetAll()
	if err != nil {
		slog.Error("Error get all product", sl.ErrWithOP(err, op))
		return status.Error(codes.Internal, err.Error())
	}

	for _, found := range founds {
		err := stream.Send(&servergrpc.ProductDomain{
			Id:           found.Id,
			ProductName:  found.ProductName,
			ProductGroup: found.ProductGroup,
			Image:        found.Image,
			Stock:        found.Stock,
			Price:        found.Price,
		})
		if err != nil {
			slog.Error("Error sending product throw stream grpc server", sl.ErrWithOP(err, op))
		}
	}
	return nil
}

func (p *ProductGRPCServer) GetExistIds(
	ctx context.Context, productIDs *servergrpc.ListProductIDs) (*servergrpc.ListProductIDs, error) {
	const op = " delivery-service.server_grpc.service.product.GetByName"

	ids, err := p.productServer.GetExistIds(productIDs.GetProductId())
	if err != nil {
		slog.Error("failed to find product ids", sl.ErrWithOP(err, op))
		return nil, status.Error(codes.Internal, "failed to find product ids")
	}
	return &servergrpc.ListProductIDs{ProductId: ids}, nil
}

func (p *ProductGRPCServer) UpdateById(
	ctx context.Context, req *servergrpc.UpdateProductByRequestID) (*emptypb.Empty, error) {
	const op = " delivery-service.server_grpc.service.product.UpdateById"

	if req.ProductToUpdate == nil {
		slog.Info("request body to update is empty")
		return nil, status.Error(codes.InvalidArgument, "request body is empty")
	}
	productToUpdate := model.ProductUpdate{
		ProductName:  &req.ProductToUpdate.ProductName,
		ProductGroup: &req.ProductToUpdate.ProductGroup,
		Image:        nil,
		Stock:        &req.ProductToUpdate.Stock,
		Price:        &req.ProductToUpdate.Price,
	}
	err := p.productServer.UpdateById(req.GetProductId(), productToUpdate)
	if err != nil {
		slog.Error("failed to update product", sl.ErrWithOP(err, op))
		return nil, types.HandleError(err)
	}
	return nil, nil
}

func (p *ProductGRPCServer) DeleteById(
	ctx context.Context, req *servergrpc.RequestProductID) (*emptypb.Empty, error) {
	const op = " delivery-service.server_grpc.service.product.DeleteById"

	err := p.productServer.DeleteById(req.GetId())
	if err != nil {
		slog.Error("failed to DeleteById product", sl.ErrWithOP(err, op))
		return nil, types.HandleError(err)
	}
	return nil, nil
}
