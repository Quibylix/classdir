package main

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, errInvalidJSON, errMsgInvalidJSON)
		return
	}

	adminPass := os.Getenv(envAdminPass)
	if subtle.ConstantTimeCompare([]byte(body.Password), []byte(adminPass)) != 1 {
		writeError(w, http.StatusUnauthorized, errUnauthorized, errMsgInvalidPass)
		return
	}

	secret := os.Getenv(envJWTSecret)
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   jwtSubject,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(jwtExpiry)),
	})

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		writeError(w, http.StatusInternalServerError, errInternalError, errMsgCreateToken)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    signed,
		Path:     cookiePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   cookieMaxAge,
	})

	w.WriteHeader(http.StatusNoContent)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			writeError(w, http.StatusUnauthorized, errUnauthorized, errMsgMissingToken)
			return
		}

		secret := os.Getenv(envJWTSecret)
		token, err := jwt.ParseWithClaims(cookie.Value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, errUnauthorized, errMsgInvalidToken)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     cookiePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusNoContent)
}
