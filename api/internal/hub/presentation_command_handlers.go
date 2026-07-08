package hub

import (
	"encoding/json"

	"classdir/api/internal/presentation"
	"classdir/api/internal/shared/cfg"
	"classdir/api/internal/shared/validate"
)

type InitParams struct {
	PresentationID string `json:"presentation_id"`
}

type JoinParams struct {
	RoomCode string `json:"room_code"`
}

type presentationStatus struct {
	PresentationID string               `json:"presentation_id"`
	Slides         []presentation.Slide `json:"slides"`
	CurrentIndex   int                  `json:"current_index"`
	RoomCode       string               `json:"room_code,omitempty"`
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

	room, err := ctx.Hub.GetOrCreateRoom(p.PresentationID)

	if err == presentation.ErrNotFound {
		ctx.Client.writeError(cfg.ErrNotFound, cfg.ErrMsgNotFound)
		return
	}

	if err != nil {
		ctx.Client.writeError(cfg.ErrInternalError, cfg.ErrMsgCreateRoom)
		return
	}

	ctx.Client.room = room
	select {
	case room.register <- &registrationRequest{client: ctx.Client, isController: true}:
	case <-room.done:
		ctx.Client.writeError(cfg.ErrInternalError, cfg.ErrMsgRoomClosed)
		return
	}

}

type JoinHandler struct{}

func (h JoinHandler) Name() string { return CmdJoinRoom }

func (h JoinHandler) Handle(ctx CommandContext, params json.RawMessage) {
	var p JoinParams
	if err := json.Unmarshal(params, &p); err != nil {
		ctx.Client.writeError(cfg.ErrInvalidJSON, cfg.ErrMsgInvalidJSON)
		return
	}

	room := ctx.Hub.GetRoomByCode(p.RoomCode)
	if room == nil {
		ctx.Client.writeError(cfg.ErrNotFound, cfg.ErrMsgNotFound)
		return
	}

	ctx.Client.room = room
	select {
	case room.register <- &registrationRequest{client: ctx.Client, isController: false}:
	case <-room.done:
		ctx.Client.writeError(cfg.ErrInternalError, cfg.ErrMsgRoomClosed)
		return
	}
}

func init() {
	register(InitHandler{}, false)
	register(JoinHandler{}, false)
}
