package hub

import (
	"context"
	"encoding/json"

	"github.com/coder/websocket"

	"classdir/api/internal/presentation"
	"classdir/api/internal/shared/cfg"
	"classdir/api/internal/shared/response"
	"classdir/api/internal/shared/validate"
)

const maxMessageSize = 65536

type Command struct {
	Command    string          `json:"command"`
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

type InitParams struct {
	PresentationID string `json:"presentation_id"`
}

type JoinParams struct {
	PresentationID string `json:"presentation_id"`
}

type SlideParams struct {
	SlideNumber int `json:"slide_number"`
}

type presentationStatus struct {
	PresentationID string               `json:"presentation_id"`
	Slides         []presentation.Slide `json:"slides"`
	CurrentIndex   int                  `json:"current_index"`
}

type Client struct {
	hub  *Hub
	conn wsConn
	send chan []byte
	room *Room
}

func NewClient(hub *Hub, conn wsConn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, channelBuffer),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		if c.room != nil {
			c.room.unregister <- c
		}
		c.conn.Close(websocket.StatusNormalClosure, "connection closed")
	}()

	for {
		_, msg, err := c.conn.Read(context.Background())
		if err != nil {
			break
		}

		var cmd Command
		if err := json.Unmarshal(msg, &cmd); err != nil {
			continue
		}

		c.handleCommand(cmd, msg)
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close(websocket.StatusNormalClosure, "connection closed")

	for msg := range c.send {
		err := c.conn.Write(context.Background(), websocket.MessageText, msg)
		if err != nil {
			break
		}
	}
}

func (c *Client) handleCommand(cmd Command, raw []byte) {
	switch cmd.Command {
	case CmdInitPresentation:
		c.handleInit(cmd.Parameters)
	case CmdJoinRoom:
		c.handleJoin(cmd.Parameters)
	default:
		if c.room != nil {
			c.room.commands <- roomCommand{msg: raw, sender: c}
		}
	}
}

func (c *Client) handleInit(paramsJSON json.RawMessage) {
	var params InitParams
	if err := json.Unmarshal(paramsJSON, &params); err != nil {
		c.writeError(cfg.ErrInvalidJSON, cfg.ErrMsgInvalidJSON)
		return
	}

	if !validate.IsValidUUIDv7(params.PresentationID) {
		c.writeError(cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
		return
	}

	pres, err := c.hub.Store().GetByID(context.Background(), params.PresentationID)
	if err != nil {
		c.writeError(cfg.ErrInternalError, cfg.ErrMsgGetPresentation)
		return
	}
	if pres == nil {
		c.writeError(cfg.ErrNotFound, cfg.ErrMsgNotFound)
		return
	}

	room := c.hub.GetOrCreateRoom(params.PresentationID)
	room.slides = pres.Slides
	room.currentIndex = 0

	c.room = room
	room.controller = c
	room.register <- c

	data, _ := json.Marshal(presentationStatus{
		PresentationID: params.PresentationID,
		Slides:         pres.Slides,
		CurrentIndex:   0,
	})
	c.writeData(data)
}

func (c *Client) handleJoin(paramsJSON json.RawMessage) {
	var params JoinParams
	if err := json.Unmarshal(paramsJSON, &params); err != nil {
		c.writeError(cfg.ErrInvalidJSON, cfg.ErrMsgInvalidJSON)
		return
	}

	if !validate.IsValidUUIDv7(params.PresentationID) {
		c.writeError(cfg.ErrInvalidUUID, cfg.ErrMsgInvalidID)
		return
	}

	room := c.hub.GetRoom(params.PresentationID)
	if room == nil {
		c.writeError(cfg.ErrNotFound, cfg.ErrMsgNotFound)
		return
	}

	c.room = room
	room.register <- c

	data, _ := json.Marshal(presentationStatus{
		PresentationID: params.PresentationID,
		Slides:         room.slides,
		CurrentIndex:   room.currentIndex,
	})
	c.writeData(data)
}

func (c *Client) writeData(data json.RawMessage) {
	resp, _ := json.Marshal(response.JSONResponse{Data: data})
	c.send <- resp
}

func (c *Client) writeError(code, message string) {
	resp, _ := json.Marshal(response.ErrorResponse{
		Error: response.ErrorData{Code: code, Message: message},
	})
	c.send <- resp
}
