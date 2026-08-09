package cfg

import "time"

const (
	EnvPort        = "PORT"
	EnvAdminPass   = "ADMIN_PASSWORD"
	EnvJWTSecret   = "JWT_SECRET"
	EnvDatabaseURL = "DATABASE_URL"
	EnvWSOrigin    = "WS_ORIGIN"
)

const (
	ErrInvalidJSON   = "INVALID_JSON"
	ErrUnauthorized  = "UNAUTHORIZED"
	ErrInternalError = "INTERNAL_ERROR"
	ErrInvalidUUID   = "INVALID_UUID"
	ErrMissingField  = "MISSING_FIELD"
	ErrConflict      = "CONFLICT"
	ErrNotFound      = "NOT_FOUND"
	ErrRateLimit     = "RATE_LIMITED"
)

const (
	ErrMsgInvalidJSON          = "invalid request body"
	ErrMsgInvalidPass          = "invalid password"
	ErrMsgCreateToken          = "failed to create token"
	ErrMsgInvalidID            = "id must be a valid UUID v7"
	ErrMsgMissingTitle         = "title is required"
	ErrMsgMissingToken         = "missing token"
	ErrMsgInvalidToken         = "invalid or expired token"
	ErrMsgCreatePresentation   = "failed to create presentation"
	ErrMsgGetPresentation      = "failed to retrieve presentation"
	ErrMsgUpdatePresentation   = "failed to update presentation"
	ErrMsgDeletePresentation   = "failed to delete presentation"
	ErrMsgNotFound             = "presentation not found"
	ErrMsgDuplicateID          = "a presentation with this id already exists"
	ErrMsgListPresentation     = "failed to list presentations"
	ErrMsgMissingName          = "name is required"
	ErrMsgCreateStudent        = "failed to create student"
	ErrMsgGetStudent           = "failed to retrieve student"
	ErrMsgUpdateStudent        = "failed to update student"
	ErrMsgDeleteStudent        = "failed to delete student"
	ErrMsgListStudents         = "failed to list students"
	ErrMsgStudentNotFound      = "student not found"
	ErrMsgDuplicateStudentID   = "a student with this id already exists"
	ErrMsgDuplicateStudentName = "a student with this name already exists"
	ErrMsgUpdateContent        = "failed to update content"
	ErrMsgCreateRoom           = "failed to create room"
	ErrMsgRoomClosed           = "room is closed"
	ErrMsgRateLimit            = "too many requests"
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

const PgErrForeignKeyViolation = "23503"
