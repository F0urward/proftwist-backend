package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/pkg/jwt"
	"github.com/F0urward/proftwist-backend/services/auth"
	"github.com/F0urward/proftwist-backend/services/auth/dto"
	"github.com/google/uuid"
)

type AuthUsecase struct {
	cfg          *config.Config
	postgresRepo auth.PostgresRepository
	redisRepo    auth.RedisRepository
	awsRepo      auth.AWSRepository
	vkWebapi     auth.VKWebapi
}

func NewAuthUsecase(postgresRepo auth.PostgresRepository, redisRepo auth.RedisRepository, awsRepo auth.AWSRepository, vkWebapi auth.VKWebapi, cfg *config.Config) auth.Usecase {
	return &AuthUsecase{
		cfg:          cfg,
		postgresRepo: postgresRepo,
		redisRepo:    redisRepo,
		awsRepo:      awsRepo,
		vkWebapi:     vkWebapi,
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

	newUser.AvatarUrl = uc.generateAWSMinioURL(uc.cfg.AWS.AvatarBucketName, "default.jpg")

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

func (uc *AuthUsecase) GetMe(ctx context.Context, userID uuid.UUID) (*dto.UserDTO, error) {
	const op = "AuthUsecase.GetMe"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID,
	})

	user, err := uc.postgresRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errs.IsNotFoundError(err) {
			logger.WithError(err).Warn("user not found")
			return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
		logger.WithError(err).Error("failed to get user by id")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userDTO := dto.UserEntityToDTO(user)

	logger.Info("successfully retrieved user info")
	return &userDTO, nil
}

func (uc *AuthUsecase) GetByID(ctx context.Context, userID uuid.UUID) (*dto.GetUserByIDResponseDTO, error) {
	const op = "AuthUsecase.GetByID"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID.String(),
	})

	user, err := uc.postgresRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errs.IsNotFoundError(err) {
			logger.WithError(err).Warn("user not found")
			return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
		logger.WithError(err).Error("failed to get user by id")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userDTO := dto.UserEntityToDTO(user)
	response := &dto.GetUserByIDResponseDTO{
		User: userDTO,
	}

	logger.Info("successfully retrieved user by ID")
	return response, nil
}

func (uc *AuthUsecase) Update(ctx context.Context, userID uuid.UUID, request *dto.UpdateUserRequestDTO) error {
	const op = "AuthUsecase.Update"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID.String(),
	})

	existingUser, err := uc.postgresRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errs.IsNotFoundError(err) {
			logger.WithError(err).Warn("user not found")
			return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
		logger.WithError(err).Error("failed to get user by id")
		return fmt.Errorf("%s: %w", op, err)
	}

	if request.Email != "" && request.Email != existingUser.Email {
		userWithEmail, err := uc.postgresRepo.GetUserByEmail(ctx, request.Email)
		if err != nil && !errs.IsNotFoundError(err) {
			logger.WithError(err).Error("failed to check email availability")
			return fmt.Errorf("%s: %w", op, err)
		}

		if userWithEmail != nil && userWithEmail.ID != userID {
			logger.WithField("email", request.Email).Warn("email already taken by another user")
			return fmt.Errorf("%s: %w", op, errs.ErrAlreadyExists)
		}
	}

	updatedUser := dto.UpdateUserRequestToEntity(existingUser, request)

	err = uc.postgresRepo.UpdateUser(ctx, updatedUser)
	if err != nil {
		logger.WithError(err).Error("failed to update user")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("successfully updated user")
	return nil
}

