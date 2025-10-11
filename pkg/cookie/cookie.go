package cookie

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/config"
)

type CookieProvider struct {
	cfg *config.JwtCookieConfig
}

func NewCookieProvider(cfg *config.JwtCookieConfig) *CookieProvider {
	return &CookieProvider{
		cfg: cfg,
	}
}

func (p *CookieProvider) GetAuthTokenCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(p.cfg.Name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (p *CookieProvider) SetAuthTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     p.cfg.Name,
		Value:    token,
		Path:     "/",
		HttpOnly: p.cfg.HttpOnly,
		Secure:   p.cfg.Secure,
		MaxAge:   p.cfg.MaxAge,
	})
}

func (p *CookieProvider) ClearAuthTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     p.cfg.Name,
		Value:    "",
		Path:     "/",
		HttpOnly: p.cfg.HttpOnly,
		Secure:   p.cfg.Secure,
		MaxAge:   -1,
	})
}
