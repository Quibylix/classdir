package hub

import (
	"context"
	"net/http"

	"github.com/coder/websocket"
)

type wsConn interface {
	Read(ctx context.Context) (websocket.MessageType, []byte, error)
	Write(ctx context.Context, typ websocket.MessageType, p []byte) error
	Close(code websocket.StatusCode, reason string) error
	SetReadLimit(n int64)
}

type wsAcceptor interface {
	Accept(w http.ResponseWriter, r *http.Request, opts *websocket.AcceptOptions) (wsConn, error)
}

type DefaultAcceptor struct{}

func (DefaultAcceptor) Accept(w http.ResponseWriter, r *http.Request, opts *websocket.AcceptOptions) (wsConn, error) {
	return websocket.Accept(w, r, opts)
}
