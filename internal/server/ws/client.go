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
		c.mu.Lock()
		if c.Conn != nil && !c.closed {
			c.closed = true
			if err := c.Conn.Close(); err != nil {
				c.Server.logger.WithFields(logrus.Fields{
					"client_id": c.ID,
					"error":     err,
				}).Warn("Error closing WebSocket connection in readPump")
			}
		}
		c.mu.Unlock()
		c.Server.logger.WithField("client_id", c.ID).Info("WebSocket client disconnected")
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
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

		c.Server.logger.WithField("raw_data", string(message.Data)).Debug("Raw message data")

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

				c.Server.logger.WithField("message_type", msg.Type).Debug("Processing message")
				if err := handler(c, msg); err != nil {
					c.Server.logger.WithError(err).WithFields(logrus.Fields{
						"client_id":    c.ID,
						"message_type": msg.Type,
					}).Error("Failed to handle message")
				} else {
					c.Server.logger.WithField("message_type", msg.Type).Debug("Message processed successfully")
				}
			}(message)
		} else {
			c.Server.logger.WithField("message_type", message.Type).Warn("No handler for message type")
		}
	}
}

func (c *WsClient) writePump() {
	ticker := time.NewTicker(c.Server.config.WebSocket.PingPeriod)
	defer func() {
		ticker.Stop()
		c.mu.Lock()
		if c.Conn != nil && !c.closed {
			c.closed = true
			if err := c.Conn.Close(); err != nil {
				c.Server.logger.WithFields(logrus.Fields{
					"client_id": c.ID,
					"error":     err,
				}).Warn("Error closing WebSocket connection in writePump")
			}
		}
		c.mu.Unlock()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.mu.Lock()
			if c.Conn == nil || c.closed {
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

			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					c.mu.Unlock()
					c.Server.logger.WithFields(logrus.Fields{
						"client_id": c.ID,
						"error":     err,
					}).Warn("Failed to write close message")
				}
				c.mu.Unlock()
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				c.mu.Unlock()
				c.Server.logger.WithError(err).WithField("client_id", c.ID).Error("Failed to write message")
				return
			}
			c.mu.Unlock()

		case <-ticker.C:
			c.mu.Lock()
			if c.Conn == nil || c.closed {
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
