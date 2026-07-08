package hub

import (
	"context"
	"encoding/json"

	"github.com/coder/websocket"
	"golang.org/x/time/rate"

	"classdir/api/internal/shared/cfg"
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
	limiter       *rate.Limiter

	ctx    context.Context
	cancel context.CancelFunc
}

func NewClient(hub *Hub, conn wsConn) *Client {
	cxt, cancel := context.WithCancel(context.Background())

	return &Client{
		hub:           hub,
		conn:          conn,
		send:          make(chan []byte, channelBuffer),
		Authenticated: false,
		ctx:           cxt,
		cancel:        cancel,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		if c.room != nil {
			select {
			case c.room.unregister <- c:
			case <-c.room.done:
			}
		}
		c.cancel()
		c.conn.Close(websocket.StatusNormalClosure, "connection closed")
	}()

	for {
		_, msg, err := c.conn.Read(c.ctx)
		if err != nil {
			break
		}

		if !c.limiter.Allow() {
			c.writeError(cfg.ErrRateLimit, cfg.ErrMsgRateLimit)
			continue
		}

		var cmd Command
		if err := json.Unmarshal(msg, &cmd); err != nil {
			continue
		}

		c.handleCommand(cmd)
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.cancel()
		c.conn.Close(websocket.StatusNormalClosure, "connection closed")
	}()

	for {
		select {
		case msg := <-c.send:
			err := c.conn.Write(c.ctx, websocket.MessageText, msg)
			if err != nil {
				return
			}
		case <-c.ctx.Done():
			return
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
			select {
			case c.room.commands <- roomCommand{handler: h, params: cmd.Parameters, sender: c}:
			case <-c.room.done:
			}
		}
	}
}

func (c *Client) writeData(data json.RawMessage) {
	resp, _ := json.Marshal(response.JSONResponse{Data: data})
	select {
	case c.send <- resp:
	case <-c.ctx.Done():
	}
}

func (c *Client) writeError(code, message string) {
	resp, _ := json.Marshal(response.ErrorResponse{
		Error: response.ErrorData{Code: code, Message: message},
	})
	select {
	case c.send <- resp:
	default:
	}
}
