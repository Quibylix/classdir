package hub

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"

	"classdir/api/internal/shared/cfg"
)

func WSHandler(hub *Hub, acceptor wsAcceptor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := acceptor.Accept(w, r, nil)
		if err != nil {
			return
		}
		conn.SetReadLimit(maxMessageSize)
		client := NewClient(hub, conn)

		if cookie, err := r.Cookie(cfg.CookieName); err == nil {
			secret := os.Getenv(cfg.EnvJWTSecret)
			token, err := jwt.ParseWithClaims(cookie.Value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
				return []byte(secret), nil
			})
			if err == nil && token.Valid {
				client.Authenticated = true
			}
		}

		go client.WritePump()
		go client.ReadPump()
	}
}
