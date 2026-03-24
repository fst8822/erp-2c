package controller

import (
	"api-gateway/lib/observability/app_metrics"
	"api-gateway/lib/response"
	"api-gateway/lib/sl"
	"api-gateway/model"
	notifygrpc "api-gateway/server_grpc/proto/v1/notify"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
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
	notifyClientGRPC notifygrpc.NotifyServiceClient
}

func NewNotifyController(notifyClientGRPC notifygrpc.NotifyServiceClient) *NotifyController {
	return &NotifyController{notifyClientGRPC: notifyClientGRPC}
}

func (n *NotifyController) UpgradeConnection(resp http.ResponseWriter, r *http.Request) {
	const OP = "controller.notify.NotifyController.UpgradeConnection"
	logger := slog.With("OP", OP)

	userID, ok := model.UserIDFromContext(r.Context())
	if !ok {
		logger.Error("user id not found in context")
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
	app_metrics.WebsocketActiveConn.Inc()

	client := &model.ClientWS{
		UUID:   uuid.New(),
		UserID: userID,
		Conn:   conn,
		Cn:     make(chan model.Notification, model.NotificationBufferSize),
	}
	_ = client
	//todo спросить: как можно передать сложную структуру или кто должен создавать Upgrade
	//n.notifyService.Subscribe(context.TODO(), client)
	//n.notifyService.AddClient(client)
}
