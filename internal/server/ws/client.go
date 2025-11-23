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
}

func (c *WsClient) readPump() {
	defer func() {
		c.Server.unregister <- c
		if err := c.Conn.Close(); err != nil {
			c.Server.logger.WithFields(logrus.Fields{
				"client_id": c.ID,
				"error":     err,
			}).Warn("Error closing WebSocket connection in readPump")
		}
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

		if handler, exists := c.Server.messageHandlers[message.Type]; exists {
			go func() {
				c.Server.logger.WithField("message_type", message.Type).Debug("Processing message")
				if err := handler(c, message); err != nil {
					c.Server.logger.WithError(err).WithFields(logrus.Fields{
						"client_id":    c.ID,
						"message_type": message.Type,
					}).Error("Failed to handle message")
				} else {
					c.Server.logger.WithField("message_type", message.Type).Debug("Message processed successfully")
				}
			}()
		} else {
			c.Server.logger.WithField("message_type", message.Type).Warn("No handler for message type")
		}
	}
}

func (c *WsClient) writePump() {
	ticker := time.NewTicker(c.Server.config.WebSocket.PingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.Conn.Close(); err != nil {
			c.Server.logger.WithFields(logrus.Fields{
				"client_id": c.ID,
				"error":     err,
			}).Warn("Error closing WebSocket connection in writePump")
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(c.Server.config.WebSocket.WriteWait)); err != nil {
				c.Server.logger.WithFields(logrus.Fields{
					"client_id": c.ID,
					"error":     err,
				}).Warn("Failed to set write deadline")
				return
			}

			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					c.Server.logger.WithFields(logrus.Fields{
						"client_id": c.ID,
						"error":     err,
					}).Warn("Failed to write close message")
				}
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				c.Server.logger.WithError(err).WithField("client_id", c.ID).Error("Failed to write message")
				return
			}

		case <-ticker.C:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(c.Server.config.WebSocket.WriteWait)); err != nil {
				c.Server.logger.WithFields(logrus.Fields{
					"client_id": c.ID,
					"error":     err,
				}).Warn("Failed to set write deadline for ping")
				return
			}

			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.Server.logger.WithFields(logrus.Fields{
					"client_id": c.ID,
					"error":     err,
				}).Warn("Failed to write ping message")
				return
			}
		}
	}
}
