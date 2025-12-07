package ws

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/server/ws/dto"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type WsRegistrar interface {
	RegisterHandlers(s *WsServer)
}

type MessageHandler func(*WsClient, dto.WebSocketMessage) error

type WsServer struct {
	config          *config.Config
	upgrader        websocket.Upgrader
	clients         map[*WsClient]bool
	clientsByUserID map[string][]*WsClient
	register        chan *WsClient
	unregister      chan *WsClient
	broadcast       chan dto.WebSocketMessage
	messageHandlers map[dto.WebSocketMessageType]MessageHandler
	mutex           sync.RWMutex
	logger          *logrus.Logger
	Registrars      []WsRegistrar
}

func (s *WsServer) RegisterHandlers() {
	for _, registrar := range s.Registrars {
		registrar.RegisterHandlers(s)
	}
}

func New(cfg *config.Config, registrars ...WsRegistrar) *WsServer {
	server := &WsServer{
		config:          cfg,
		clients:         make(map[*WsClient]bool),
		clientsByUserID: make(map[string][]*WsClient),
		register:        make(chan *WsClient),
		unregister:      make(chan *WsClient),
		broadcast:       make(chan dto.WebSocketMessage),
		messageHandlers: make(map[dto.WebSocketMessageType]MessageHandler),
		logger:          logrus.New(),
		Registrars:      registrars,
	}

	server.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return server
}

func (s *WsServer) RegisterMessageHandler(messageType dto.WebSocketMessageType, handler MessageHandler) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.messageHandlers[messageType] = handler
}

func (s *WsServer) EnableDebugLogging() {
	s.logger.SetLevel(logrus.InfoLevel)
}

func (s *WsServer) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID string) error {
	s.logger.WithFields(logrus.Fields{
		"user_id":     userID,
		"remote_addr": r.RemoteAddr,
	}).Info("WebSocket connection attempt")

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.WithError(err).Error("WebSocket upgrade failed")
		return err
	}

	client := &WsClient{
		ID:     generateClientID(),
		UserID: userID,
		Conn:   conn,
		Server: s,
		Send:   make(chan dto.WebSocketMessage, 256),
		mu:     sync.Mutex{},
		closed: false,
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

func (s *WsServer) Run() {
	s.RegisterHandlers()
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

			client.mu.Lock()
			if !client.closed && client.Conn != nil {
				client.closed = true
				if err := client.Conn.Close(); err != nil {
					s.logger.WithFields(logrus.Fields{
						"client_id": client.ID,
						"error":     err,
					}).Warn("Error closing WebSocket connection")
				}
			}
			client.mu.Unlock()

			s.logger.WithFields(logrus.Fields{
				"client_id": client.ID,
				"user_id":   client.UserID,
			}).Info("Client disconnected")

		case message := <-s.broadcast:
			s.broadcastMessage(message)
		}
	}
}

func (s *WsServer) Broadcast(message dto.WebSocketMessage) error {
	s.broadcast <- message
	return nil
}

func (s *WsServer) SendToUser(userID string, message dto.WebSocketMessage) error {
	s.mutex.RLock()
	clients := make([]*WsClient, len(s.clientsByUserID[userID]))
	n := copy(clients, s.clientsByUserID[userID])
	s.mutex.RUnlock()

	if n == 0 {
		return nil
	}

	clients = clients[:n]

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
	return nil
}

func (s *WsServer) SendToUsers(userIDs []string, message dto.WebSocketMessage) error {
	for _, userID := range userIDs {
		if err := s.SendToUser(userID, message); err != nil {
			s.logger.WithFields(logrus.Fields{
				"user_id": userID,
			}).Warn("Failed to send message to user")
		}
	}
	return nil
}

func (s *WsServer) broadcastMessage(message dto.WebSocketMessage) {
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

func (s *WsServer) closeClient(client *WsClient) {
	s.unregister <- client
}

func generateClientID() string {
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
}

func (s *WsServer) Logger() *logrus.Logger {
	return s.logger
}
