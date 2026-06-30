package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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

func createPresentationHandler(store presentationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, errInvalidJSON, errMsgInvalidJSON)
			return
		}

		if !isValidUUIDv7(body.ID) {
			writeError(w, http.StatusBadRequest, errInvalidUUID, errMsgInvalidID)
			return
		}

		if strings.TrimSpace(body.Title) == "" {
			writeError(w, http.StatusBadRequest, errMissingField, errMsgMissingTitle)
			return
		}

		if err := store.create(r.Context(), body.ID, body.Title); err != nil {
			if errors.Is(err, ErrDuplicateKey) {
				writeError(w, http.StatusConflict, errConflict, errMsgDuplicateID)
				return
			}
			writeError(w, http.StatusInternalServerError, errInternalError, errMsgCreatePresentation)
			return
		}

		data, err := json.Marshal(Presentation{
			ID:     body.ID,
			Title:  body.Title,
			Slides: []Slide{},
		})
		if err != nil {
			writeError(w, http.StatusInternalServerError, errInternalError, errMsgCreatePresentation)
			return
		}
		writeJSON(w, http.StatusCreated, data)
	}
}
