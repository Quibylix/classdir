package hub

import (
	"context"
	"encoding/json"

	"classdir/api/internal/presentation"
	"classdir/api/internal/shared/cfg"
	"classdir/api/internal/shared/validate"
)

type InitParams struct {
	PresentationID string `json:"presentation_id"`
}

type JoinParams struct {
	PresentationID string `json:"presentation_id"`
}

type presentationStatus struct {
	PresentationID string               `json:"presentation_id"`
	Slides         []presentation.Slide `json:"slides"`
	CurrentIndex   int                  `json:"current_index"`
}

type InitHandler struct{}

func (h InitHandler) Name() string { return CmdInitPresentation }

func (h InitHandler) Handle(ctx CommandContext, params json.RawMessage) {
	if !ctx.Client.Authenticated {
		ctx.Client.writeError(cfg.ErrUnauthorized, cfg.ErrMsgInvalidToken)
		return
	}

	var p InitParams
	if err := json.Unmarshal(params, &p); err != nil {
		ctx.Client.writeError(cfg.ErrInvalidJSON, cfg.ErrMsgInvalidJSON)
		return
	}

	if !validate.IsValidUUIDv7(p.PresentationID) {
		ctx.Client.writeError(cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
		return
	}

	pres, err := ctx.Hub.Store().GetByID(context.Background(), p.PresentationID)
	if err != nil {
		ctx.Client.writeError(cfg.ErrInternalError, cfg.ErrMsgGetPresentation)
		return
	}
	if pres == nil {
		ctx.Client.writeError(cfg.ErrNotFound, cfg.ErrMsgNotFound)
		return
	}

	room := ctx.Hub.GetOrCreateRoom(p.PresentationID)
	room.slides = pres.Slides
	if room.currentIndex >= len(room.slides) {
		room.currentIndex = 0
	}

	ctx.Client.room = room
	room.controller = ctx.Client
	room.register <- ctx.Client

	data, _ := json.Marshal(presentationStatus{
		PresentationID: p.PresentationID,
		Slides:         pres.Slides,
		CurrentIndex:   room.currentIndex,
	})
	ctx.Client.writeData(data)
}

type JoinHandler struct{}

func (h JoinHandler) Name() string { return CmdJoinRoom }

func (h JoinHandler) Handle(ctx CommandContext, params json.RawMessage) {
	var p JoinParams
	if err := json.Unmarshal(params, &p); err != nil {
		ctx.Client.writeError(cfg.ErrInvalidJSON, cfg.ErrMsgInvalidJSON)
		return
	}

	if !validate.IsValidUUIDv7(p.PresentationID) {
		ctx.Client.writeError(cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
		return
	}

	room := ctx.Hub.GetRoom(p.PresentationID)
	if room == nil {
		ctx.Client.writeError(cfg.ErrNotFound, cfg.ErrMsgNotFound)
		return
	}

	ctx.Client.room = room
	room.register <- ctx.Client

	data, _ := json.Marshal(presentationStatus{
		PresentationID: p.PresentationID,
		Slides:         room.slides,
		CurrentIndex:   room.currentIndex,
	})
	ctx.Client.writeData(data)
}

func init() {
	register(InitHandler{}, false)
	register(JoinHandler{}, false)
}
