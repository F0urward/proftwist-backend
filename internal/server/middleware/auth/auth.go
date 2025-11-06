package auth

import (
	"context"
	"net/http"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
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
		const op = "AuthUsecase.Register"
		logger := logctx.GetLogger(context.Background()).WithFields(map[string]interface{}{
			"op": op,
		})

		ctx := r.Context()

		cookie, err := r.Cookie(a.cfg.Auth.Jwt.Cookie.Name)
		if err != nil {
			logger.Warnf("jwt token not found in cookies: %v", err)
			utils.JSONError(r.Context(), w, http.StatusUnauthorized, "unauthorized")
			return
		}

		claims, err := jwt.ParseJWT(&a.cfg.Auth.Jwt, cookie.Value)
		if err != nil {
			logger.Warnf("failed to parse jwt: %v", err)
			utils.JSONError(r.Context(), w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if claims.IsExpired() {
			logger.Warnf("jwt is expired: %v", err)
			utils.JSONError(r.Context(), w, http.StatusUnauthorized, "jwt is expired")
			return
		}

		isInBlackList, err := a.authRedisRepo.IsInBlacklist(ctx, claims.UserID, cookie.Value)
		if err != nil {
			logger.Warnf("failed to check token: %v", err)
			utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to check token")
			return
		}
		if isInBlackList {
			logger.Warnf("unauthorized token: %v", err)
			utils.JSONError(r.Context(), w, http.StatusUnauthorized, "unauthorized")
			return
		}

		ctx = context.WithValue(ctx, utils.UserIDKey{}, claims.UserID)
		ctx = context.WithValue(ctx, utils.RoleKey{}, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
