package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/F0urward/proftwist-backend/internal/server/ws/dto"
)

type WsClient struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
	Server *WsServer
	Send   chan dto.WebSocketMessage
	mu     sync.Mutex
	closed bool
}

func (c *WsClient) readPump() {
	defer func() {
		c.Server.unregister <- c
	}()

	c.Conn.SetReadLimit(c.Server.config.WebSocket.MaxMessageSize)
	if err := c.Conn.SetReadDeadline(time.Now().Add(c.Server.config.WebSocket.PongWait)); err != nil {
		c.Server.logger.WithFields(logrus.Fields{
			"client_id": c.ID,
			"error":     err,
		}).Warn("Failed to set read deadline")
	}

	c.Conn.SetPongHandler(func(string) error {
		if err := c.Conn.SetReadDeadline(time.Now().Add(c.Server.config.WebSocket.PongWait)); err != nil {
			c.Server.logger.WithFields(logrus.Fields{
				"client_id": c.ID,
				"error":     err,
			}).Warn("Failed to set read deadline in pong handler")
		}
		return nil
	})

	for {
		var message dto.WebSocketMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				c.Server.logger.WithError(err).WithField("client_id", c.ID).Error("WebSocket read error")
			}
			break
		}

		message.Timestamp = time.Now()

		c.Server.logger.WithFields(logrus.Fields{
			"client_id":    c.ID,
			"user_id":      c.UserID,
			"message_type": message.Type,
			"data_length":  len(message.Data),
		}).Info("WebSocket message received")

		c.Server.mutex.RLock()
		handler, exists := c.Server.messageHandlers[message.Type]
		c.Server.mutex.RUnlock()

		if exists && handler != nil {
			go func(msg dto.WebSocketMessage) {
				defer func() {
					if r := recover(); r != nil {
						c.Server.logger.WithFields(logrus.Fields{
							"client_id":    c.ID,
							"message_type": msg.Type,
							"recover":      r,
						}).Error("Panic in message handler")
					}
				}()

				if err := handler(c, msg); err != nil {
					c.Server.logger.WithError(err).WithFields(logrus.Fields{
						"client_id":    c.ID,
						"message_type": msg.Type,
					}).Error("Failed to handle message")
				} else {
					c.Server.logger.WithField("message_type", msg.Type).Info("Message processed successfully")
				}
			}(message)
		} else {
			c.Server.logger.WithField("message_type", message.Type).Warn("No handler for message type")
		}
	}
}

func (c *WsClient) writePump() {
	ticker := time.NewTicker(c.Server.config.WebSocket.PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				c.Server.logger.WithField("client_id", c.ID).Info("Send channel closed received in writePump")
				return
			}

			c.mu.Lock()
			if c.closed {
				c.mu.Unlock()
				return
			}

			if err := c.Conn.SetWriteDeadline(time.Now().Add(c.Server.config.WebSocket.WriteWait)); err != nil {
				c.mu.Unlock()
				c.Server.logger.WithFields(logrus.Fields{
					"client_id": c.ID,
					"error":     err,
				}).Warn("Failed to set write deadline")
				return
			}

			if err := c.Conn.WriteJSON(msg); err != nil {
				c.mu.Unlock()
				c.Server.unregister <- c
				c.Server.logger.WithFields(logrus.Fields{
					"client_id": c.ID,
					"error":     err,
				}).Error("Failed to write message to client")
				return
			}
			c.mu.Unlock()

		case <-ticker.C:
			c.mu.Lock()
			if c.closed {
				c.mu.Unlock()
				return
			}

			if err := c.Conn.SetWriteDeadline(time.Now().Add(c.Server.config.WebSocket.WriteWait)); err != nil {
				c.mu.Unlock()
				c.Server.logger.WithFields(logrus.Fields{
					"client_id": c.ID,
					"error":     err,
				}).Warn("Failed to set write deadline for ping")
				return
			}

			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.mu.Unlock()
				c.Server.logger.WithFields(logrus.Fields{
					"client_id": c.ID,
					"error":     err,
				}).Warn("Failed to write ping message")
				return
			}
			c.mu.Unlock()
		}
	}
}
