package main

import "time"

// Env var names
const (
	envPort        = "PORT"
	envAdminPass   = "ADMIN_PASSWORD"
	envJWTSecret   = "JWT_SECRET"
	envDatabaseURL = "DATABASE_URL"
)

// Error codes
const (
	errInvalidJSON   = "INVALID_JSON"
	errUnauthorized  = "UNAUTHORIZED"
	errInternalError = "INTERNAL_ERROR"
	errInvalidUUID   = "INVALID_UUID"
	errMissingField  = "MISSING_FIELD"
	errConflict      = "CONFLICT"
)

// Error messages
const (
	errMsgInvalidJSON        = "invalid request body"
	errMsgInvalidPass        = "invalid password"
	errMsgCreateToken        = "failed to create token"
	errMsgInvalidID          = "id must be a valid UUID v7"
	errMsgMissingTitle       = "title is required"
	errMsgMissingToken       = "missing token"
	errMsgInvalidToken       = "invalid or expired token"
	errMsgCreatePresentation = "failed to create presentation"
	errMsgDuplicateID       = "a presentation with this id already exists"
)

// Cookie
const (
	cookieName   = "token"
	cookiePath   = "/"
	cookieMaxAge = 86400 // 24h in seconds
)

// JWT
const (
	jwtSubject = "admin"
	jwtExpiry  = 24 * time.Hour
)

// Server
const defaultPort = "8080"

// Database
const dbTimeout = 5 * time.Second

// PostgreSQL error codes
const pgErrUniqueViolation = "23505"
