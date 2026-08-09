package student

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

const (
	validPresentationUUID = "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f"
	validStudentUUID      = "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80"
)

type mockStudentStore struct {
	createFunc  func(ctx context.Context, presentationID, id, name string) error
	getByIDFunc func(ctx context.Context, presentationID, id string) (*Student, error)
	updateFunc  func(ctx context.Context, presentationID, id, name string) error
	deleteFunc  func(ctx context.Context, presentationID, id string) error
	listFunc    func(ctx context.Context, presentationID string) ([]*Student, error)
}

func (m *mockStudentStore) Create(ctx context.Context, presentationID, id, name string) error {
	return m.createFunc(ctx, presentationID, id, name)
}

func (m *mockStudentStore) GetByID(ctx context.Context, presentationID, id string) (*Student, error) {
	return m.getByIDFunc(ctx, presentationID, id)
}

func (m *mockStudentStore) Update(ctx context.Context, presentationID, id, name string) error {
	return m.updateFunc(ctx, presentationID, id, name)
}

func (m *mockStudentStore) Delete(ctx context.Context, presentationID, id string) error {
	return m.deleteFunc(ctx, presentationID, id)
}

func (m *mockStudentStore) List(ctx context.Context, presentationID string) ([]*Student, error) {
	return m.listFunc(ctx, presentationID)
}

func TestCreateStudent_ValidInput(t *testing.T) {
	var called bool
	store := &mockStudentStore{
		createFunc: func(ctx context.Context, presentationID, id, name string) error {
			called = true
			if presentationID != validPresentationUUID {
				t.Errorf("got presentationID %q, want %q", presentationID, validPresentationUUID)
			}
			if id != validStudentUUID {
				t.Errorf("got id %q, want %q", id, validStudentUUID)
			}
			if name != "Alice" {
				t.Errorf("got name %q, want %q", name, "Alice")
			}
			return nil
		},
	}

	handler := createStudentHandler(store)
	body := `{"id":"0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80","name":"Alice"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusCreated)
	}
	if !called {
		t.Error("expected store.Create to be called")
	}

	var payload struct {
		Data Student `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected valid JSON, got:", rec.Body.String())
	}
	if payload.Data.ID != validStudentUUID {
		t.Errorf("got id %q, want %q", payload.Data.ID, validStudentUUID)
	}
	if payload.Data.Name != "Alice" {
		t.Errorf("got name %q, want %q", payload.Data.Name, "Alice")
	}
}

