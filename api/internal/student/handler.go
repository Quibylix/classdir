package student

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"classdir/api/internal/shared/cfg"
	"classdir/api/internal/shared/response"
	"classdir/api/internal/shared/validate"
)

type Student struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func RegisterRoutes(mux *http.ServeMux, store Store) {
	mux.HandleFunc("POST /api/v1/presentation/{"+pathKeyPresentationID+"}/students", createStudentHandler(store))
	mux.HandleFunc("GET /api/v1/presentation/{"+pathKeyPresentationID+"}/students", listStudentsHandler(store))
	mux.HandleFunc("GET /api/v1/presentation/{"+pathKeyPresentationID+"}/students/{"+pathKeyStudentID+"}", getStudentHandler(store))
	mux.HandleFunc("PUT /api/v1/presentation/{"+pathKeyPresentationID+"}/students/{"+pathKeyStudentID+"}", updateStudentHandler(store))
	mux.HandleFunc("DELETE /api/v1/presentation/{"+pathKeyPresentationID+"}/students/{"+pathKeyStudentID+"}", deleteStudentHandler(store))
}

func createStudentHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidJSON, cfg.ErrMsgInvalidJSON)
			return
		}

		presentationID := r.PathValue(pathKeyPresentationID)
		if !validate.IsValidUUIDv7(presentationID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		if !validate.IsValidUUIDv7(body.ID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		if strings.TrimSpace(body.Name) == "" {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrMissingField, cfg.ErrMsgMissingName)
			return
		}

		if err := store.Create(r.Context(), presentationID, body.ID, body.Name); err != nil {
			if errors.Is(err, ErrDuplicateID) {
				response.WriteError(w, http.StatusConflict, cfg.ErrConflict, cfg.ErrMsgDuplicateStudentID)
				return
			}
			if errors.Is(err, ErrDuplicateName) {
				response.WriteError(w, http.StatusConflict, cfg.ErrConflict, cfg.ErrMsgDuplicateStudentName)
				return
			}
			if errors.Is(err, ErrNotFound) {
				response.WriteError(w, http.StatusNotFound, cfg.ErrNotFound, cfg.ErrMsgNotFound)
				return
			}
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgCreateStudent)
			return
		}

		data, err := json.Marshal(Student{
			ID:   body.ID,
			Name: body.Name,
		})
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgCreateStudent)
			return
		}
		response.WriteJSON(w, http.StatusCreated, data)
	}
}

func listStudentsHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presentationID := r.PathValue(pathKeyPresentationID)

		if !validate.IsValidUUIDv7(presentationID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		students, err := store.List(r.Context(), presentationID)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgListStudents)
			return
		}

		data, err := json.Marshal(students)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgListStudents)
			return
		}
		response.WriteJSON(w, http.StatusOK, data)
	}
}

func getStudentHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presentationID := r.PathValue(pathKeyPresentationID)

		if !validate.IsValidUUIDv7(presentationID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		studentID := r.PathValue(pathKeyStudentID)
		if !validate.IsValidUUIDv7(studentID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		st, err := store.GetByID(r.Context(), presentationID, studentID)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgGetStudent)
			return
		}
		if st == nil {
			response.WriteError(w, http.StatusNotFound, cfg.ErrNotFound, cfg.ErrMsgStudentNotFound)
			return
		}

		data, err := json.Marshal(st)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgGetStudent)
			return
		}
		response.WriteJSON(w, http.StatusOK, data)
	}
}

func updateStudentHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presentationID := r.PathValue(pathKeyPresentationID)

		if !validate.IsValidUUIDv7(presentationID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		studentID := r.PathValue(pathKeyStudentID)
		if !validate.IsValidUUIDv7(studentID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		var body struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidJSON, cfg.ErrMsgInvalidJSON)
			return
		}

		if strings.TrimSpace(body.Name) == "" {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrMissingField, cfg.ErrMsgMissingName)
			return
		}

		if err := store.Update(r.Context(), presentationID, studentID, body.Name); err != nil {
			if errors.Is(err, ErrNotFound) {
				response.WriteError(w, http.StatusNotFound, cfg.ErrNotFound, cfg.ErrMsgStudentNotFound)
				return
			}
			if errors.Is(err, ErrDuplicateName) {
				response.WriteError(w, http.StatusConflict, cfg.ErrConflict, cfg.ErrMsgDuplicateStudentName)
				return
			}
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgUpdateStudent)
			return
		}

		st, err := store.GetByID(r.Context(), presentationID, studentID)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgUpdateStudent)
			return
		}

		data, err := json.Marshal(st)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgUpdateStudent)
			return
		}
		response.WriteJSON(w, http.StatusOK, data)
	}
}

func deleteStudentHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presentationID := r.PathValue(pathKeyPresentationID)

		if !validate.IsValidUUIDv7(presentationID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		studentID := r.PathValue(pathKeyStudentID)
		if !validate.IsValidUUIDv7(studentID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		if err := store.Delete(r.Context(), presentationID, studentID); err != nil {
			if errors.Is(err, ErrNotFound) {
				response.WriteError(w, http.StatusNotFound, cfg.ErrNotFound, cfg.ErrMsgStudentNotFound)
				return
			}
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgDeleteStudent)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
