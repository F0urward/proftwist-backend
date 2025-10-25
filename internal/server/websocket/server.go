package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/websocketclient/dto"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type MessageHandler func(*Client, dto.WebSocketMessage) error

type Server struct {
	config          *config.WebSocketConfig
	upgrader        websocket.Upgrader
	clients         map[*Client]bool
	clientsByUserID map[string][]*Client
	register        chan *Client
	unregister      chan *Client
	broadcast       chan dto.WebSocketMessage
	messageHandlers map[dto.WebSocketMessageType]MessageHandler
	mutex           sync.RWMutex
	logger          *logrus.Logger
}

type Client struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
	Server *Server
	Send   chan dto.WebSocketMessage
	mu     sync.Mutex
}

func NewWebSocketServer(cfg *config.WebSocketConfig) *Server {
	return &Server{
		config: cfg,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients:         make(map[*Client]bool),
		clientsByUserID: make(map[string][]*Client),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		broadcast:       make(chan dto.WebSocketMessage),
		messageHandlers: make(map[dto.WebSocketMessageType]MessageHandler),
		logger:          logrus.New(),
	}
}

func (s *Server) RegisterMessageHandler(messageType dto.WebSocketMessageType, handler MessageHandler) {
	s.messageHandlers[messageType] = handler
}

func (s *Server) EnableDebugLogging() {
	s.logger.SetLevel(logrus.DebugLevel)
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID string) error {
	s.logger.WithFields(logrus.Fields{
		"user_id":     userID,
		"remote_addr": r.RemoteAddr,
	}).Info("🔌 WebSocket connection attempt")

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.WithError(err).Error("❌ WebSocket upgrade failed")
		return err
	}

	client := &Client{
		ID:     generateClientID(),
		UserID: userID,
		Conn:   conn,
		Server: s,
		Send:   make(chan dto.WebSocketMessage, 256),
	}

	s.register <- client

	s.logger.WithFields(logrus.Fields{
		"client_id": client.ID,
		"user_id":   userID,
	}).Info("✅ WebSocket client connected")

	go client.writePump()
	go client.readPump()

	return nil
}

func (c *Client) readPump() {
	defer func() {
		c.Server.unregister <- c
		c.Conn.Close()
		c.Server.logger.WithField("client_id", c.ID).Info("🔌 WebSocket client disconnected")
	}()

	c.Conn.SetReadLimit(c.Server.config.MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(c.Server.config.PongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(c.Server.config.PongWait))
		return nil
	})

	for {
		var message dto.WebSocketMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Server.logger.WithError(err).WithField("client_id", c.ID).Error("❌ WebSocket read error")
			}
			break
		}

		message.UserID = c.UserID
		message.Timestamp = time.Now()

		c.Server.logger.WithFields(logrus.Fields{
			"client_id":    c.ID,
			"user_id":      c.UserID,
			"message_type": message.Type,
			"data_length":  len(message.Data),
		}).Info("📨 WebSocket message received")

		c.Server.logger.WithField("raw_data", string(message.Data)).Debug("Raw message data")

		if handler, exists := c.Server.messageHandlers[message.Type]; exists {
			go func() {
				c.Server.logger.WithField("message_type", message.Type).Debug("🔄 Processing message")
				if err := handler(c, message); err != nil {
					c.Server.logger.WithError(err).WithFields(logrus.Fields{
						"client_id":    c.ID,
						"message_type": message.Type,
					}).Error("❌ Failed to handle message")

					errorMsg := dto.WebSocketMessage{
						Type: dto.WebSocketMessageTypeError,
						Data: mustMarshal(dto.ErrorMessageData{
							Code:    "HANDLER_ERROR",
							Message: err.Error(),
						}),
						Timestamp: time.Now(),
					}
					c.Send <- errorMsg
				} else {
					c.Server.logger.WithField("message_type", message.Type).Debug("✅ Message processed successfully")
				}
			}()
		} else {
			c.Server.logger.WithField("message_type", message.Type).Warn("⚠️ No handler for message type")
		}
	}
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client] = true
			s.clientsByUserID[client.UserID] = append(s.clientsByUserID[client.UserID], client)
			s.mutex.Unlock()
			s.logger.WithFields(logrus.Fields{
				"client_id": client.ID,
				"user_id":   client.UserID,
			}).Info("Client connected")

		case client := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.Send)

				if userClients, exists := s.clientsByUserID[client.UserID]; exists {
					for i, c := range userClients {
						if c == client {
							s.clientsByUserID[client.UserID] = append(userClients[:i], userClients[i+1:]...)
							break
						}
					}
					if len(s.clientsByUserID[client.UserID]) == 0 {
						delete(s.clientsByUserID, client.UserID)
					}
				}
			}
			s.mutex.Unlock()
			s.logger.WithFields(logrus.Fields{
				"client_id": client.ID,
				"user_id":   client.UserID,
			}).Info("Client disconnected")

		case message := <-s.broadcast:
			s.broadcastMessage(message)
		}
	}
}

func (s *Server) Broadcast(message dto.WebSocketMessage) error {
	s.broadcast <- message
	return nil
}

func (s *Server) SendToUser(userID string, message dto.WebSocketMessage) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if clients, exists := s.clientsByUserID[userID]; exists {
		for _, client := range clients {
			select {
			case client.Send <- message:
			default:
				s.logger.WithFields(logrus.Fields{
					"client_id": client.ID,
					"user_id":   userID,
				}).Warn("Failed to send message to client - channel full")
				go s.closeClient(client)
			}
		}
	}
	return nil
}

func (s *Server) SendToUsers(userIDs []string, message dto.WebSocketMessage) error {
	for _, userID := range userIDs {
		s.SendToUser(userID, message)
	}
	return nil
}

func (s *Server) broadcastMessage(message dto.WebSocketMessage) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for client := range s.clients {
		select {
		case client.Send <- message:
		default:
			s.logger.WithField("client_id", client.ID).Warn("Failed to broadcast message - channel full")
			go s.closeClient(client)
		}
	}
}

func (s *Server) closeClient(client *Client) {
	client.mu.Lock()
	defer client.mu.Unlock()

	if client.Conn != nil {
		client.Conn.Close()
	}
	s.unregister <- client
}

func (c *Client) writePump() {
	ticker := time.NewTicker(c.Server.config.PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(c.Server.config.WriteWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				c.Server.logger.WithError(err).WithField("client_id", c.ID).Error("Failed to write message")
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(c.Server.config.WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func generateClientID() string {
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
}

func mustMarshal(v interface{}) json.RawMessage {
	bytes, _ := json.Marshal(v)
	return bytes
}

func (s *Server) Logger() *logrus.Logger {
	return s.logger
}