func (uc *AuthUsecase) UploadAvatar(ctx context.Context, request *dto.UploadAvatarRequestDTO) (*dto.UploadAvatarResponseDTO, error) {
	const op = "AuthUsecase.UploadAvatar"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": request.UserID.String(),
	})

	existingUser, err := uc.postgresRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		if errs.IsNotFoundError(err) {
			logger.WithError(err).Warn("user not found")
			return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
		logger.WithError(err).Error("failed to get user by id")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if existingUser.AvatarUrl != "" {
		oldKey := extractKeyFromURL(existingUser.AvatarUrl)
		if oldKey != "" {
			if err := uc.awsRepo.RemoveObject(ctx, request.BucketName, oldKey); err != nil {
				logger.WithError(err).Warn("failed to remove old avatar, but continuing with upload")
			}
		}
	}

	uploadInput := dto.UploadAvatarRequestToUploadInputEntity(request)

	uploadInfo, err := uc.awsRepo.PutObject(ctx, *uploadInput)
	if err != nil {
		logger.WithError(err).Error("failed to upload avatar to storage")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	avatarURL := uc.generateAWSMinioURL(request.BucketName, uploadInfo.Key)

	existingUser.AvatarUrl = avatarURL
	if err := uc.postgresRepo.UpdateUser(ctx, existingUser); err != nil {
		logger.WithError(err).Error("failed to update user avatar URL")

		if cleanupErr := uc.awsRepo.RemoveObject(ctx, request.BucketName, uploadInfo.Key); cleanupErr != nil {
			logger.WithError(cleanupErr).Error("failed to cleanup uploaded avatar after user update failure")
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	response := &dto.UploadAvatarResponseDTO{
		AvatarUrl: avatarURL,
	}

	logger.WithField("avatar_url", avatarURL).Info("successfully uploaded avatar")
	return response, nil
}

func (uc *AuthUsecase) generateAWSMinioURL(bucket string, key string) string {
	return fmt.Sprintf("%s/%s/%s", uc.cfg.AWS.FilesEndpoint, bucket, key)
}

func extractKeyFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func (uc *AuthUsecase) VKOauthLink(ctx context.Context) (*dto.VKOauthLinkResponse, error) {
	const op = "AuthUsecase.OAuthLink"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	stateValue, err := generateState()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	codeChallenge := generateCodeChallenge(codeVerifier)

	if err := uc.redisRepo.StoreState(ctx, stateValue, codeVerifier); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	authURL := buildVKOauthURL(uc.cfg.Auth.VK.IntegrationID, uc.cfg.Auth.VK.RedirectURL, stateValue, codeChallenge)

	response := &dto.VKOauthLinkResponse{
		VKOauthURL: authURL,
	}

	logger.Info("successfully generated oauth link")
	return response, nil
}

func (uc *AuthUsecase) VKOAuthCallback(ctx context.Context, request *dto.VKCallbackRequestDTO) (*dto.UserTokenDTO, error) {
	const op = "AuthUsecase.VKOAuthCallback"
	logger := logctx.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":    op,
		"state": request.State,
	})

	codeVerifier, err := uc.redisRepo.GetCodeVerifierByState(ctx, request.State)
	if err != nil {
		logger.WithError(err).Error("failed to validate oauth state")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := uc.vkWebapi.ExchangeCodeForTokens(ctx, request.Code, codeVerifier, request.DeviceID)
	if err != nil {
		logger.WithError(err).Error("failed to exchange code for tokens")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userInfo, err := uc.vkWebapi.GetUserInfo(ctx, tokens.AccessToken)
	if err != nil {
		logger.WithError(err).Error("failed to get user info from vk")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := uc.postgresRepo.GetUserByEmail(ctx, userInfo.Email)
	if err != nil && !errs.IsNotFoundError(err) {
		logger.WithError(err).Error("failed to get user by email")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if user == nil {
		avatar := userInfo.Avatar
		if avatar == "" {
			avatar = uc.generateAWSMinioURL(uc.cfg.AWS.AvatarBucketName, "default.jpg")
		}
		newUser := &entities.User{
			Username:  userInfo.FirstName,
			Email:     userInfo.Email,
			AvatarUrl: avatar,
		}

		user, err = uc.postgresRepo.CreateUser(ctx, newUser)
		if err != nil {
			logger.WithError(err).Error("failed to create user from vk oauth")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		logger.WithField("user_id", user.ID.String()).Info("created new user from vk oauth")
	}

	existingVKUser, err := uc.postgresRepo.GetVKUserByUserID(ctx, user.ID)
	if err != nil && !errs.IsNotFoundError(err) {
		logger.WithError(err).Error("failed to check existing vk user")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	vkUserData := &entities.VKUser{
		UserID:       user.ID,
		VKUserID:     userInfo.VKUserID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
		DeviceID:     request.DeviceID,
	}

	if existingVKUser == nil {
		err = uc.postgresRepo.CreateVKUser(ctx, vkUserData)
		if err != nil {
			logger.WithError(err).Error("failed to create vk user data")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		logger.WithField("user_id", user.ID.String()).Info("created new vk user record")
	} else {
		vkUserData.ID = existingVKUser.ID
		err = uc.postgresRepo.UpdateVKUser(ctx, vkUserData)
		if err != nil {
			logger.WithError(err).Error("failed to update vk user data")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		logger.WithField("user_id", user.ID.String()).Info("updated existing vk user record")
	}

	if err := uc.redisRepo.DeleteState(ctx, request.State); err != nil {
		logger.WithError(err).Warn("failed to delete oauth state")
	}

	token, err := jwt.GenerateJWT(&uc.cfg.Auth.Jwt, user)
	if err != nil {
		logger.WithError(err).Error("failed to generate JWT token")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	response := dto.UserTokenToDTO(user, token)

	logger.WithField("user_id", user.ID.String()).Info("successfully processed vk oauth callback")
	return response, nil
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

func generateState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func generateCodeVerifier() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func buildVKOauthURL(clientID, redirectURI, state, codeChallenge string) string {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("state", state)
	params.Add("code_challenge", codeChallenge)
	params.Add("code_challenge_method", "s256")
	params.Add("scope", "email")

	return "https://id.vk.com/authorize?" + params.Encode()
}
