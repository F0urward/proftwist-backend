package http

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/entities/errs"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/internal/utils"
	"github.com/F0urward/proftwist-backend/pkg/cookie"
	"github.com/F0urward/proftwist-backend/services/auth"
	"github.com/F0urward/proftwist-backend/services/auth/dto"

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

func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandlers.Logout"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	err := h.uc.Logout(r.Context())
	if err != nil {
		logger.WithError(err).Error("failed to logout user")
		utils.JSONError(r.Context(), w, http.StatusInternalServerError, "failed to logout user")
		return
	}

	cookieProvider := cookie.NewCookieProvider(&h.cfg.Auth.Jwt.Cookie)
	cookieProvider.ClearAuthTokenCookie(w)

	logger.Info("successfully logged out user")
	w.WriteHeader(http.StatusOK)
}
