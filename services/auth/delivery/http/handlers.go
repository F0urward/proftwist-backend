package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/pkg/cookie"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/pkg/image"
	"github.com/F0urward/proftwist-backend/services/auth"
	"github.com/F0urward/proftwist-backend/services/auth/dto"
)

type AuthHandlers struct {
	cfg *config.Config
	uc  auth.Usecase
}

func NewAuthHandlers(authUC auth.Usecase, cfg *config.Config) auth.Handlers {
	return &AuthHandlers{
		cfg: cfg,
		uc:  authUC,
	}
}

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.Register"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	var req dto.RegisterRequestDTO

	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Username == "" || req.Password == "" {
		logger.Warn("email, username and password required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "email, username and password required")
		return
	}

	logger = logger.WithField("email", req.Email)

	res, err := h.uc.Register(r.Context(), &req)
	if err != nil {
		logger.WithError(err).Error("failed to register user")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to register user"

		switch {
		case errs.IsAlreadyExistsError(err):
			statusCode = http.StatusConflict
			errorMsg = "email already registered"

		case errs.IsBusinessLogicError(err):
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	cookieProvider := cookie.NewCookieProvider(&h.cfg.Auth.Jwt.Cookie)
	cookieProvider.SetAuthTokenCookie(w, res.Token)

	logger.WithField("user_id", res.User.ID.String()).Info("successfully registered user")
	utils.JSONResponse(r.Context(), w, http.StatusCreated, res)
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.Login"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	var req dto.LoginRequestDTO

	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "email and password are required")
		return
	}

	logger = logger.WithField("email", req.Email)

	res, err := h.uc.Login(r.Context(), &req)
	if err != nil {
		logger.WithError(err).Error("failed to login user")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to login user"

		if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		} else if errs.IsNotFoundError(err) {
			statusCode = http.StatusUnauthorized
			errorMsg = "invalid credentials"
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	cookieProvider := cookie.NewCookieProvider(&h.cfg.Auth.Jwt.Cookie)
	cookieProvider.SetAuthTokenCookie(w, res.Token)

	logger.WithField("user_id", res.User.ID.String()).Info("successfully logged in user")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.Logout"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	cookieProvider := cookie.NewCookieProvider(&h.cfg.Auth.Jwt.Cookie)
	token, err := cookieProvider.GetAuthTokenCookie(r)
	if err != nil {
		logger.WithError(err).Warn("failed to extract token from cookie")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "failed to extract token from cookie")
		return
	}

	err = h.uc.Logout(r.Context(), token)
	if err != nil {
		logger.WithError(err).Error("failed to logout user")
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to logout user")
		return
	}

	cookieProvider.ClearAuthTokenCookie(w)

	logger.Info("successfully logged out user")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandlers) GetMe(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.GetMe"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	logger = logger.WithField("user_id", userID.String())

	user, err := h.uc.GetMe(r.Context(), userID)
	if err != nil {
		logger.WithError(err).Error("failed to get user info")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get user info"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "user not found"
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	response := dto.GetMeResponseDTO{
		User: *user,
	}

	logger.Info("successfully retrieved user info")
	utils.JSONResponse(r.Context(), w, http.StatusOK, response)
}

func (h *AuthHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.GetByID"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	if userIDStr == "" {
		logger.Warn("user_id parameter is required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "user_id parameter is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	logger = logger.WithField("user_id", userID.String())

	res, err := h.uc.GetByID(r.Context(), userID)
	if err != nil {
		logger.WithError(err).Error("failed to get user by ID")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to get user"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "user not found"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Info("successfully retrieved user")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *AuthHandlers) Update(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.Update"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	var req dto.UpdateUserRequestDTO

	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Warn("invalid request body")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	logger = logger.WithField("user_id", userID.String())

	err = h.uc.Update(r.Context(), userID, &req)
	if err != nil {
		logger.WithError(err).Error("failed to update user")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to update user"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "user not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Info("successfully updated user")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandlers) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.UploadAvatar"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		logger.WithError(err).Warn("failed to get avatar file from form")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "avatar file is required")
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.WithError(err).Warn("failed to close avatar file")
		}
	}()

	if header.Size > h.cfg.Upload.Avatar.MaxSize {
		logger.WithFields(map[string]interface{}{
			"file_size": header.Size,
			"max_size":  h.cfg.Upload.Avatar.MaxSize,
		}).Warn("avatar file too large")
		utils.JSONError(r.Context(), w, http.StatusBadRequest,
			fmt.Sprintf("avatar file too large. Maximum size is %d MB", h.cfg.Upload.Avatar.MaxSize/(1024*1024)))
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		logger.WithError(err).Warn("failed to read file data")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "failed to read file")
		return
	}

	contentType, err := image.CheckImageFileContentType(fileData)
	if err != nil {
		logger.WithError(err).Warn("invalid image content type")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid image format")
		return
	}
	reader := bytes.NewReader(fileData)

	uploadInput := dto.UploadAvatarRequestDTO{
		UserID:      userID,
		File:        reader,
		Name:        header.Filename,
		Size:        header.Size,
		ContentType: contentType,
		BucketName:  h.cfg.AWS.AvatarBucketName,
	}

	logger = logger.WithFields(map[string]interface{}{
		"user_id":      userID.String(),
		"filename":     header.Filename,
		"size":         header.Size,
		"content_type": contentType,
		"bucket":       uploadInput.BucketName,
	})

	res, err := h.uc.UploadAvatar(r.Context(), &uploadInput)
	if err != nil {
		logger.WithError(err).Error("failed to upload avatar")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to upload avatar"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "user not found"
		} else if errs.IsForbiddenError(err) {
			statusCode = http.StatusForbidden
			errorMsg = "access denied"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.Info("successfully uploaded avatar")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *AuthHandlers) SearchUsers(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.SearchUsers"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	userIDStr, ok := r.Context().Value(utils.UserIDKey{}).(string)
	if !ok || userIDStr == "" {
		logger.Warn("user ID not found in context")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).WithField("user_id", userIDStr).Warn("invalid user_id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	query := r.URL.Query().Get("q")

	logger = logger.WithField("query", query)

	res, err := h.uc.SearchUsers(r.Context(), userID, query)
	if err != nil {
		logger.WithError(err).Error("failed to search users")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to search users"

		if errs.IsNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorMsg = "no users found for this search query"
		} else if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	logger.WithField("count", len(res.Users)).Info("successfully searched users")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}

func (h *AuthHandlers) VKOauthLink(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.OAuthLink"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	response, err := h.uc.VKOauthLink(r.Context())
	if err != nil {
		logger.WithError(err).Error("failed to generate oauth link")
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to generate oauth link")
		return
	}

	logger.Info("successfully generated oauth link")
	utils.JSONResponse(r.Context(), w, http.StatusOK, response)
}

func (h *AuthHandlers) VKOAuthCallback(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.VKOAuthCallback"
	logger := ctxutil.GetLogger(r.Context()).WithField("op", op)

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	deviceID := r.URL.Query().Get("device_id")

	if code == "" || state == "" {
		logger.Warn("code and state are required")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "code and state are required")
		return
	}

	request := dto.VKCallbackRequestDTO{
		Code:     code,
		State:    state,
		DeviceID: deviceID,
	}

	res, err := h.uc.VKOAuthCallback(r.Context(), &request)
	if err != nil {
		logger.WithError(err).Error("failed to process vk oauth callback")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to process vk oauth callback"

		if errs.IsBusinessLogicError(err) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		} else if errs.IsNotFoundError(err) {
			statusCode = http.StatusUnauthorized
			errorMsg = "invalid oauth state"
		}

		utils.JSONError(r.Context(), w, statusCode, errorMsg)
		return
	}

	cookieProvider := cookie.NewCookieProvider(&h.cfg.Auth.Jwt.Cookie)
	cookieProvider.SetAuthTokenCookie(w, res.Token)

	logger.WithField("user_id", res.User.ID.String()).Info("successfully processed vk oauth callback")
	utils.JSONResponse(r.Context(), w, http.StatusOK, res)
}
