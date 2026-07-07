package hub

import (
	"context"
	"encoding/json"

	"github.com/coder/websocket"

	"classdir/api/internal/shared/response"
)

const maxMessageSize = 65536

type Command struct {
	Command    string          `json:"command"`
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

type Client struct {
	hub  *Hub
	conn wsConn
	send chan []byte
	room *Room

	Authenticated bool
}

func NewClient(hub *Hub, conn wsConn) *Client {
	return &Client{
		hub:           hub,
		conn:          conn,
		send:          make(chan []byte, channelBuffer),
		Authenticated: false,
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

		c.handleCommand(cmd)
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

func (c *Client) handleCommand(cmd Command) {
	if h, ok := clientHandlers[cmd.Command]; ok {
		h.Handle(CommandContext{Client: c, Hub: c.hub}, cmd.Parameters)
		return
	}
	if c.room != nil {
		if h, ok := roomHandlers[cmd.Command]; ok {
			c.room.commands <- roomCommand{handler: h, params: cmd.Parameters, sender: c}
		}
	}
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
