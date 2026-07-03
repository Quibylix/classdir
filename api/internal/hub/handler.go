package hub

import (
	"net/http"

	"github.com/coder/websocket"
)

func WSHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		conn.SetReadLimit(maxMessageSize)
		client := NewClient(hub, conn)
		go client.WritePump()
		go client.ReadPump()
	}
}