func TestCreateStudent_InvalidJSON(t *testing.T) {
	store := &mockStudentStore{
		createFunc: func(ctx context.Context, presentationID, id, name string) error {
			t.Error("store.Create should not be called")
			return nil
		},
	}

	handler := createStudentHandler(store)
	body := `{bad`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
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

func TestCreateStudent_InvalidPresentationUUID(t *testing.T) {
	store := &mockStudentStore{
		createFunc: func(ctx context.Context, presentationID, id, name string) error {
			t.Error("store.Create should not be called")
			return nil
		},
	}

	handler := createStudentHandler(store)
	body := `{"id":"0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80","name":"Alice"}`
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

func TestCreateStudent_InvalidStudentUUID(t *testing.T) {
	store := &mockStudentStore{
		createFunc: func(ctx context.Context, presentationID, id, name string) error {
			t.Error("store.Create should not be called")
			return nil
		},
	}

	handler := createStudentHandler(store)
	body := `{"id":"bad","name":"Alice"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
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

func TestCreateStudent_EmptyName(t *testing.T) {
	store := &mockStudentStore{
		createFunc: func(ctx context.Context, presentationID, id, name string) error {
			t.Error("store.Create should not be called")
			return nil
		},
	}

	handler := createStudentHandler(store)
	body := `{"id":"0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80","name":"  "}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
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

func TestCreateStudent_DuplicateID(t *testing.T) {
	store := &mockStudentStore{
		createFunc: func(ctx context.Context, presentationID, id, name string) error {
			return ErrDuplicateID
		},
	}

	handler := createStudentHandler(store)
	body := `{"id":"0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80","name":"Alice"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusConflict)
	}

	var payload response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected error JSON, got:", rec.Body.String())
	}
	if payload.Error.Code != cfg.ErrConflict {
		t.Errorf("got code %q, want %q", payload.Error.Code, cfg.ErrConflict)
	}
	if payload.Error.Message != cfg.ErrMsgDuplicateStudentID {
		t.Errorf("got message %q, want %q", payload.Error.Message, cfg.ErrMsgDuplicateStudentID)
	}
}

func TestCreateStudent_DuplicateName(t *testing.T) {
	store := &mockStudentStore{
		createFunc: func(ctx context.Context, presentationID, id, name string) error {
			return ErrDuplicateName
		},
	}

	handler := createStudentHandler(store)
	body := `{"id":"0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80","name":"Alice"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusConflict)
	}

	var payload response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected error JSON, got:", rec.Body.String())
	}
	if payload.Error.Code != cfg.ErrConflict {
		t.Errorf("got code %q, want %q", payload.Error.Code, cfg.ErrConflict)
	}
	if payload.Error.Message != cfg.ErrMsgDuplicateStudentName {
		t.Errorf("got message %q, want %q", payload.Error.Message, cfg.ErrMsgDuplicateStudentName)
	}
}

func TestCreateStudent_PresentationNotFound(t *testing.T) {
	store := &mockStudentStore{
		createFunc: func(ctx context.Context, presentationID, id, name string) error {
			return ErrNotFound
		},
	}

	handler := createStudentHandler(store)
	body := `{"id":"0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80","name":"Alice"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
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
	if payload.Error.Message != cfg.ErrMsgNotFound {
		t.Errorf("got message %q, want %q", payload.Error.Message, cfg.ErrMsgNotFound)
	}
}

func TestCreateStudent_StoreError(t *testing.T) {
	store := &mockStudentStore{
		createFunc: func(ctx context.Context, presentationID, id, name string) error {
			return errors.New("db error")
		},
	}

	handler := createStudentHandler(store)
	body := `{"id":"0192e5a0-7b7f-7b7f-8b7f-0192e5a07b80","name":"Alice"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
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

func TestListStudents_Success(t *testing.T) {
	store := &mockStudentStore{
		listFunc: func(ctx context.Context, presentationID string) ([]*Student, error) {
			if presentationID != validPresentationUUID {
				t.Errorf("got presentationID %q, want %q", presentationID, validPresentationUUID)
			}
			return []*Student{
				{ID: validStudentUUID, Name: "Alice"},
				{ID: "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b81", Name: "Bob"},
			}, nil
		},
	}

	handler := listStudentsHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data []*Student `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected valid JSON, got:", rec.Body.String())
	}
	if len(payload.Data) != 2 {
		t.Fatalf("got %d students, want 2", len(payload.Data))
	}
	if payload.Data[0].Name != "Alice" {
		t.Errorf("got name %q, want %q", payload.Data[0].Name, "Alice")
	}
	if payload.Data[1].Name != "Bob" {
		t.Errorf("got name %q, want %q", payload.Data[1].Name, "Bob")
	}
}

func TestListStudents_Empty(t *testing.T) {
	store := &mockStudentStore{
		listFunc: func(ctx context.Context, presentationID string) ([]*Student, error) {
			return []*Student{}, nil
		},
	}

	handler := listStudentsHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data []*Student `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected valid JSON, got:", rec.Body.String())
	}
	if payload.Data == nil {
		t.Error("expected data to be an empty array, got null")
	}
	if len(payload.Data) != 0 {
		t.Errorf("got %d students, want 0", len(payload.Data))
	}
}

func TestListStudents_InvalidPresentationUUID(t *testing.T) {
	store := &mockStudentStore{
		listFunc: func(ctx context.Context, presentationID string) ([]*Student, error) {
			t.Error("store.List should not be called")
			return nil, nil
		},
	}

	handler := listStudentsHandler(store)
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

func TestListStudents_StoreError(t *testing.T) {
	store := &mockStudentStore{
		listFunc: func(ctx context.Context, presentationID string) ([]*Student, error) {
			return nil, errors.New("db error")
		},
	}

	handler := listStudentsHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
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

func TestGetStudent_Found(t *testing.T) {
	store := &mockStudentStore{
		getByIDFunc: func(ctx context.Context, presentationID, id string) (*Student, error) {
			if presentationID != validPresentationUUID {
				t.Errorf("got presentationID %q, want %q", presentationID, validPresentationUUID)
			}
			return &Student{ID: id, Name: "Alice"}, nil
		},
	}

	handler := getStudentHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var payload struct {
		Data Student `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected valid JSON, got:", rec.Body.String())
	}
	if payload.Data.ID != validStudentUUID {
		t.Errorf("got id %q, want %q", payload.Data.ID, validStudentUUID)
	}
	if payload.Data.Name != "Alice" {
		t.Errorf("got name %q, want %q", payload.Data.Name, "Alice")
	}
}

func TestGetStudent_NotFound(t *testing.T) {
	store := &mockStudentStore{
		getByIDFunc: func(ctx context.Context, presentationID, id string) (*Student, error) {
			return nil, nil
		},
	}

	handler := getStudentHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
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
	if payload.Error.Message != cfg.ErrMsgStudentNotFound {
		t.Errorf("got message %q, want %q", payload.Error.Message, cfg.ErrMsgStudentNotFound)
	}
}

func TestGetStudent_InvalidPresentationUUID(t *testing.T) {
	store := &mockStudentStore{
		getByIDFunc: func(ctx context.Context, presentationID, id string) (*Student, error) {
			t.Error("store.GetByID should not be called")
			return nil, nil
		},
	}

	handler := getStudentHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, "bad")
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
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

func TestGetStudent_InvalidStudentUUID(t *testing.T) {
	store := &mockStudentStore{
		getByIDFunc: func(ctx context.Context, presentationID, id string) (*Student, error) {
			t.Error("store.GetByID should not be called")
			return nil, nil
		},
	}

	handler := getStudentHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, "bad")
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

func TestGetStudent_StoreError(t *testing.T) {
	store := &mockStudentStore{
		getByIDFunc: func(ctx context.Context, presentationID, id string) (*Student, error) {
			return nil, errors.New("db error")
		},
	}

	handler := getStudentHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
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

func TestUpdateStudent_Valid(t *testing.T) {
	var updateCalled bool
	store := &mockStudentStore{
		updateFunc: func(ctx context.Context, presentationID, id, name string) error {
			updateCalled = true
			if presentationID != validPresentationUUID {
				t.Errorf("got presentationID %q, want %q", presentationID, validPresentationUUID)
			}
			if id != validStudentUUID {
				t.Errorf("got id %q, want %q", id, validStudentUUID)
			}
			if name != "Updated" {
				t.Errorf("got name %q, want %q", name, "Updated")
			}
			return nil
		},
		getByIDFunc: func(ctx context.Context, presentationID, id string) (*Student, error) {
			return &Student{ID: id, Name: "Updated"}, nil
		},
	}

	handler := updateStudentHandler(store)
	body := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if !updateCalled {
		t.Error("expected store.Update to be called")
	}

	var payload struct {
		Data Student `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected valid JSON, got:", rec.Body.String())
	}
	if payload.Data.Name != "Updated" {
		t.Errorf("got name %q, want %q", payload.Data.Name, "Updated")
	}
}

func TestUpdateStudent_InvalidJSON(t *testing.T) {
	store := &mockStudentStore{
		updateFunc: func(ctx context.Context, presentationID, id, name string) error {
			t.Error("store.Update should not be called")
			return nil
		},
	}

	handler := updateStudentHandler(store)
	body := `{bad`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
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

func TestUpdateStudent_InvalidPresentationUUID(t *testing.T) {
	store := &mockStudentStore{
		updateFunc: func(ctx context.Context, presentationID, id, name string) error {
			t.Error("store.Update should not be called")
			return nil
		},
	}

	handler := updateStudentHandler(store)
	body := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, "bad")
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
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

func TestUpdateStudent_InvalidStudentUUID(t *testing.T) {
	store := &mockStudentStore{
		updateFunc: func(ctx context.Context, presentationID, id, name string) error {
			t.Error("store.Update should not be called")
			return nil
		},
	}

	handler := updateStudentHandler(store)
	body := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, "bad")
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

func TestUpdateStudent_EmptyName(t *testing.T) {
	store := &mockStudentStore{
		updateFunc: func(ctx context.Context, presentationID, id, name string) error {
			t.Error("store.Update should not be called")
			return nil
		},
	}

	handler := updateStudentHandler(store)
	body := `{"name":"  "}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
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

func TestUpdateStudent_NotFound(t *testing.T) {
	store := &mockStudentStore{
		updateFunc: func(ctx context.Context, presentationID, id, name string) error {
			return ErrNotFound
		},
	}

	handler := updateStudentHandler(store)
	body := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
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
	if payload.Error.Message != cfg.ErrMsgStudentNotFound {
		t.Errorf("got message %q, want %q", payload.Error.Message, cfg.ErrMsgStudentNotFound)
	}
}

func TestUpdateStudent_DuplicateName(t *testing.T) {
	store := &mockStudentStore{
		updateFunc: func(ctx context.Context, presentationID, id, name string) error {
			return ErrDuplicateName
		},
	}

	handler := updateStudentHandler(store)
	body := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusConflict)
	}

	var payload response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatal("expected error JSON, got:", rec.Body.String())
	}
	if payload.Error.Code != cfg.ErrConflict {
		t.Errorf("got code %q, want %q", payload.Error.Code, cfg.ErrConflict)
	}
	if payload.Error.Message != cfg.ErrMsgDuplicateStudentName {
		t.Errorf("got message %q, want %q", payload.Error.Message, cfg.ErrMsgDuplicateStudentName)
	}
}

func TestUpdateStudent_StoreError(t *testing.T) {
	store := &mockStudentStore{
		updateFunc: func(ctx context.Context, presentationID, id, name string) error {
			return errors.New("db error")
		},
	}

	handler := updateStudentHandler(store)
	body := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
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

func TestDeleteStudent_Valid(t *testing.T) {
	var called bool
	store := &mockStudentStore{
		deleteFunc: func(ctx context.Context, presentationID, id string) error {
			called = true
			if presentationID != validPresentationUUID {
				t.Errorf("got presentationID %q, want %q", presentationID, validPresentationUUID)
			}
			if id != validStudentUUID {
				t.Errorf("got id %q, want %q", id, validStudentUUID)
			}
			return nil
		},
	}

	handler := deleteStudentHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusNoContent)
	}
	if !called {
		t.Error("expected store.Delete to be called")
	}
}

func TestDeleteStudent_InvalidPresentationUUID(t *testing.T) {
	store := &mockStudentStore{
		deleteFunc: func(ctx context.Context, presentationID, id string) error {
			t.Error("store.Delete should not be called")
			return nil
		},
	}

	handler := deleteStudentHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.SetPathValue(pathKeyPresentationID, "bad")
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
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

func TestDeleteStudent_InvalidStudentUUID(t *testing.T) {
	store := &mockStudentStore{
		deleteFunc: func(ctx context.Context, presentationID, id string) error {
			t.Error("store.Delete should not be called")
			return nil
		},
	}

	handler := deleteStudentHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, "bad")
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

func TestDeleteStudent_NotFound(t *testing.T) {
	store := &mockStudentStore{
		deleteFunc: func(ctx context.Context, presentationID, id string) error {
			return ErrNotFound
		},
	}

	handler := deleteStudentHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
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
	if payload.Error.Message != cfg.ErrMsgStudentNotFound {
		t.Errorf("got message %q, want %q", payload.Error.Message, cfg.ErrMsgStudentNotFound)
	}
}

func TestDeleteStudent_StoreError(t *testing.T) {
	store := &mockStudentStore{
		deleteFunc: func(ctx context.Context, presentationID, id string) error {
			return errors.New("db error")
		},
	}

	handler := deleteStudentHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.SetPathValue(pathKeyPresentationID, validPresentationUUID)
	req.SetPathValue(pathKeyStudentID, validStudentUUID)
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
