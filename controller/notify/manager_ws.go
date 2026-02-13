package notify

import (
	"context"
	"erp-2c/lib/response"
	"erp-2c/lib/sl"
	"erp-2c/service/use_cases"
	"log/slog"
	"net/http"

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
}

func NewManagerWS(services *use_cases.Manager) *ManagerWS {
	return &ManagerWS{services: services}
}

func (m *ManagerWS) UpgradeConnection(resp http.ResponseWriter, r *http.Request) {
	const OP = "controller.notify.manager_ws.UpgradeConnection"
	logger := slog.With("OP", OP)
	ctxTODO := context.TODO()

	conn, err := websocketUpgrade.Upgrade(resp, r, nil)
	if err != nil {
		logger.Error("Error websocket Upgrade connection", sl.Err(err))
		response.InternalServerError().SendResponse(resp, r)
		return
	}
	logger.Info("has new connection", slog.Any("LocalAddr", conn.LocalAddr()))

	client := NewClientWS(conn, m.services)
	m.services.NotifyService.Subscribe(ctxTODO, client)
}
