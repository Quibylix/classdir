package hub

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"

	"classdir/api/internal/shared/cfg"
)

func WSHandler(hub *Hub, acceptor wsAcceptor, rlp rateLimitProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := acceptor.Accept(w, r, nil)
		if err != nil {
			return
		}
		conn.SetReadLimit(maxMessageSize)
		client := NewClient(hub, conn)

		authenticated := false
		if cookie, err := r.Cookie(cfg.CookieName); err == nil {
			secret := os.Getenv(cfg.EnvJWTSecret)
			token, err := jwt.ParseWithClaims(cookie.Value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
				return []byte(secret), nil
			})
			if err == nil && token.Valid {
				authenticated = true
			}
		}
		client.Authenticated = authenticated

		limit, burst := rlp.Limits(client.Authenticated)
		client.limiter = rate.NewLimiter(limit, burst)

		go client.WritePump()
		go client.ReadPump()
	}
}
