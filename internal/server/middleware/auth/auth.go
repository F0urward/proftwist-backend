package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/internal/server/websocket"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/pkg/jwt"
	"github.com/F0urward/proftwist-backend/services/auth"
)

type AuthMiddleware struct {
	authRedisRepo auth.RedisRepository
	cfg           *config.Config
}

func NewAuthMiddleware(authRedisRepo auth.RedisRepository, cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		authRedisRepo: authRedisRepo,
		cfg:           cfg,
	}
}

func (a *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "AuthMiddleware.AuthMiddleware"
		logger := logctx.GetLogger(context.Background()).WithFields(map[string]interface{}{
			"op": op,
		})

		ctx := r.Context()
		var tokenString string

		cookie, err := r.Cookie(a.cfg.Auth.Jwt.Cookie.Name)
		if err == nil {
			tokenString = cookie.Value
			logger.Debug("token found in cookies")
		} else {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
				logger.Debug("token found in Authorization header")
			}
		}

		if tokenString == "" {
			logger.Warn("jwt token not found in cookies or Authorization header")
			utils.JSONError(r.Context(), w, http.StatusUnauthorized, "unauthorized")
			return
		}

		userID, err := a.validateToken(ctx, tokenString)
		if err != nil {
			logger.Warnf("token validation failed: %v", err)
			utils.JSONError(r.Context(), w, http.StatusUnauthorized, "unauthorized")
			return
		}

		ctx = context.WithValue(ctx, utils.UserIDKey{}, userID)
		//ctx = context.WithValue(ctx, utils.RoleKey{}, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WebSocketAuthMiddleware middleware для WebSocket аутентификации
func (a *AuthMiddleware) WebSocketAuthMiddleware(wsServer *websocket.Server, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "AuthMiddleware.WebSocketAuthMiddleware"
		logger := logctx.GetLogger(context.Background()).WithField("op", op)

		// Извлекаем токен из query параметров или заголовков
		token := r.URL.Query().Get("token")
		if token == "" {
			token = r.Header.Get("Authorization")
			if strings.HasPrefix(token, "Bearer ") {
				token = strings.TrimPrefix(token, "Bearer ")
			}
		}

		if token == "" {
			logger.Warn("jwt token not found in query parameters or Authorization header")
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		userID, err := a.validateToken(r.Context(), token)
		if err != nil {
			logger.Warnf("websocket token validation failed: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Сохраняем userID в контекст
		ctx := context.WithValue(r.Context(), utils.UserIDKey{}, userID)
		r = r.WithContext(ctx)

		logger.WithField("user_id", userID).Debug("websocket authentication successful")
		next.ServeHTTP(w, r)
	})
}

// validateToken универсальная функция валидации токена
func (a *AuthMiddleware) validateToken(ctx context.Context, tokenString string) (string, error) {
	if tokenString == "" {
		return "", fmt.Errorf("empty token")
	}

	// Парсим JWT
	claims, err := jwt.ParseJWT(&a.cfg.Auth.Jwt, tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to parse jwt: %w", err)
	}

	if claims.IsExpired() {
		return "", fmt.Errorf("jwt is expired")
	}

	// Проверяем blacklist
	isInBlackList, err := a.authRedisRepo.IsInBlacklist(ctx, claims.UserID, tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to check token: %w", err)
	}
	if isInBlackList {
		return "", fmt.Errorf("token in blacklist")
	}

	return claims.UserID, nil
}
