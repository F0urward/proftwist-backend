package usecase

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/pkg/jwt"
	"github.com/F0urward/proftwist-backend/services/auth"
	"github.com/F0urward/proftwist-backend/services/auth/dto"
)

type AuthUsecase struct {
	cfg          *config.Config
	postgresRepo auth.PostgresRepository
	redisRepo    auth.RedisRepository
}

func NewAuthUsecase(postgresRepo auth.PostgresRepository, redisRepo auth.RedisRepository, cfg *config.Config) auth.Usecase {
	return &AuthUsecase{
		cfg:          cfg,
		postgresRepo: postgresRepo,
		redisRepo:    redisRepo,
	}
}

func (uc *AuthUsecase) Register(ctx context.Context, request *dto.RegisterRequestDTO) (*dto.UserTokenDTO, error) {
	const op = "AuthUsecase.Register"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":       op,
		"email":    request.Email,
		"username": request.Username,
	})

	existingUser, err := uc.postgresRepo.GetUserByEmail(ctx, request.Email)
	if err != nil && !errs.IsNotFoundError(err) {
		logger.WithError(err).Error("failed to check existing user")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if existingUser != nil {
		logger.Warn("user with this email already exists")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrAlreadyExists)
	}

	passwordHash, err := hashPassword(request.Password)
	if err != nil {
		logger.WithError(err).Error("failed to hash password")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	newUser := dto.RegisterRequestToEntity(request, passwordHash)

	createdUser, err := uc.postgresRepo.CreateUser(ctx, newUser)
	if err != nil {
		logger.WithError(err).Error("failed to create user")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.GenerateJWT(&uc.cfg.Auth.Jwt, createdUser)
	if err != nil {
		logger.WithError(err).Error("failed to generate JWT token")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	response := dto.UserTokenToDTO(createdUser, token)

	logger.WithField("user_id", createdUser.ID.String()).Info("user registered successfully")
	return response, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, request *dto.LoginRequestDTO) (*dto.UserTokenDTO, error) {
	const op = "AuthUsecase.Login"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":    op,
		"email": request.Email,
	})

	user, err := uc.postgresRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		if errs.IsNotFoundError(err) {
			logger.Warn("user not found")
			return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
		logger.WithError(err).Error("failed to get user by email")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = checkPasswordHash(request.Password, user.PasswordHash)
	if err != nil {
		logger.WithError(err).Warn("invalid password")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrInvalidCredentials)
	}

	token, err := jwt.GenerateJWT(&uc.cfg.Auth.Jwt, user)
	if err != nil {
		logger.WithError(err).Error("failed to generate JWT token")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	response := dto.UserTokenToDTO(user, token)

	logger.WithField("user_id", user.ID.String()).Info("user logged in successfully")
	return response, nil
}

func (uc *AuthUsecase) Logout(ctx context.Context, token string) error {
	const op = "AuthUsecase.Logout"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	claims, err := jwt.ParseJWT(&uc.cfg.Auth.Jwt, token)
	if err != nil {
		logger.WithError(err).Error("failed to parse token")
		return fmt.Errorf("%s: %w", op, errs.ErrInvalidToken)
	}

	if err := uc.redisRepo.AddToBlacklist(ctx, claims.UserID, token); err != nil {
		logger.WithError(err).Error("failed to add token to blacklist")
		return fmt.Errorf("%s: %w", op, errs.ErrInternal)
	}

	logger.Info("user logged out successfully")
	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// func generateJWT(user *entities.User) (string, error) {
// 	return "jwt_token_for_user_" + user.ID.String(), nil
// }
