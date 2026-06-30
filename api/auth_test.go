package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/auth/login", loginHandler)
	mux.HandleFunc("POST /api/v1/auth/logout", logoutHandler)
	return mux
}

func TestLoginHandler_ValidPassword(t *testing.T) {
	t.Setenv(envAdminPass, "testpass")
	t.Setenv(envJWTSecret, "testsecret")

	mux := setupMux()

	body := `{"password":"testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusNoContent)
	}

	cookies := rec.Result().Cookies()
	var token *http.Cookie
	for _, c := range cookies {
		if c.Name == cookieName {
			token = c
			break
		}
	}
	if token == nil {
		t.Fatal("expected token cookie, got none")
	}
	if token.Value == "" {
		t.Error("expected non-empty token value")
	}
	if !token.HttpOnly {
		t.Error("expected HttpOnly flag")
	}
	if !token.Secure {
		t.Error("expected Secure flag")
	}
	if token.SameSite != http.SameSiteStrictMode {
		t.Error("expected SameSite=Strict")
	}
	if token.MaxAge != cookieMaxAge {
		t.Errorf("got MaxAge %d, want %d", token.MaxAge, cookieMaxAge)
	}
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	t.Setenv(envAdminPass, "testpass")
	t.Setenv(envJWTSecret, "testsecret")

	mux := setupMux()

	body := `{"password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	var payload ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected error JSON, got:", rec.Body.String())
	}
	if payload.Error.Code != errUnauthorized {
		t.Errorf("got code %q, want %q", payload.Error.Code, errUnauthorized)
	}
	if payload.Error.Message == "" {
		t.Error("expected non-empty error message")
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	t.Setenv(envAdminPass, "testpass")
	t.Setenv(envJWTSecret, "testsecret")

	mux := setupMux()

	body := `{bad`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}

	var payload ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected error JSON, got:", rec.Body.String())
	}
	if payload.Error.Code != errInvalidJSON {
		t.Errorf("got code %q, want %q", payload.Error.Code, errInvalidJSON)
	}
	if payload.Error.Message == "" {
		t.Error("expected non-empty error message")
	}
}

func TestLoginHandler_EmptyBody(t *testing.T) {
	t.Setenv(envAdminPass, "testpass")
	t.Setenv(envJWTSecret, "testsecret")

	mux := setupMux()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestLogoutHandler(t *testing.T) {
	mux := setupMux()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusNoContent)
	}

	cookies := rec.Result().Cookies()
	var token *http.Cookie
	for _, c := range cookies {
		if c.Name == cookieName {
			token = c
			break
		}
	}
	if token == nil {
		t.Fatal("expected token cookie, got none")
	}
	if token.MaxAge != -1 {
		t.Errorf("got MaxAge %d, want -1", token.MaxAge)
	}
	if token.Value != "" {
		t.Error("expected empty token value")
	}
}
