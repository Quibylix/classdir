package auth

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"classdir/api/internal/shared/cfg"
	"classdir/api/internal/shared/response"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/login", loginHandler)
	mux.HandleFunc("POST /api/v1/auth/logout", logoutHandler)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidJSON, cfg.ErrMsgInvalidJSON)
		return
	}

	adminPass := os.Getenv(cfg.EnvAdminPass)
	if subtle.ConstantTimeCompare([]byte(body.Password), []byte(adminPass)) != 1 {
		response.WriteError(w, http.StatusUnauthorized, cfg.ErrUnauthorized, cfg.ErrMsgInvalidPass)
		return
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   cfg.JwtSubject,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(cfg.JwtExpiry)),
	}

	secret := os.Getenv(cfg.EnvJWTSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgCreateToken)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cfg.CookieName,
		Value:    signed,
		Path:     cfg.CookiePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   cfg.CookieMaxAge,
	})
	w.WriteHeader(http.StatusNoContent)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.CookieName,
		Value:    "",
		Path:     cfg.CookiePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
	w.WriteHeader(http.StatusNoContent)
}
