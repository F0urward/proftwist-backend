package http

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/pkg/cookie"
	"github.com/F0urward/proftwist-backend/pkg/jwt"
	"github.com/F0urward/proftwist-backend/services/auth"
	"github.com/F0urward/proftwist-backend/services/auth/dto"
	"github.com/google/uuid"

	"github.com/mailru/easyjson"
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
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

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

	if req.Role == "" {
		req.Role = "regular"
	}

	logger = logger.WithField("email", req.Email)

	res, err := h.uc.Register(r.Context(), &req)
	if err != nil {
		logger.WithError(err).Error("failed to register user")

		statusCode := http.StatusInternalServerError
		errorMsg := "failed to register user"

		if errs.IsBusinessLogicError(err) || errs.IsAlreadyExistsError(err) {
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
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

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

func (h *AuthHandlers) GetMe(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.GetMe"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	cookieProvider := cookie.NewCookieProvider(&h.cfg.Auth.Jwt.Cookie)
	tokenString, err := cookieProvider.GetAuthTokenCookie(r)
	if err != nil {
		logger.WithError(err).Warn("failed to extract token from cookie")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "authentication required")
		return
	}

	claims, err := jwt.ParseJWT(&h.cfg.Auth.Jwt, tokenString)
	if err != nil {
		logger.WithError(err).Warn("invalid token")
		utils.JSONError(r.Context(), w, http.StatusUnauthorized, "invalid token")
		return
	}

	logger = logger.WithField("user_id", claims.UserID)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		logger.WithError(err).Warn("invalid user id format")
		utils.JSONError(r.Context(), w, http.StatusBadRequest, "invalid user_id format")
		return
	}

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

func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.Logout"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

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

func (h *AuthHandlers) VKOauthLink(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.OAuthLink"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

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
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

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
