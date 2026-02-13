package notify

import (
	"context"
	"erp-2c/lib/response"
	"erp-2c/lib/sl"
	"erp-2c/service/use_cases"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrade = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type ManagerWS struct {
	services *use_cases.Manager
	sync.RWMutex
	clients clientList
}

func NewManagerWS(services *use_cases.Manager) *ManagerWS {
	return &ManagerWS{services: services, clients: make(clientList)}
}

func (m *ManagerWS) ServerWS(resp http.ResponseWriter, r *http.Request) {
	const OP = "controller.notify.manager_ws.ServerWS"
	logger := slog.With("OP", OP)
	ctxTODO := context.TODO()

	conn, err := websocketUpgrade.Upgrade(resp, r, nil)
	if err != nil {
		logger.Error("Error websocket Upgrade connection", sl.Err(err))
		response.InternalServerError().SendResponse(resp, r)
		return
	}
	logger.Info("has new connection", slog.Any("LocalAddr", conn.LocalAddr()))
	client := NewClientWS(conn, m)

	m.addClient(client)

	go client.aliveConnection(ctxTODO)
}

func (m *ManagerWS) addClient(client *ClientWS) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
}

func (m *ManagerWS) removeClint(client *ClientWS) {
	m.Lock()
	defer m.Unlock()
	if ok := m.clients[client]; ok {
		err := client.conn.Close()
		if err != nil {
			slog.Error("Can not close client connection", sl.Err(err))
			return
		}
		delete(m.clients, client)
	}
}
