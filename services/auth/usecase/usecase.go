package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/friendclient"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/pkg/jwt"
	"github.com/F0urward/proftwist-backend/services/auth"
	"github.com/F0urward/proftwist-backend/services/auth/dto"
)

type AuthUsecase struct {
	cfg          *config.Config
	postgresRepo auth.PostgresRepository
	redisRepo    auth.RedisRepository
	awsRepo      auth.AWSRepository
	vkWebapi     auth.VKWebapi
	friendClient friendclient.FriendServiceClient
}

func NewAuthUsecase(postgresRepo auth.PostgresRepository, redisRepo auth.RedisRepository, awsRepo auth.AWSRepository, vkWebapi auth.VKWebapi, friendClient friendclient.FriendServiceClient, cfg *config.Config) auth.Usecase {
	return &AuthUsecase{
		cfg:          cfg,
		postgresRepo: postgresRepo,
		redisRepo:    redisRepo,
		awsRepo:      awsRepo,
		vkWebapi:     vkWebapi,
		friendClient: friendClient,
	}
}

func (uc *AuthUsecase) Register(ctx context.Context, request *dto.RegisterRequestDTO) (*dto.UserTokenDTO, error) {
	const op = "AuthUsecase.Register"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":       op,
		"email":    request.Email,
		"username": request.Username,
	})

	existingUser, err := uc.postgresRepo.GetUserByEmail(ctx, request.Email)
	if err != nil && !errs.IsNotFoundError(err) {
		logger.WithError(err).Error("failed to check existing user")
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		logger.Warn("user with this email already exists")
		return nil, errs.ErrAlreadyExists
	}

	passwordHash, err := hashPassword(request.Password)
	if err != nil {
		logger.WithError(err).Error("failed to hash password")
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := dto.RegisterRequestToEntity(request, passwordHash)

	// newUser.AvatarUrl = uc.generateAWSMinioURL(uc.cfg.AWS.AvatarBucketName, "default.jpg")
	newUser.AvatarUrl = ""

	createdUser, err := uc.postgresRepo.CreateUser(ctx, newUser)
	if err != nil {
		logger.WithError(err).Error("failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := jwt.GenerateJWT(&uc.cfg.Auth.Jwt, createdUser)
	if err != nil {
		logger.WithError(err).Error("failed to generate JWT token")
		return nil, fmt.Errorf("failed to generate authentication token: %w", err)
	}

	response := dto.UserTokenToDTO(createdUser, token)

	logger.WithField("user_id", createdUser.ID.String()).Info("user registered successfully")
	return response, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, request *dto.LoginRequestDTO) (*dto.UserTokenDTO, error) {
	const op = "AuthUsecase.Login"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":    op,
		"email": request.Email,
	})

	user, err := uc.postgresRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		if errs.IsNotFoundError(err) {
			logger.Warn("user not found")
			return nil, errs.ErrInvalidCredentials
		}
		logger.WithError(err).Error("failed to get user by email")
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	err = checkPasswordHash(request.Password, user.PasswordHash)
	if err != nil {
		logger.WithError(err).Warn("invalid password")
		return nil, errs.ErrInvalidCredentials
	}

	token, err := jwt.GenerateJWT(&uc.cfg.Auth.Jwt, user)
	if err != nil {
		logger.WithError(err).Error("failed to generate JWT token")
		return nil, fmt.Errorf("failed to generate authentication token: %w", err)
	}

	response := dto.UserTokenToDTO(user, token)

	logger.WithField("user_id", user.ID.String()).Info("user logged in successfully")
	return response, nil
}

func (uc *AuthUsecase) Logout(ctx context.Context, token string) error {
	const op = "AuthUsecase.Logout"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	claims, err := jwt.ParseJWT(&uc.cfg.Auth.Jwt, token)
	if err != nil {
		logger.WithError(err).Error("failed to parse token")
		return errs.ErrInvalidToken
	}

	if err := uc.redisRepo.AddToBlacklist(ctx, claims.UserID, token); err != nil {
		logger.WithError(err).Error("failed to add token to blacklist")
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}

	logger.Info("user logged out successfully")
	return nil
}

