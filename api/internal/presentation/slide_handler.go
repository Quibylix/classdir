package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"classdir/api/internal/shared/cfg"
	"classdir/api/internal/shared/response"
	"classdir/api/internal/shared/sanitize"
	"classdir/api/internal/shared/validate"
)

func createSlideHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presID := r.PathValue(pathKeyPresentationID)
		if !validate.IsValidUUIDv7(presID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		var body struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidJSON, cfg.ErrMsgInvalidJSON)
			return
		}

		if !validate.IsValidUUIDv7(body.ID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		if strings.TrimSpace(body.Content) == "" {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrMissingField, cfg.ErrMsgMissingContent)
			return
		}

		sanitized := sanitize.RevealPolicy.Sanitize(body.Content)

		if err := store.CreateSlide(r.Context(), presID, body.ID, sanitized); err != nil {
			if errors.Is(err, ErrNotFound) {
				response.WriteError(w, http.StatusNotFound, cfg.ErrNotFound, cfg.ErrMsgNotFound)
				return
			}
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgCreateSlide)
			return
		}

		data, err := json.Marshal(Slide{
			ID:      body.ID,
			Content: sanitized,
		})
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgCreateSlide)
			return
		}
		response.WriteJSON(w, http.StatusCreated, data)
	}
}

func getSlideHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presID := r.PathValue(pathKeyPresentationID)
		if !validate.IsValidUUIDv7(presID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		slideID := r.PathValue(pathKeySlideID)
		if !validate.IsValidUUIDv7(slideID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		slide, err := store.GetSlide(r.Context(), presID, slideID)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgGetSlide)
			return
		}
		if slide == nil {
			response.WriteError(w, http.StatusNotFound, cfg.ErrNotFound, cfg.ErrMsgNotFound)
			return
		}

		data, err := json.Marshal(slide)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgGetSlide)
			return
		}
		response.WriteJSON(w, http.StatusOK, data)
	}
}

func updateSlideHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presID := r.PathValue(pathKeyPresentationID)
		if !validate.IsValidUUIDv7(presID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		slideID := r.PathValue(pathKeySlideID)
		if !validate.IsValidUUIDv7(slideID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		var body struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidJSON, cfg.ErrMsgInvalidJSON)
			return
		}

		if strings.TrimSpace(body.Content) == "" {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrMissingField, cfg.ErrMsgMissingContent)
			return
		}

		sanitized := sanitize.RevealPolicy.Sanitize(body.Content)

		if err := store.UpdateSlide(r.Context(), presID, slideID, sanitized); err != nil {
			if errors.Is(err, ErrNotFound) {
				response.WriteError(w, http.StatusNotFound, cfg.ErrNotFound, cfg.ErrMsgNotFound)
				return
			}
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgUpdateSlide)
			return
		}

		data, err := json.Marshal(Slide{
			ID:      slideID,
			Content: sanitized,
		})
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgUpdateSlide)
			return
		}
		response.WriteJSON(w, http.StatusOK, data)
	}
}

func deleteSlideHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presID := r.PathValue(pathKeyPresentationID)
		if !validate.IsValidUUIDv7(presID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		slideID := r.PathValue(pathKeySlideID)
		if !validate.IsValidUUIDv7(slideID) {
			response.WriteError(w, http.StatusBadRequest, cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
			return
		}

		if err := store.DeleteSlide(r.Context(), presID, slideID); err != nil {
			if errors.Is(err, ErrNotFound) {
				response.WriteError(w, http.StatusNotFound, cfg.ErrNotFound, cfg.ErrMsgNotFound)
				return
			}
			response.WriteError(w, http.StatusInternalServerError, cfg.ErrInternalError, cfg.ErrMsgDeleteSlide)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
