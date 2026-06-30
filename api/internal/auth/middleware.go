package auth

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"

	"classdir/api/internal/shared/cfg"
	"classdir/api/internal/shared/response"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cfg.CookieName)
		if err != nil {
			response.WriteError(w, http.StatusUnauthorized, cfg.ErrUnauthorized, cfg.ErrMsgMissingToken)
			return
		}

		secret := os.Getenv(cfg.EnvJWTSecret)
		token, err := jwt.ParseWithClaims(cookie.Value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			response.WriteError(w, http.StatusUnauthorized, cfg.ErrUnauthorized, cfg.ErrMsgInvalidToken)
			return
		}

		next.ServeHTTP(w, r)
	})
}