func (uc *AuthUsecase) GetMe(ctx context.Context, userID uuid.UUID) (*dto.UserDTO, error) {
	const op = "AuthUsecase.GetMe"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID,
	})

	user, err := uc.postgresRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errs.IsNotFoundError(err) {
			logger.WithError(err).Warn("user not found")
			return nil, errs.ErrNotFound
		}
		logger.WithError(err).Error("failed to get user by id")
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	userDTO := dto.UserToDTO(user)

	logger.Info("successfully retrieved user info")
	return &userDTO, nil
}

func (uc *AuthUsecase) GetByID(ctx context.Context, userID uuid.UUID) (*dto.GetUserByIDResponseDTO, error) {
	const op = "AuthUsecase.GetByID"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID.String(),
	})

	user, err := uc.postgresRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errs.IsNotFoundError(err) {
			logger.WithError(err).Warn("user not found")
			return nil, errs.ErrNotFound
		}
		logger.WithError(err).Error("failed to get user by id")
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	userDTO := dto.UserToDTO(user)
	response := &dto.GetUserByIDResponseDTO{
		User: userDTO,
	}

	logger.Info("successfully retrieved user by ID")
	return response, nil
}

func (uc *AuthUsecase) GetByIDs(ctx context.Context, userIDs []uuid.UUID) (*dto.GetUsersByIDsResponseDTO, error) {
	const op = "AuthUsecase.GetByIDs"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":         op,
		"user_ids":   userIDs,
		"user_count": len(userIDs),
	})

	if len(userIDs) == 0 {
		logger.Warn("empty user IDs list provided")
		return &dto.GetUsersByIDsResponseDTO{Users: []dto.UserDTO{}}, nil
	}

	users, err := uc.postgresRepo.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		logger.WithError(err).Error("failed to get users by IDs")
		return nil, fmt.Errorf("failed to get users by IDs: %w", err)
	}

	userDTOs := dto.UserListToDTO(users)

	response := &dto.GetUsersByIDsResponseDTO{
		Users: userDTOs,
	}

	logger.WithField("found_count", len(userDTOs)).Info("successfully retrieved users by IDs")
	return response, nil
}

func (uc *AuthUsecase) Update(ctx context.Context, userID uuid.UUID, request *dto.UpdateUserRequestDTO) error {
	const op = "AuthUsecase.Update"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID.String(),
	})

	existingUser, err := uc.postgresRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errs.IsNotFoundError(err) {
			logger.WithError(err).Warn("user not found")
			return errs.ErrNotFound
		}
		logger.WithError(err).Error("failed to get user by id")
		return fmt.Errorf("failed to get user by id: %w", err)
	}

	if request.Email != "" && request.Email != existingUser.Email {
		userWithEmail, err := uc.postgresRepo.GetUserByEmail(ctx, request.Email)
		if err != nil && !errs.IsNotFoundError(err) {
			logger.WithError(err).Error("failed to check email availability")
			return fmt.Errorf("failed to check email availability: %w", err)
		}

		if userWithEmail != nil && userWithEmail.ID != userID {
			logger.WithField("email", request.Email).Warn("email already taken by another user")
			return errs.ErrAlreadyExists
		}
	}

	updatedUser := dto.UpdateUserRequestToEntity(existingUser, request)

	err = uc.postgresRepo.UpdateUser(ctx, updatedUser)
	if err != nil {
		logger.WithError(err).Error("failed to update user")
		return fmt.Errorf("failed to update user: %w", err)
	}

	logger.Info("successfully updated user")
	return nil
}

