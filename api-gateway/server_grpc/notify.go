package server_grpc

import (
	"api-gateway/lib/sl"
	"api-gateway/model"
	pb "api-gateway/server_grpc/proto/v1/notify"
	"io"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type notifyServiceInt interface {
	SendNotify(notification model.Notification)
}

type NotifyGRPCServer struct {
	pb.UnimplementedNotifyServiceServer
	notify notifyServiceInt
}

func NewNotifyGRPCServer(notify notifyServiceInt) *NotifyGRPCServer {
	return &NotifyGRPCServer{notify: notify}
}

func (n *NotifyGRPCServer) SendNotify(stream grpc.ClientStreamingServer[pb.Notification, emptypb.Empty]) error {
	const op = "server_grpc.notify.SendNotify"
	sLogger := slog.With("op", op)

	var notifications []model.Notification
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			sLogger.Info("stream close", sl.Err(err))
			break
			//todo stream.SendAndClose()
		}
		if err != nil {
			sLogger.Error("failed get response from stream", sl.Err(err))
			return status.Error(codes.Internal, err.Error())
		}

		notification := model.Notification{
			DeliveryId: resp.DeliveryId,
			Status:     model.DeliveryStatus(resp.Status),
			CreatedAt:  resp.CreatedAt.AsTime(),
			UserID:     resp.UserId,
		}
		notifications = append(notifications, notification)
	}
	go func(notifications []model.Notification) {
		for _, notification := range notifications {
			n.notify.SendNotify(notification)
		}
	}(notifications)
	return nil
}
