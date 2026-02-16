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

type NotifyService struct {
	sync.RWMutex
	clients map[int64]*websocket.Conn
	Event   chan model.DeliveryDB
}

const (
	writeLimit = 5 * time.Second
)

func NewNotifyService() *NotifyService {
	return &NotifyService{Event: make(chan model.DeliveryDB)}
}

func (n *NotifyService) Subscribe(ctx context.Context, conn *websocket.Conn, userId int64) {
	const OP = "controller.notify.client_ws.subscribe"
	log := slog.With("OP", OP)

	n.AddClient(conn, userId)
	defer func() {
		n.RemoveClint(userId)
	}()

	for {
		select {
		case <-ctx.Done():
			log.Info("Context is don, exit from subscribe")
			return

		case event, ok := <-n.Event:
			if !ok {
				if err := conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Error("Connection is closed", sl.Err(err))
				}
				return
			}

			log.Info("Send message to client", slog.AnyValue(event))
			err := conn.SetWriteDeadline(time.Now().Add(writeLimit))
			if err != nil {
				log.Error("error set WriteDeadline", sl.Err(err))
				return
			}
			err = conn.WriteJSON(event)
			if err != nil {
				log.Error("Error send message to client, connection is closed", sl.Err(err))
				return
			}
		}
	}
}

func (n *NotifyService) AddClient(conn *websocket.Conn, userId int64) {
	n.Lock()
	defer n.Unlock()
	n.clients[userId] = conn
}

func (n *NotifyService) RemoveClint(userId int64) {
	n.Lock()
	defer n.Unlock()
	if conn, ok := n.clients[userId]; ok {
		err := conn.Close()
		if err != nil {
			slog.Error("Can not close client connection", sl.Err(err))
			return
		}
		delete(n.clients, userId)
	}
}
