package hub

import (
	"net/http"
)

func WSHandler(hub *Hub, acceptor wsAcceptor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := acceptor.Accept(w, r, nil)
		if err != nil {
			return
		}
		conn.SetReadLimit(maxMessageSize)
		client := NewClient(hub, conn)
		go client.WritePump()
		go client.ReadPump()
	}
}
