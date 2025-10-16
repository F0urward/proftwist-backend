package cors

import (
	"fmt"
	"net/http"

	"github.com/F0urward/proftwist-backend/config"
)

type CORSMiddleware struct {
	cfg *config.Config
}

func NewCORSMiddleware(cfg *config.Config) *CORSMiddleware {
	return &CORSMiddleware{
		cfg: cfg,
	}
}

func (c *CORSMiddleware) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", c.cfg.Service.CORS.AllowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", c.cfg.Service.CORS.AllowMethods)
		w.Header().Set("Access-Control-Allow-Headers", c.cfg.Service.CORS.AllowHeaders)
		w.Header().Set("Access-Control-Expose-Headers", c.cfg.Service.CORS.ExposeHeaders)
		w.Header().Set("Access-Control-Allow-Credentials", fmt.Sprintf("%t", c.cfg.Service.CORS.AllowCredentials))

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
