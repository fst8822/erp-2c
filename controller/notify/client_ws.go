package notify

import (
	"context"
	"erp-2c/lib/sl"
	"erp-2c/model"
	"erp-2c/service"
	"erp-2c/service/use_cases"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

type ClientList map[ClientWS]bool

type Client interface {
	subscribe(ctx context.Context)
}

const (
	writeLimit   = 5 * time.Second
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
	MessageSize  = 512
)

type ClientWS struct {
	Conn     *websocket.Conn
	Services *use_cases.Manager
	Event    chan model.DeliveryDB
}

func NewClientWS(conn *websocket.Conn, services *use_cases.Manager) *ClientWS {
	return &ClientWS{
		Conn:     conn,
		Services: services,
		Event:    make(chan model.DeliveryDB),
	}
}

func (c *ClientWS) subscribe(ctx context.Context) {
	const OP = "controller.notify.client_ws.subscribe"
	log := slog.With("OP", OP)

	defer func() {
		c.Services.NotifyService.RemoveClint(c)
	}()

	for {
		select {
		case <-ctx.Done():
			log.Info("Context is don, exit from subscribe")
			return

		case event, ok := <-c.Event:
			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Error("Connection is closed", sl.Err(err))
				}
				return
			}

			log.Info("Send message to client", slog.AnyValue(event))
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeLimit))
			if err != nil {
				log.Error("error set WriteDeadline", sl.Err(err))
				return
			}
			err = c.Conn.WriteJSON(event)
			if err != nil {
				log.Error("Error send message to client, connection is closed", sl.Err(err))
				return
			}
		}
	}
}

func (c *ClientWS) aliveConnection(ctx context.Context) {
	const OP = "controller.notify.client_ws.aliveConnection"
	log := slog.With("OP", OP)

	defer func() {
		c.Services.NotifyService.RemoveClint(c)
	}()
	if err := c.initLimitRead(); err != nil {
		log.Error("Error set  init limit Read", sl.Err(err))
		return
	}
	ticket := time.NewTicker(pingInterval)
	defer ticket.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("Context is don, exit from aliveConnection")
				return
			case <-ticket.C:
				log.Info("PING")
				err := c.Conn.SetWriteDeadline(time.Now().Add(writeLimit))
				if err != nil {
					log.Error("error set WriteDeadline", sl.Err(err))
					return
				}
				if err = c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Error("Error send PING connection is closed", sl.Err(err))
					return
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Info("Context is done, exit from aliveConnection")
			return
		default:
			_, _, err := c.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseMessage,
					websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Error("Error read message from client, connection is closed", sl.Err(err))
				}
				return
			}
		}
	}
}

func (c *ClientWS) initLimitRead() error {
	const OP = "controller.notify.client_ws.initLimitRead"
	log := slog.With("OP", OP)

	c.Conn.SetReadLimit(MessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Error("error set ReadDeadline", sl.Err(err))
		return err
	}

	c.Conn.SetPongHandler(func(appData string) error {
		log.Info("PONG RECEIVED")
		err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Error("error Reset the ReadDeadline", sl.Err(err))
		}
		return err
	})
	return nil
}
