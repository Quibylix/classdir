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

const validPresID = "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f"
const validSlideID = "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80"

func TestCreateSlide_Valid(t *testing.T) {
	var called bool
	store := &mockPresentationStore{
		createSlideFunc: func(ctx context.Context, presID, slideID, content string) error {
			called = true
			if presID != validPresID {
				t.Errorf("got presID %q, want %q", presID, validPresID)
			}
			if slideID != validSlideID {
				t.Errorf("got slideID %q, want %q", slideID, validSlideID)
			}
			if strings.Contains(content, "<script") {
				t.Error("content should not contain script tags after sanitization")
			}
			if !strings.Contains(content, "<h1>") {
				t.Error("content should contain allowed tags after sanitization")
			}
			return nil
		},
	}

	handler := createSlideHandler(store)
	body := `{"id":"` + validSlideID + `","content":"<script>alert('xss')</script><h1>Hello</h1>"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusCreated)
	}
	if !called {
		t.Error("expected store.CreateSlide to be called")
	}

	var payload struct {
		Data Slide `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected valid JSON, got:", rec.Body.String())
	}
	if payload.Data.ID != validSlideID {
		t.Errorf("got id %q, want %q", payload.Data.ID, validSlideID)
	}
	if strings.Contains(payload.Data.Content, "<script") {
		t.Error("response content should not contain script tags after sanitization")
	}
	if !strings.Contains(payload.Data.Content, "<h1>") {
		t.Error("response content should contain allowed tags after sanitization")
	}
}

func TestCreateSlide_InvalidSlideID(t *testing.T) {
	store := &mockPresentationStore{
		createSlideFunc: func(ctx context.Context, presID, slideID, content string) error {
			t.Error("store.CreateSlide should not be called")
			return nil
		},
	}

	handler := createSlideHandler(store)
	body := `{"id":"bad","content":"<h1>Hello</h1>"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresID)
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

func TestCreateSlide_InvalidPresID(t *testing.T) {
	store := &mockPresentationStore{
		createSlideFunc: func(ctx context.Context, presID, slideID, content string) error {
			t.Error("store.CreateSlide should not be called")
			return nil
		},
	}

	handler := createSlideHandler(store)
	body := `{"id":"` + validSlideID + `","content":"<h1>Hello</h1>"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, "bad")
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

func TestCreateSlide_EmptyContent(t *testing.T) {
	store := &mockPresentationStore{
		createSlideFunc: func(ctx context.Context, presID, slideID, content string) error {
			t.Error("store.CreateSlide should not be called")
			return nil
		},
	}

	handler := createSlideHandler(store)
	body := `{"id":"` + validSlideID + `","content":"  "}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresID)
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

func TestCreateSlide_StoreError(t *testing.T) {
	store := &mockPresentationStore{
		createSlideFunc: func(ctx context.Context, presID, slideID, content string) error {
			return errors.New("db error")
		},
	}

	handler := createSlideHandler(store)
	body := `{"id":"` + validSlideID + `","content":"<h1>Hello</h1>"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresID)
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

func TestCreateSlide_NotFound(t *testing.T) {
	store := &mockPresentationStore{
		createSlideFunc: func(ctx context.Context, presID, slideID, content string) error {
			return ErrNotFound
		},
	}

	handler := createSlideHandler(store)
	body := `{"id":"` + validSlideID + `","content":"<h1>Hello</h1>"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestGetSlide_Found(t *testing.T) {
	store := &mockPresentationStore{
		getSlideFunc: func(ctx context.Context, presID, slideID string) (*Slide, error) {
			return &Slide{ID: slideID, Content: "<h1>Hi</h1>"}, nil
		},
	}

	handler := getSlideHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, validSlideID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data Slide `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected valid JSON, got:", rec.Body.String())
	}
	if payload.Data.ID != validSlideID {
		t.Errorf("got id %q, want %q", payload.Data.ID, validSlideID)
	}
	if payload.Data.Content != "<h1>Hi</h1>" {
		t.Errorf("got content %q, want %q", payload.Data.Content, "<h1>Hi</h1>")
	}
}

func TestGetSlide_NotFound(t *testing.T) {
	store := &mockPresentationStore{
		getSlideFunc: func(ctx context.Context, presID, slideID string) (*Slide, error) {
			return nil, nil
		},
	}

	handler := getSlideHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, validSlideID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}

	var payload response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected error JSON, got:", rec.Body.String())
	}
	if payload.Error.Code != cfg.ErrNotFound {
		t.Errorf("got code %q, want %q", payload.Error.Code, cfg.ErrNotFound)
	}
}