func (uc *AuthUsecase) UploadAvatar(ctx context.Context, request *dto.UploadAvatarRequestDTO) (*dto.UploadAvatarResponseDTO, error) {
	const op = "AuthUsecase.UploadAvatar"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": request.UserID.String(),
	})

	existingUser, err := uc.postgresRepo.GetUserByID(ctx, request.UserID)
	if err != nil {
		if errs.IsNotFoundError(err) {
			logger.WithError(err).Warn("user not found")
			return nil, errs.ErrNotFound
		}
		logger.WithError(err).Error("failed to get user by id")
		return nil, fmt.Errorf("failed to get user by id: %w", err)
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
		return nil, fmt.Errorf("failed to upload avatar to storage: %w", err)
	}

	avatarURL := uc.generateAWSMinioURL(request.BucketName, uploadInfo.Key)

	existingUser.AvatarUrl = avatarURL
	if err := uc.postgresRepo.UpdateUser(ctx, existingUser); err != nil {
		logger.WithError(err).Error("failed to update user avatar URL")

		if cleanupErr := uc.awsRepo.RemoveObject(ctx, request.BucketName, uploadInfo.Key); cleanupErr != nil {
			logger.WithError(cleanupErr).Error("failed to cleanup uploaded avatar after user update failure")
		}

		return nil, fmt.Errorf("failed to update user avatar URL: %w", err)
	}

	response := &dto.UploadAvatarResponseDTO{
		AvatarUrl: avatarURL,
	}

	logger.WithField("avatar_url", avatarURL).Info("successfully uploaded avatar")
	return response, nil
}

func (uc *AuthUsecase) IsInBlacklist(ctx context.Context, userID string, token string) (bool, error) {
	const op = "AuthUsecase.IsInBlacklist"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"user_id": userID,
	})

	isBlacklisted, err := uc.redisRepo.IsInBlacklist(ctx, userID, token)
	if err != nil {
		logger.WithError(err).Error("failed to check token blacklist status")
		return false, fmt.Errorf("failed to check token blacklist status: %w", err)
	}

	logger.WithField("is_blacklisted", isBlacklisted).Info("token blacklist status checked")
	return isBlacklisted, nil
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

func (uc *AuthUsecase) SearchUsers(ctx context.Context, userID uuid.UUID, query string) (*dto.SearchUsersResponseDTO, error) {
	const op = "AuthUsecase.SearchUsers"
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":      op,
		"query":   query,
		"user_id": userID.String(),
	})

	if query == "" {
		return &dto.SearchUsersResponseDTO{Users: []dto.UserPublicDTO{}}, nil
	}

	users, err := uc.postgresRepo.SearchUsers(ctx, query)
	if err != nil {
		logger.WithError(err).Error("failed to search users")
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	if users == nil {
		users = []*entities.User{}
	}

	if len(users) == 0 {
		logger.Info("no users found for search query")
		return &dto.SearchUsersResponseDTO{Users: []dto.UserPublicDTO{}}, nil
	}

	friendshipStatusMap := uc.fetchFriendshipStatus(ctx, userID, users)
	userPublicDTOs := dto.UserListToPublicDTO(users, friendshipStatusMap)

	logger.WithField("count", len(userPublicDTOs)).Info("successfully searched users")
	return &dto.SearchUsersResponseDTO{Users: userPublicDTOs}, nil
}

func (uc *AuthUsecase) fetchFriendshipStatus(ctx context.Context, currentUserID uuid.UUID, users []*entities.User) map[uuid.UUID]*dto.FriendshipStatusDTO {
	if len(users) == 0 {
		return make(map[uuid.UUID]*dto.FriendshipStatusDTO)
	}

	statusMap := make(map[uuid.UUID]*dto.FriendshipStatusDTO, len(users))

	for _, user := range users {
		if user == nil || user.ID == currentUserID {
			continue
		}

		friendshipStatus, err := uc.friendClient.GetFriendshipStatus(ctx, &friendclient.GetFriendshipStatusRequest{
			UserId:       currentUserID.String(),
			TargetUserId: user.ID.String(),
		})
		if err != nil {
			continue
		}

		if friendshipStatus != nil {
			statusMap[user.ID] = &dto.FriendshipStatusDTO{
				Status:    friendshipStatus.Status,
				RequestID: friendshipStatus.RequestId,
				IsSender:  friendshipStatus.IsSender,
			}
		}
	}

	return statusMap
}

