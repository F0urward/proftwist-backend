package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/server/websocket/dto"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type MessageHandler func(*Client, dto.WebSocketMessage) error

type Server struct {
	config          *config.Config
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

func NewWebSocketServer(cfg *config.Config) *Server {
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
	}).Info("WebSocket connection attempt")

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.WithError(err).Error("WebSocket upgrade failed")
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
	}).Info("WebSocket client connected")

	go client.writePump()
	go client.readPump()

	return nil
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
	var firstErr error
	for _, userID := range userIDs {
		if err := s.SendToUser(userID, message); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			s.logger.WithFields(logrus.Fields{
				"user_id": userID,
				"error":   err,
			}).Warn("Failed to send message to user")
		}
	}
	return firstErr
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
		if err := client.Conn.Close(); err != nil {
			s.logger.WithFields(logrus.Fields{
				"client_id": client.ID,
				"error":     err,
			}).Warn("Error closing WebSocket connection")
		}
	}
	s.unregister <- client
}

func generateClientID() string {
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
}

func (s *Server) Logger() *logrus.Logger {
	return s.logger
}
