package hub

import (
	"context"
	"net/http"

	"github.com/coder/websocket"
	"golang.org/x/time/rate"
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

type DefaultAcceptor struct {
	OriginPatterns []string
}

func (a DefaultAcceptor) Accept(w http.ResponseWriter, r *http.Request, opts *websocket.AcceptOptions) (wsConn, error) {
	if opts == nil {
		opts = &websocket.AcceptOptions{}
	}
	opts.OriginPatterns = a.OriginPatterns
	return websocket.Accept(w, r, opts)
}

type rateLimitProvider interface {
	Limits(authenticated bool) (rate.Limit, int)
}

type DefaultRateLimitProvider struct{}

func (DefaultRateLimitProvider) Limits(authenticated bool) (rate.Limit, int) {
	if authenticated {
		return 20, 50
	}
	return 2, 5
}
