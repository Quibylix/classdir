package presentation

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"classdir/api/internal/shared/cfg"
	"classdir/api/internal/shared/response"
)

type mockPresentationStore struct {
	createFunc func(ctx context.Context, id, title string) error
}

func (m *mockPresentationStore) Create(ctx context.Context, id, title string) error {
	return m.createFunc(ctx, id, title)
}

func TestCreatePresentation_ValidInput(t *testing.T) {
	var called bool
	store := &mockPresentationStore{
		createFunc: func(ctx context.Context, id, title string) error {
			called = true
			return nil
		},
	}

	handler := createPresentationHandler(store)
	body := `{"id":"0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f","title":"My Presentation"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusCreated)
	}
	if !called {
		t.Error("expected store.Create to be called")
	}
}

func TestCreatePresentation_InvalidUUID(t *testing.T) {
	store := &mockPresentationStore{
		createFunc: func(ctx context.Context, id, title string) error {
			t.Error("store.Create should not be called")
			return nil
		},
	}

	handler := createPresentationHandler(store)
	body := `{"id":"bad","title":"My Presentation"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}

	var payload response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected error JSON, got:", rec.Body.String())
	}
	if payload.Error.Code != cfg.ErrInvalidUUID {
		t.Errorf("got code %q, want %q", payload.Error.Code, cfg.ErrInvalidUUID)
	}
}

func TestCreatePresentation_EmptyTitle(t *testing.T) {
	store := &mockPresentationStore{
		createFunc: func(ctx context.Context, id, title string) error {
			t.Error("store.Create should not be called")
			return nil
		},
	}

	handler := createPresentationHandler(store)
	body := `{"id":"0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f","title":"  "}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}

	var payload response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected error JSON, got:", rec.Body.String())
	}
	if payload.Error.Code != cfg.ErrMissingField {
		t.Errorf("got code %q, want %q", payload.Error.Code, cfg.ErrMissingField)
	}
}

func TestCreatePresentation_InvalidJSON(t *testing.T) {
	store := &mockPresentationStore{
		createFunc: func(ctx context.Context, id, title string) error {
			t.Error("store.Create should not be called")
			return nil
		},
	}

	handler := createPresentationHandler(store)
	body := `{bad`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}

	var payload response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected error JSON, got:", rec.Body.String())
	}
	if payload.Error.Code != cfg.ErrInvalidJSON {
		t.Errorf("got code %q, want %q", payload.Error.Code, cfg.ErrInvalidJSON)
	}
}

func TestCreatePresentation_DuplicateID(t *testing.T) {
	store := &mockPresentationStore{
		createFunc: func(ctx context.Context, id, title string) error {
			return ErrDuplicateKey
		},
	}

	handler := createPresentationHandler(store)
	body := `{"id":"0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f","title":"My Presentation"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestCreatePresentation_StoreError(t *testing.T) {
	store := &mockPresentationStore{
		createFunc: func(ctx context.Context, id, title string) error {
			return errors.New("db error")
		},
	}

	handler := createPresentationHandler(store)
	body := `{"id":"0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f","title":"My Presentation"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}

	var payload response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected error JSON, got:", rec.Body.String())
	}
	if payload.Error.Code != cfg.ErrInternalError {
		t.Errorf("got code %q, want %q", payload.Error.Code, cfg.ErrInternalError)
	}
}
