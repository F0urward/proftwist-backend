package jwt

import (
	"time"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(cfg *config.JwtConfig, user *entities.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"exp":     time.Now().Add(cfg.Expire * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}
