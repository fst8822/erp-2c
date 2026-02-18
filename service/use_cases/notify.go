package use_cases

import (
	"context"
	"erp-2c/lib/sl"
	"erp-2c/model"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

const PingInterval = 10 * time.Second

type NotifyService struct {
	mu      sync.RWMutex
	clients map[int64]*model.ClientWS
}

func NewNotifyService() *NotifyService {
	return &NotifyService{
		clients: make(map[int64]*model.ClientWS),
	}
}

func (n *NotifyService) Subscribe(ctx context.Context, client *model.ClientWS) {
	const OP = "services.use_cases.notify.NewNotifyService.Subscribe"
	log := slog.With("OP", OP, "UserID", client.UserID)

	defer func() {
		n.RemoveClient(client)
	}()

	ticket := time.NewTicker(PingInterval)
	defer ticket.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("Context is don, exit from subscribe")
			client.Conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "context is done"))
			return
		case notification, ok := <-client.Cn:
			if !ok {
				log.Info("Channel is close, exit from subscribe")
				return
			}
			if err := client.Conn.WriteJSON(notification); err != nil {
				log.Error("websocket is close, exit from subscribe", sl.Err(err))
				return
			}
			slog.Info("Send message successful", slog.Any("notification", notification))
		case <-ticket.C:
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error("websocket is close, exit from subscribe", sl.Err(err))
				return
			}
		}
	}
}

func (n *NotifyService) SendNotify(notification model.Notification) {
	const OP = "services.use_cases.notify.NewNotifyService.SendNotify"
	log := slog.With("OP", OP, "UserID", notification.UserID)

	n.mu.RLock()
	clientWS, ok := n.clients[notification.UserID]
	n.mu.RUnlock()

	if !ok {
		slog.Error("Websocket is already close")
		return
	}

	log.Info("Begin write notification in notifications")
	select {
	case clientWS.Cn <- notification:
	default:
		log.Info("client channel full, dropping notification", slog.Any("notification", notification))
	}
	log.Info("End write notification in notifications channel")
}

func (n *NotifyService) AddClient(client *model.ClientWS) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.clients[client.UserID] = client
}

func (n *NotifyService) RemoveClient(client *model.ClientWS) {
	n.mu.Lock()
	defer n.mu.Unlock()

	found, ok := n.clients[client.UserID]
	if !ok || found.UUID != client.UUID {
		client.Conn.Close()
		close(client.Cn)
		return
	}
	found.Conn.Close()
	delete(n.clients, found.UserID)
	close(found.Cn)
}

func (n *NotifyService) Shutdown(ctx context.Context) {

}
