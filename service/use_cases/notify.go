package use_cases

import (
	"context"
	"erp-2c/lib/observability/app_metrics"
	"erp-2c/lib/sl"
	"erp-2c/model"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

type NotifyService struct {
	mu      sync.RWMutex
	clients map[int64]*model.ClientWS
	done    chan struct{}
	wg      sync.WaitGroup
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

	ctxDone, cancel := context.WithCancel(ctx)
	go n.checkClientConn(cancel, client)

	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		defer func() {
			n.RemoveClient(client)
			defer app_metrics.WebsocketActiveConn.Dec()
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

func (n *NotifyService) SendNotify(notification model.Notification) {
	const OP = "services.use_cases.notify.NotifyService.SendNotify"
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
	_, _, err := client.Conn.ReadMessage()
	defer cancel()
	if err != nil {
		return
	}
}
