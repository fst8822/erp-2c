package notify

import (
	"context"
	"erp-2c/lib/sl"
	"fmt"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

type clientList map[*clientWS]bool

type Client interface {
	readMessage(ctx context.Context)
	writeMessage(ctx context.Context)
}

const (
	writeLimit   = 5 * time.Second
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
	MessageSize  = 512
)

type clientWS struct {
	conn      *websocket.Conn
	managerWS *ManagerWS
	msgCH     chan string
}

func NewClientWS(conn *websocket.Conn, managerWS *ManagerWS) *clientWS {
	return &clientWS{
		conn:      conn,
		managerWS: managerWS,
		msgCH:     make(chan string),
	}
}

// todo where check error
func (c *clientWS) readMessage(ctx context.Context) {
	const OP = "controller.notify.client_ws.readMessage"
	log := slog.With("OP", OP)

	defer func() {
		c.managerWS.removeClint(c)
	}()

	if err := c.initLimitRead(); err != nil {
		log.Error("Error init limit to Read message", sl.Err(err))
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Info("Context is don, exit from readMessage")
			return
		default:
			messageType, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseMessage,
					websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Error("Error read message from client, connection is closed", sl.Err(err))
				}
				return
			}

			log.Info("read message from client",
				slog.Int("messageType", messageType),
				slog.String("message", string(message)),
			)
			c.msgCH <- string(message)
		}
	}
}

func (c *clientWS) writeMessage(ctx context.Context) {
	const OP = "controller.notify.client_ws.writeMessage"
	log := slog.With("OP", OP)

	defer func() {
		c.managerWS.removeClint(c)
	}()
	ticket := time.NewTicker(pingInterval)
	defer ticket.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info("Context is don, exit from writeMessage")
			return
		case message, ok := <-c.msgCH:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Error("Connection is closed", sl.Err(err))
				}
				return
			}

			log.Info("Send message to client", slog.StringValue(message))
			str := fmt.Sprintf(message + "from server")

			err := c.conn.SetWriteDeadline(time.Now().Add(writeLimit))
			if err != nil {
				log.Error("error set WriteDeadline", sl.Err(err))
				return
			}
			err = c.conn.WriteMessage(websocket.TextMessage, []byte(str))
			if err != nil {
				log.Error("Error send message to client, connection is closed", sl.Err(err))
				return
			}

		case <-ticket.C:
			log.Info("PING")
			err := c.conn.SetWriteDeadline(time.Now().Add(writeLimit))
			if err != nil {
				log.Error("error set WriteDeadline", sl.Err(err))
				return
			}
			err = c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Error("Error send PING connection is closed", sl.Err(err))
				return
			}
		}
	}
}

func (c *clientWS) initLimitRead() error {
	const OP = "controller.notify.client_ws.initLimitRead"
	log := slog.With("OP", OP)

	c.conn.SetReadLimit(MessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Error("error set ReadDeadline", sl.Err(err))
		return err
	}

	c.conn.SetPongHandler(func(appData string) error {
		log.Info("PONG RECEIVED")
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Error("error Reset the ReadDeadline", sl.Err(err))
		}
		return err
	})
	return nil
}
