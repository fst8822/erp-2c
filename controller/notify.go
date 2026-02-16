package controller

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

type NotifyController struct {
	services *use_cases.Manager
}

func NewNotifyController(services *use_cases.Manager) *NotifyController {
	return &NotifyController{services: services}
}

func (n *NotifyController) UpgradeConnection(resp http.ResponseWriter, r *http.Request) {
	const OP = "controller.notify.manager_ws.UpgradeConnection"
	logger := slog.With("OP", OP)
	ctxTODO := context.TODO()

	userId := r.Context().Value("userIdKey")
	id, ok := userId.(int64)
	if !ok {
		logger.Error("User id dont exist on context")
		response.Unauthorized("Unauthorized").SendResponse(resp, r)
		return
	}

	conn, err := websocketUpgrade.Upgrade(resp, r, nil)
	if err != nil {
		logger.Error("Error websocket Upgrade connection", sl.Err(err))
		response.InternalServerError().SendResponse(resp, r)
		return
	}
	logger.Info("has new connection", slog.Any("LocalAddr", conn.LocalAddr()))

	n.services.NotifyService.Subscribe(ctxTODO, conn, id)
}
