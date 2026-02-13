package use_cases

import (
	"context"
	"erp-2c/controller/notify"
	"erp-2c/lib/sl"
	"erp-2c/model"
	"sync"

	"golang.org/x/exp/slog"
)

type NotifyService struct {
	sync.RWMutex
	clients notify.ClientList
	Event   chan model.DeliveryDB
}

func NewNotifyService() *NotifyService {
	return &NotifyService{Event: make(chan model.DeliveryDB)}
}

func (n *NotifyService) Subscribe(ctx context.Context, client notify.Client) {
	//TODO implement me
	panic("implement me")
}

func (n *NotifyService) AddClient(client notify.ClientWS) {
	n.Lock()
	defer n.Unlock()
	n.clients[client] = true
}

func (n *NotifyService) RemoveClint(client notify.ClientWS) {
	n.Lock()
	defer n.Unlock()
	if ok := n.clients[client]; ok {
		err := client.Conn.Close()
		if err != nil {
			slog.Error("Can not close client connection", sl.Err(err))
			return
		}
		delete(n.clients, client)
	}
}
