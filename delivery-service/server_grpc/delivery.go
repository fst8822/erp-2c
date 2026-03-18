package server_grpc

import (
	"context"
	"delivery-service/lib/sl"
	"delivery-service/model"
	servergrpc "delivery-service/server_grpc/proto/v1/delivery"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type deliveryServiceInt interface {
	Save(delivery model.DeliveryItemsDomain) (*model.DeliveryItemsDomain, error)
	GetById(deliveryId int64) (*model.DeliveryItemsDomain, error)
	GetByStatus(status model.DeliveryStatus) (*model.DeliveryItemListDomain, error)
	GetAll() (*model.DeliveryItemListDomain, error)
	UpdateById(deliveryId int64, update model.UpdateStatus) error
	DeleteById(deliveryId int64) error
}

type DeliveryGRPCServer struct {
	servergrpc.UnimplementedDeliveryServiceServer
	deliveryServiceInt
}

func NewDeliveryGRPCServer(deliveryServiceInt deliveryServiceInt) *DeliveryGRPCServer {
	return &DeliveryGRPCServer{deliveryServiceInt: deliveryServiceInt}
}

func (d *DeliveryGRPCServer) Save(
	ctx context.Context, req *servergrpc.DeliveryItemsDomain) (*servergrpc.DeliveryItemsDomain, error) {
	const op = " delivery-service.server_grpc.service.delivery.Save"

	deliveryStatus := model.DeliveryStatus(req.Delivery.GetStatus().Status)
	if err := deliveryStatus.IsValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	itemsToSave := make([]model.ItemDomain, 0, len(req.Items))
	for i, item := range req.Items {
		itemsToSave[i] = model.ItemDomain{
			ProductID:  item.GetProductId(),
			ItemPrice:  item.GetItemPrice(),
			Quantity:   item.GetQuantity(),
			ItemAmount: item.GetItemAmount(),
		}
	}
	deliveryToSave := model.DeliveryItemsDomain{
		DeliverDomain: model.DeliverDomain{
			Recipient:     req.Delivery.GetRecipient(),
			Address:       req.Delivery.GetAddress(),
			Status:        deliveryStatus,
			CreatedAt:     req.Delivery.GetCreatedAt().AsTime(),
			UserID:        req.Delivery.GetUserId(),
			DeliverAmount: req.Delivery.GetDeliverAmount(),
		},
		Items: itemsToSave,
	}
	saved, err := d.deliveryServiceInt.Save(deliveryToSave)
	if err != nil {
		slog.Error("failed save delivery", sl.ErrWithOP(err, op))
		return nil, status.Error(codes.Internal, err.Error())
	}

	deliverDomain := &servergrpc.DeliverDomain{
		Id:            saved.ID,
		Recipient:     saved.Recipient,
		Address:       saved.Address,
		Status:        &servergrpc.DeliveryStatus{Status: string(saved.Status)},
		CreatedAt:     timestamppb.New(saved.CreatedAt),
		UserId:        saved.UserID,
		DeliverAmount: saved.DeliverAmount,
	}
	itemsToResponse := make([]*servergrpc.ItemDomain, 0, len(saved.Items))
	for i, item := range saved.Items {
		itemsToResponse[i] = &servergrpc.ItemDomain{
			DeliveryId: item.DeliveryID,
			ProductId:  item.ProductID,
			ItemPrice:  item.ItemPrice,
			Quantity:   item.Quantity,
			ItemAmount: item.ItemAmount,
		}
	}
	deliveryItemsToResponse := &servergrpc.DeliveryItemsDomain{
		Delivery: deliverDomain,
		Items:    itemsToResponse,
	}
	return deliveryItemsToResponse, nil
}

func (d *DeliveryGRPCServer) GetById(
	ctx context.Context, req *servergrpc.RequestDeliveryID) (*servergrpc.DeliveryItemsDomain, error) {
	const op = " delivery-service.server_grpc.service.delivery.GetById"
	found, err := d.deliveryServiceInt.GetById(req.GetDeliveryId())
	if err != nil {
		slog.Error("failed get delivery",
			slog.Int64("delivery id", req.GetDeliveryId()), sl.ErrWithOP(err, op))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	deliverDomain := &servergrpc.DeliverDomain{
		Id:            found.ID,
		Recipient:     found.Recipient,
		Address:       found.Address,
		Status:        &servergrpc.DeliveryStatus{Status: string(found.Status)},
		CreatedAt:     timestamppb.New(found.CreatedAt),
		UserId:        found.UserID,
		DeliverAmount: found.DeliverAmount,
	}
	itemsToResponse := make([]*servergrpc.ItemDomain, 0, len(found.Items))
	for i, item := range found.Items {
		itemsToResponse[i] = &servergrpc.ItemDomain{
			DeliveryId: item.DeliveryID,
			ProductId:  item.ProductID,
			ItemPrice:  item.ItemPrice,
			Quantity:   item.Quantity,
			ItemAmount: item.ItemAmount,
		}
	}
	deliveryItemsToResponse := &servergrpc.DeliveryItemsDomain{
		Delivery: deliverDomain,
		Items:    itemsToResponse,
	}
	return deliveryItemsToResponse, nil
}

func (d *DeliveryGRPCServer) GetByStatus(
	req *servergrpc.DeliveryStatus, stream grpc.ServerStreamingServer[servergrpc.DeliveryItemsDomain]) error {
	const op = " delivery-service.server_grpc.service.delivery.GetByStatus"

	deliveryStatus := model.DeliveryStatus(req.GetStatus())
	if err := deliveryStatus.IsValid(); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	founds, err := d.deliveryServiceInt.GetByStatus(deliveryStatus)
	if err != nil {
		slog.Error("failed get delivery",
			slog.String("delivery status", string(deliveryStatus)), sl.ErrWithOP(err, op))
		return status.Error(codes.NotFound, err.Error())
	}

	deliveryItemsToResponse := make([]*servergrpc.DeliveryItemsDomain, 0, len(founds.DeliveryItemsDomain))
	for i, it := range founds.DeliveryItemsDomain {
		deliverToResponse := &servergrpc.DeliverDomain{
			Id:            it.ID,
			Recipient:     it.Recipient,
			Address:       it.Address,
			Status:        &servergrpc.DeliveryStatus{Status: string(it.Status)},
			CreatedAt:     timestamppb.New(it.CreatedAt),
			UserId:        it.UserID,
			DeliverAmount: it.DeliverAmount,
		}
		itemsToResponse := make([]*servergrpc.ItemDomain, 0, len(it.Items))
		for i, item := range it.Items {
			itemsToResponse[i] = &servergrpc.ItemDomain{
				DeliveryId: item.DeliveryID,
				ProductId:  item.ProductID,
				ItemPrice:  item.ItemPrice,
				Quantity:   item.Quantity,
				ItemAmount: item.ItemAmount,
			}
		}
		deliveryItemsToResponse[i] = &servergrpc.DeliveryItemsDomain{
			Delivery: deliverToResponse,
			Items:    itemsToResponse,
		}
	}
	for _, it := range deliveryItemsToResponse {
		err := stream.Send(it)
		if err != nil {
			slog.Error("Error sending product throw stream grpc server", sl.ErrWithOP(err, op))
		}
	}
	return nil
}

func (d *DeliveryGRPCServer) GetAll(
	empty *emptypb.Empty, stream grpc.ServerStreamingServer[servergrpc.DeliveryItemsDomain]) error {
	const op = " delivery-service.server_grpc.service.delivery.GetAll"

	founds, err := d.deliveryServiceInt.GetAll()
	if err != nil {
		slog.Error("failed get all delivery", sl.ErrWithOP(err, op))
		return status.Error(codes.Internal, err.Error())
	}

	deliveryItemsToResponse := make([]*servergrpc.DeliveryItemsDomain, 0, len(founds.DeliveryItemsDomain))
	for i, it := range founds.DeliveryItemsDomain {
		deliverToResponse := &servergrpc.DeliverDomain{
			Id:            it.ID,
			Recipient:     it.Recipient,
			Address:       it.Address,
			Status:        &servergrpc.DeliveryStatus{Status: string(it.Status)},
			CreatedAt:     timestamppb.New(it.CreatedAt),
			UserId:        it.UserID,
			DeliverAmount: it.DeliverAmount,
		}
		itemsToResponse := make([]*servergrpc.ItemDomain, 0, len(it.Items))
		for i, item := range it.Items {
			itemsToResponse[i] = &servergrpc.ItemDomain{
				DeliveryId: item.DeliveryID,
				ProductId:  item.ProductID,
				ItemPrice:  item.ItemPrice,
				Quantity:   item.Quantity,
				ItemAmount: item.ItemAmount,
			}
		}
		deliveryItemsToResponse[i] = &servergrpc.DeliveryItemsDomain{
			Delivery: deliverToResponse,
			Items:    itemsToResponse,
		}
	}
	for _, it := range deliveryItemsToResponse {
		err := stream.Send(it)
		if err != nil {
			slog.Error("Error sending product throw stream grpc server", sl.ErrWithOP(err, op))
		}
	}
	return nil
}

func (d *DeliveryGRPCServer) UpdateById(
	ctx context.Context, req *servergrpc.UpdateStatus) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateById not implemented")
}

// DeleteById todo в разработке
func (d *DeliveryGRPCServer) DeleteById(
	ctx context.Context, req *servergrpc.RequestDeliveryID) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteById not implemented")
}
