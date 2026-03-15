package model

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const NotificationBufferSize = 61

type ClientWS struct {
	UUID   uuid.UUID
	UserID int64
	Conn   *websocket.Conn
	Cn     chan Notification
}