func (uc *AuthUsecase) VKOauthLink(ctx context.Context) (*dto.VKOauthLinkResponse, error) {
	const op = "AuthUsecase.OAuthLink"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	stateValue, err := generateState()
	if err != nil {
		logger.WithError(err).Error("failed to generate state")
		return nil, fmt.Errorf("failed to generate oauth state: %w", err)
	}

	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		logger.WithError(err).Error("failed to generate code verifier")
		return nil, fmt.Errorf("failed to generate oauth code verifier: %w", err)
	}
	codeChallenge := generateCodeChallenge(codeVerifier)

	if err := uc.redisRepo.StoreState(ctx, stateValue, codeVerifier); err != nil {
		logger.WithError(err).Error("failed to store state")
		return nil, fmt.Errorf("failed to store state: %w", err)
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
	logger := ctxutil.GetLogger(ctx).WithFields(map[string]interface{}{
		"op":    op,
		"state": request.State,
	})

	codeVerifier, err := uc.redisRepo.GetCodeVerifierByState(ctx, request.State)
	if err != nil {
		logger.WithError(err).Error("failed to validate oauth state")
		return nil, errs.ErrInvalidOAuthState
	}

	tokens, err := uc.vkWebapi.ExchangeCodeForTokens(ctx, request.Code, codeVerifier, request.DeviceID)
	if err != nil {
		logger.WithError(err).Error("failed to exchange code for tokens")
		return nil, fmt.Errorf("failed to exchange oauth code for tokens: %w", err)
	}

	userInfo, err := uc.vkWebapi.GetUserInfo(ctx, tokens.AccessToken)
	if err != nil {
		logger.WithError(err).Error("failed to get user info from vk")
		return nil, fmt.Errorf("failed to get user information from VK: %w", err)
	}

	user, err := uc.postgresRepo.GetUserByEmail(ctx, userInfo.Email)
	if err != nil && !errs.IsNotFoundError(err) {
		logger.WithError(err).Error("failed to get user by email")
		return nil, fmt.Errorf("failed to get user by email: %w", err)
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
			return nil, fmt.Errorf("failed to create user from vk oauth: %w", err)
		}

		logger.WithField("user_id", user.ID.String()).Info("created new user from vk oauth")
	}

	existingVKUser, err := uc.postgresRepo.GetVKUserByUserID(ctx, user.ID)
	if err != nil && !errs.IsNotFoundError(err) {
		logger.WithError(err).Error("failed to check existing vk user")
		return nil, fmt.Errorf("failed to check existing vk user: %w", err)
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
			return nil, fmt.Errorf("failed to create vk user data: %w", err)
		}
		logger.WithField("user_id", user.ID.String()).Info("created new vk user record")
	} else {
		vkUserData.ID = existingVKUser.ID
		err = uc.postgresRepo.UpdateVKUser(ctx, vkUserData)
		if err != nil {
			logger.WithError(err).Error("failed to update vk user data")
			return nil, fmt.Errorf("failed to update vk user data: %w", err)
		}
		logger.WithField("user_id", user.ID.String()).Info("updated existing vk user record")
	}

	if err := uc.redisRepo.DeleteState(ctx, request.State); err != nil {
		logger.WithError(err).Warn("failed to delete oauth state")
	}

	token, err := jwt.GenerateJWT(&uc.cfg.Auth.Jwt, user)
	if err != nil {
		logger.WithError(err).Error("failed to generate JWT token")
		return nil, fmt.Errorf("failed to generate authentication token: %w", err)
	}

	response := dto.UserTokenToDTO(user, token)

	logger.WithField("user_id", user.ID.String()).Info("successfully processed vk oauth callback")
	return response, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func generateState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes for state: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func generateCodeVerifier() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes for code verifier: %w", err)
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
