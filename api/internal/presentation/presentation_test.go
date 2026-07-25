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
	createFunc  func(ctx context.Context, id, title string) error
	getByIDFunc func(ctx context.Context, id string) (*Presentation, error)
	updateFunc  func(ctx context.Context, id, title, content string) error
	deleteFunc  func(ctx context.Context, id string) error
	listFunc    func(ctx context.Context) ([]*PresentationPreview, error)
}

func (m *mockPresentationStore) Create(ctx context.Context, id, title string) error {
	return m.createFunc(ctx, id, title)
}

func (m *mockPresentationStore) GetByID(ctx context.Context, id string) (*Presentation, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockPresentationStore) Update(ctx context.Context, id, title, content string) error {
	return m.updateFunc(ctx, id, title, content)
}

func (m *mockPresentationStore) Delete(ctx context.Context, id string) error {
	return m.deleteFunc(ctx, id)
}

func (m *mockPresentationStore) List(ctx context.Context) ([]*PresentationPreview, error) {
	return m.listFunc(ctx)
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

func TestGetPresentation_Found(t *testing.T) {
	store := &mockPresentationStore{
		getByIDFunc: func(ctx context.Context, id string) (*Presentation, error) {
			return &Presentation{
				ID:      id,
				Title:   "Test",
				Content: "<h1>Hi</h1>",
			}, nil
		},
	}

	handler := getPresentationHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data Presentation `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected valid JSON, got:", rec.Body.String())
	}
	if payload.Data.ID != "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f" {
		t.Errorf("got id %q, want %q", payload.Data.ID, "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f")
	}
	if payload.Data.Content != "<h1>Hi</h1>" {
		t.Errorf("got content %q, want %q", payload.Data.Content, "<h1>Hi</h1>")
	}
}

func TestGetPresentation_NotFound(t *testing.T) {
	store := &mockPresentationStore{
		getByIDFunc: func(ctx context.Context, id string) (*Presentation, error) {
			return nil, nil
		},
	}

	handler := getPresentationHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f")
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

func TestGetPresentation_InvalidUUID(t *testing.T) {
	store := &mockPresentationStore{
		getByIDFunc: func(ctx context.Context, id string) (*Presentation, error) {
			t.Error("store.GetByID should not be called")
			return nil, nil
		},
	}

	handler := getPresentationHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
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

func TestUpdatePresentation_Valid(t *testing.T) {
	var called bool
	store := &mockPresentationStore{
		updateFunc: func(ctx context.Context, id, title, content string) error {
			called = true
			return nil
		},
		getByIDFunc: func(ctx context.Context, id string) (*Presentation, error) {
			return &Presentation{ID: id, Title: "Updated", Content: "<h1>Hello</h1>"}, nil
		},
	}

	handler := updatePresentationHandler(store)
	body := `{"title":"Updated","content":"<h1>Hello</h1>"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if !called {
		t.Error("expected store.Update to be called")
	}
}

func TestUpdatePresentation_InvalidUUID(t *testing.T) {
	store := &mockPresentationStore{
		updateFunc: func(ctx context.Context, id, title, content string) error {
			t.Error("store.Update should not be called")
			return nil
		},
	}

	handler := updatePresentationHandler(store)
	body := `{"title":"Updated","content":""}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
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

func TestUpdatePresentation_EmptyTitle(t *testing.T) {
	store := &mockPresentationStore{
		updateFunc: func(ctx context.Context, id, title, content string) error {
			t.Error("store.Update should not be called")
			return nil
		},
	}

	handler := updatePresentationHandler(store)
	body := `{"title":"  ","content":""}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f")
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

func TestUpdatePresentation_InvalidJSON(t *testing.T) {
	store := &mockPresentationStore{
		updateFunc: func(ctx context.Context, id, title, content string) error {
			t.Error("store.Update should not be called")
			return nil
		},
	}

	handler := updatePresentationHandler(store)
	body := `{bad`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f")
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

func TestUpdatePresentation_NotFound(t *testing.T) {
	store := &mockPresentationStore{
		updateFunc: func(ctx context.Context, id, title, content string) error {
			return ErrNotFound
		},
	}

	handler := updatePresentationHandler(store)
	body := `{"title":"Updated","content":""}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f")
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

func TestUpdatePresentation_StoreError(t *testing.T) {
	store := &mockPresentationStore{
		updateFunc: func(ctx context.Context, id, title, content string) error {
			return errors.New("db error")
		},
	}

	handler := updatePresentationHandler(store)
	body := `{"title":"Updated","content":""}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f")
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

func TestUpdatePresentation_SanitizesScriptTag(t *testing.T) {
	var capturedContent string
	store := &mockPresentationStore{
		updateFunc: func(ctx context.Context, id, title, content string) error {
			capturedContent = content
			return nil
		},
		getByIDFunc: func(ctx context.Context, id string) (*Presentation, error) {
			return &Presentation{ID: id, Title: "Updated", Content: capturedContent}, nil
		},
	}

	handler := updatePresentationHandler(store)
	body := `{"title":"Updated","content":"<script>alert('xss')</script><h1>Hello</h1>"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	if strings.Contains(capturedContent, "<script") {
		t.Error("content should not contain script tags after sanitization")
	}
	if !strings.Contains(capturedContent, "<h1>") {
		t.Error("content should contain allowed tags after sanitization")
	}

	var payload struct {
		Data Presentation `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected valid JSON, got:", rec.Body.String())
	}
	if strings.Contains(payload.Data.Content, "<script") {
		t.Error("response content should not contain script tags after sanitization")
	}
}

func TestUpdatePresentation_SanitizesIframeTag(t *testing.T) {
	var capturedContent string
	store := &mockPresentationStore{
		updateFunc: func(ctx context.Context, id, title, content string) error {
			capturedContent = content
			return nil
		},
		getByIDFunc: func(ctx context.Context, id string) (*Presentation, error) {
			return &Presentation{ID: id, Title: "Updated", Content: capturedContent}, nil
		},
	}

	handler := updatePresentationHandler(store)
	body := `{"title":"Updated","content":"<iframe src=\"evil.com\"></iframe><h1>Safe</h1>"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	if strings.Contains(capturedContent, "<iframe") {
		t.Error("content should not contain iframe tags after sanitization")
	}
	if !strings.Contains(capturedContent, "<h1>") {
		t.Error("content should contain allowed tags after sanitization")
	}
}

func TestDeletePresentation_Valid(t *testing.T) {
	var called bool
	store := &mockPresentationStore{
		deleteFunc: func(ctx context.Context, id string) error {
			called = true
			return nil
		},
	}

	handler := deletePresentationHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.SetPathValue(pathKeyPresentationID, "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusNoContent)
	}
	if !called {
		t.Error("expected store.Delete to be called")
	}
}

func TestDeletePresentation_InvalidUUID(t *testing.T) {
	store := &mockPresentationStore{
		deleteFunc: func(ctx context.Context, id string) error {
			t.Error("store.Delete should not be called")
			return nil
		},
	}

	handler := deletePresentationHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
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

func TestDeletePresentation_NotFound(t *testing.T) {
	store := &mockPresentationStore{
		deleteFunc: func(ctx context.Context, id string) error {
			return ErrNotFound
		},
	}

	handler := deletePresentationHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.SetPathValue(pathKeyPresentationID, "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f")
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

func TestDeletePresentation_StoreError(t *testing.T) {
	store := &mockPresentationStore{
		deleteFunc: func(ctx context.Context, id string) error {
			return errors.New("db error")
		},
	}

	handler := deletePresentationHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.SetPathValue(pathKeyPresentationID, "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f")
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

func TestListPresentations_Success(t *testing.T) {
	store := &mockPresentationStore{
		listFunc: func(ctx context.Context) ([]*PresentationPreview, error) {
			return []*PresentationPreview{
				{ID: "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f", Title: "First"},
				{ID: "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80", Title: "Second"},
			}, nil
		},
	}

	handler := listPresentationHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data []*PresentationPreview `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected valid JSON, got:", rec.Body.String())
	}
	if len(payload.Data) != 2 {
		t.Errorf("got %d presentations, want 2", len(payload.Data))
	}
	if payload.Data[0].Title != "First" {
		t.Errorf("got title %q, want %q", payload.Data[0].Title, "First")
	}
}

func TestListPresentations_StoreError(t *testing.T) {
	store := &mockPresentationStore{
		listFunc: func(ctx context.Context) ([]*PresentationPreview, error) {
			return nil, errors.New("db error")
		},
	}

	handler := listPresentationHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
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
