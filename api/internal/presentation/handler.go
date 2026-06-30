package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"classdir/api/internal/shared/cfg"
	"classdir/api/internal/shared/response"
	"classdir/api/internal/shared/validate"
)

type SlideMetadata struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type Slide struct {
	ID          string        `json:"id"`
	SlideNumber int           `json:"slide_number"`
	Content     string        `json:"content"`
	Metadata    SlideMetadata `json:"metadata"`
}

type Presentation struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Slides []Slide `json:"slides"`
}

func RegisterRoutes(mux *http.ServeMux, store Store) {
	mux.HandleFunc("POST /api/v1/presentation", createPresentationHandler(store))
	mux.HandleFunc("GET /api/v1/presentation/{presentationId}", getPresentationHandler(store))
}

func createPresentationHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidJSON, cfg.ErrMsgInvalidJSON)
			return
		}

		if !validate.IsValidUUIDv7(body.ID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		if strings.TrimSpace(body.Title) == "" {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrMissingField, cfg.ErrMsgMissingTitle)
			return
		}

		if err := store.Create(r.Context(), body.ID, body.Title); err != nil {
			if errors.Is(err, ErrDuplicateKey) {
				response.WriteError(w, http.StatusConflict, cfg.ErrConflict, cfg.ErrMsgDuplicateID)
				return
			}
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgCreatePresentation)
			return
		}

		data, err := json.Marshal(Presentation{
			ID:     body.ID,
			Title:  body.Title,
			Slides: []Slide{},
		})
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgCreatePresentation)
			return
		}
		response.WriteJSON(w, http.StatusCreated, data)
	}
}

func getPresentationHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("presentationId")

		if !validate.IsValidUUIDv7(id) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		pres, err := store.GetByID(r.Context(), id)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgGetPresentation)
			return
		}
		if pres == nil {
			response.WriteError(w, http.StatusNotFound, cfg.ErrNotFound, cfg.ErrMsgNotFound)
			return
		}

		data, err := json.Marshal(pres)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgGetPresentation)
			return
		}
		response.WriteJSON(w, http.StatusOK, data)
	}
}
