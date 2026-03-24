package use_cases

import (
	"context"
	"errors"
	"io"
	"notification-service/lib/observability/app_metrics"
	"notification-service/lib/sl"
	"notification-service/model"
	pb "notification-service/server_grpc/proto/v1/notify"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NotifyService struct {
	mu      sync.RWMutex
	clients map[int64]*model.ClientWS
	done    chan struct{}
	wg      sync.WaitGroup
	pb.UnimplementedNotifyServiceServer
}

func NewNotifyService() *NotifyService {
	return &NotifyService{
		clients: make(map[int64]*model.ClientWS),
		done:    make(chan struct{}),
	}
}

func (n *NotifyService) Subscribe(ctx context.Context, client *model.ClientWS) {
	const OP = "services.use_cases.notify.NotifyService.Subscribe"
	log := slog.With("OP", OP, "UserID", client.UserID)

	start := time.Now()
	ctxDone, cancel := context.WithCancel(ctx)
	go n.checkClientConn(cancel, client)

	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		defer func() {
			n.RemoveClient(client)
			defer app_metrics.WebsocketActiveConn.Dec()
			defer app_metrics.WebsocketSessionDuration.Observe(time.Since(start).Seconds())
		}()

		for {
			log.Info("Start WS worker")
			select {
			case <-ctx.Done():
				log.Info("Client disconnected, context is don, exit from subscribe")
				client.Conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Client disconnected"))
				return

			case <-n.done:
				log.Info("Server is shutdown, closing connection")
				client.Conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server is shutdown"))
				return

			case notification, ok := <-client.Cn:
				if !ok {
					log.Info("Notification channel is close, exit from subscribe")
					return
				}
				if err := client.Conn.WriteJSON(notification); err != nil {
					log.Error("websocket is close, exit from subscribe", sl.Err(err))
					return
				}
				slog.Info("Send message successful", slog.Any("notification", notification))

			case <-ctxDone.Done():
				log.Error("client connection is last, exit from subscribe")
				return
			}
		}
	}()
}

func (n *NotifyService) SendNotify(stream grpc.ClientStreamingServer[pb.Notification, emptypb.Empty]) error {
	const OP = "services.use_cases.notify.NotifyService.SendNotify"
	log := slog.With("OP", OP)

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Info("stream close", sl.Err(err))
			return stream.SendAndClose(&emptypb.Empty{})
		}
		if err != nil {
			log.Error("failed get response from stream", sl.Err(err))
			return status.Error(codes.Internal, err.Error())
		}

		notification := model.Notification{
			DeliveryId: resp.DeliveryId,
			Status:     resp.Status,
			CreatedAt:  resp.CreatedAt.AsTime(),
			UserID:     resp.UserId,
		}
		n.mu.RLock()
		clientWS, ok := n.clients[notification.UserID]
		n.mu.RUnlock()
		if !ok {
			slog.Error("no active websocket for user", slog.Int64("userID", notification.UserID))
			continue
		}
		log.Info("Begin write notification in notifications")
		select {
		case clientWS.Cn <- notification:
		default:
			log.Info("client channel full, dropping notification", slog.Any("notification", notification))
		}
		log.Info("End write notification in notifications channel")
	}
}

func (n *NotifyService) AddClient(client *model.ClientWS) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if old, exist := n.clients[client.UserID]; exist {
		old.Conn.Close()
		close(old.Cn)
		delete(n.clients, old.UserID)
	}
	n.clients[client.UserID] = client
}

func (n *NotifyService) RemoveClient(client *model.ClientWS) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if found, ok := n.clients[client.UserID]; ok && found == client {
		found.Conn.Close()
		close(found.Cn)
		delete(n.clients, found.UserID)
	}
}

func (n *NotifyService) Shutdown() {
	const OP = "services.use_cases.notify.NotifyService.Shutdown"
	log := slog.With("OP", OP)

	log.Info("Start shutdown is client ws")
	close(n.done)
	n.wg.Wait()
	log.Info("Shutdown client ws is done")
}

func (n *NotifyService) checkClientConn(cancel context.CancelFunc, client *model.ClientWS) {
	const OP = "services.use_cases.notify.NotifyService.checkClientConn"

	_, _, err := client.Conn.ReadMessage()
	defer cancel()
	if err != nil {
		slog.Info("Connection WS is close",
			slog.Int64("clientID", client.UserID),
			slog.String("OP", OP),
		)
	}
}
