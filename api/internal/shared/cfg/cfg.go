package cfg

import "time"

const (
	EnvPort        = "PORT"
	EnvAdminPass   = "ADMIN_PASSWORD"
	EnvJWTSecret   = "JWT_SECRET"
	EnvDatabaseURL = "DATABASE_URL"
)

const (
	ErrInvalidJSON   = "INVALID_JSON"
	ErrUnauthorized  = "UNAUTHORIZED"
	ErrInternalError = "INTERNAL_ERROR"
	ErrInvalidUUID   = "INVALID_UUID"
	ErrMissingField  = "MISSING_FIELD"
	ErrConflict      = "CONFLICT"
	ErrNotFound      = "NOT_FOUND"
)

const (
	ErrMsgInvalidJSON        = "invalid request body"
	ErrMsgInvalidPass        = "invalid password"
	ErrMsgCreateToken        = "failed to create token"
	ErrMsgInvalidID          = "id must be a valid UUID v7"
	ErrMsgMissingTitle       = "title is required"
	ErrMsgMissingToken       = "missing token"
	ErrMsgInvalidToken       = "invalid or expired token"
	ErrMsgCreatePresentation = "failed to create presentation"
	ErrMsgGetPresentation   = "failed to retrieve presentation"
	ErrMsgNotFound          = "presentation not found"
	ErrMsgDuplicateID       = "a presentation with this id already exists"
)

const (
	CookieName   = "token"
	CookiePath   = "/"
	CookieMaxAge = 86400
)

const (
	JwtSubject = "admin"
	JwtExpiry  = 24 * time.Hour
)

const DefaultPort = "8080"

const DbTimeout = 5 * time.Second

const PgErrUniqueViolation = "23505"
