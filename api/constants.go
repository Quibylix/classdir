package main

import "time"

// Env var names
const (
	envPort      = "PORT"
	envAdminPass = "ADMIN_PASSWORD"
	envJWTSecret = "JWT_SECRET"
)

// Error codes
const (
	errInvalidJSON   = "INVALID_JSON"
	errUnauthorized  = "UNAUTHORIZED"
	errInternalError = "INTERNAL_ERROR"
)

// Error messages
const (
	errMsgInvalidJSON = "invalid request body"
	errMsgInvalidPass = "invalid password"
	errMsgCreateToken = "failed to create token"
)

// Cookie
const (
	cookieName    = "token"
	cookiePath    = "/"
	cookieMaxAge  = 86400 // 24h in seconds
)

// JWT
const (
	jwtSubject = "admin"
	jwtExpiry  = 24 * time.Hour
)

// Server
const defaultPort = "8080"