func TestGetSlide_InvalidPresID(t *testing.T) {
	store := &mockPresentationStore{
		getSlideFunc: func(ctx context.Context, presID, slideID string) (*Slide, error) {
			t.Error("store.GetSlide should not be called")
			return nil, nil
		},
	}

	handler := getSlideHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, "bad")
	req.SetPathValue(pathKeySlideID, validSlideID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestGetSlide_InvalidSlideID(t *testing.T) {
	store := &mockPresentationStore{
		getSlideFunc: func(ctx context.Context, presID, slideID string) (*Slide, error) {
			t.Error("store.GetSlide should not be called")
			return nil, nil
		},
	}

	handler := getSlideHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, "bad")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestGetSlide_StoreError(t *testing.T) {
	store := &mockPresentationStore{
		getSlideFunc: func(ctx context.Context, presID, slideID string) (*Slide, error) {
			return nil, errors.New("db error")
		},
	}

	handler := getSlideHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, validSlideID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestUpdateSlide_Valid(t *testing.T) {
	var called bool
	store := &mockPresentationStore{
		updateSlideFunc: func(ctx context.Context, presID, slideID, content string) error {
			called = true
			if strings.Contains(content, "<script") {
				t.Error("content should not contain script tags after sanitization")
			}
			if !strings.Contains(content, "<h1>") {
				t.Error("content should contain allowed tags after sanitization")
			}
			return nil
		},
	}

	handler := updateSlideHandler(store)
	body := `{"content":"<script>alert('xss')</script><h1>Updated</h1>"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, validSlideID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if !called {
		t.Error("expected store.UpdateSlide to be called")
	}

	var payload struct {
		Data Slide `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected valid JSON, got:", rec.Body.String())
	}
	if payload.Data.ID != validSlideID {
		t.Errorf("got id %q, want %q", payload.Data.ID, validSlideID)
	}
	if strings.Contains(payload.Data.Content, "<script") {
		t.Error("response content should not contain script tags after sanitization")
	}
	if !strings.Contains(payload.Data.Content, "<h1>") {
		t.Error("response content should contain allowed tags after sanitization")
	}
}

func TestUpdateSlide_NotFound(t *testing.T) {
	store := &mockPresentationStore{
		updateSlideFunc: func(ctx context.Context, presID, slideID, content string) error {
			return ErrNotFound
		},
	}

	handler := updateSlideHandler(store)
	body := `{"content":"<h1>Updated</h1>"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, validSlideID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}

	var payload response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected error JSON, got:", rec.Body.String())
	}
	if payload.Error.Code != cfg.ErrNotFound {
		t.Errorf("got code %q, want %q", payload.Error.Code, cfg.ErrNotFound)
	}
}

func TestUpdateSlide_InvalidSlideID(t *testing.T) {
	store := &mockPresentationStore{
		updateSlideFunc: func(ctx context.Context, presID, slideID, content string) error {
			t.Error("store.UpdateSlide should not be called")
			return nil
		},
	}

	handler := updateSlideHandler(store)
	body := `{"content":"<h1>Updated</h1>"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, "bad")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestUpdateSlide_InvalidPresID(t *testing.T) {
	store := &mockPresentationStore{
		updateSlideFunc: func(ctx context.Context, presID, slideID, content string) error {
			t.Error("store.UpdateSlide should not be called")
			return nil
		},
	}

	handler := updateSlideHandler(store)
	body := `{"content":"<h1>Updated</h1>"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, "bad")
	req.SetPathValue(pathKeySlideID, validSlideID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestUpdateSlide_EmptyContent(t *testing.T) {
	store := &mockPresentationStore{
		updateSlideFunc: func(ctx context.Context, presID, slideID, content string) error {
			t.Error("store.UpdateSlide should not be called")
			return nil
		},
	}

	handler := updateSlideHandler(store)
	body := `{"content":"  "}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, validSlideID)
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

func TestUpdateSlide_StoreError(t *testing.T) {
	store := &mockPresentationStore{
		updateSlideFunc: func(ctx context.Context, presID, slideID, content string) error {
			return errors.New("db error")
		},
	}

	handler := updateSlideHandler(store)
	body := `{"content":"<h1>Updated</h1>"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, validSlideID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestUpdateSlide_InvalidJSON(t *testing.T) {
	store := &mockPresentationStore{
		updateSlideFunc: func(ctx context.Context, presID, slideID, content string) error {
			t.Error("store.UpdateSlide should not be called")
			return nil
		},
	}

	handler := updateSlideHandler(store)
	body := `{bad`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, validSlideID)
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

func TestDeleteSlide_Valid(t *testing.T) {
	var called bool
	store := &mockPresentationStore{
		deleteSlideFunc: func(ctx context.Context, presID, slideID string) error {
			called = true
			return nil
		},
	}

	handler := deleteSlideHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, validSlideID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusNoContent)
	}
	if !called {
		t.Error("expected store.DeleteSlide to be called")
	}
}

func TestDeleteSlide_NotFound(t *testing.T) {
	store := &mockPresentationStore{
		deleteSlideFunc: func(ctx context.Context, presID, slideID string) error {
			return ErrNotFound
		},
	}

	handler := deleteSlideHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, validSlideID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}

	var payload response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected error JSON, got:", rec.Body.String())
	}
	if payload.Error.Code != cfg.ErrNotFound {
		t.Errorf("got code %q, want %q", payload.Error.Code, cfg.ErrNotFound)
	}
}

func TestDeleteSlide_InvalidPresID(t *testing.T) {
	store := &mockPresentationStore{
		deleteSlideFunc: func(ctx context.Context, presID, slideID string) error {
			t.Error("store.DeleteSlide should not be called")
			return nil
		},
	}

	handler := deleteSlideHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.SetPathValue(pathKeyPresentationID, "bad")
	req.SetPathValue(pathKeySlideID, validSlideID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestDeleteSlide_InvalidSlideID(t *testing.T) {
	store := &mockPresentationStore{
		deleteSlideFunc: func(ctx context.Context, presID, slideID string) error {
			t.Error("store.DeleteSlide should not be called")
			return nil
		},
	}

	handler := deleteSlideHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, "bad")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestDeleteSlide_StoreError(t *testing.T) {
	store := &mockPresentationStore{
		deleteSlideFunc: func(ctx context.Context, presID, slideID string) error {
			return errors.New("db error")
		},
	}

	handler := deleteSlideHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresID)
	req.SetPathValue(pathKeySlideID, validSlideID)
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
