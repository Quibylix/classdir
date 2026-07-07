package hub

import (
	"encoding/json"

	"classdir/api/internal/presentation"
)

type roomCommand struct {
	handler CommandHandler
	params  json.RawMessage
	sender  *Client
}

type Room struct {
	ID           string
	Code         string
	clients      map[*Client]bool
	controller   *Client
	register     chan *Client
	unregister   chan *Client
	commands     chan roomCommand
	currentIndex int
	slides       []presentation.Slide
	hub          *Hub
}

func NewRoom(id string) *Room {
	return &Room{
		ID:         id,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		commands:   make(chan roomCommand, channelBuffer),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
				if client == r.controller {
					r.controller = nil
				}
				if len(r.clients) == 0 && r.hub != nil {
					r.hub.RemoveRoom(r)
					return
				}
			}
		case cmd := <-r.commands:
			if cmd.sender == r.controller {
				cmd.handler.Handle(CommandContext{Client: cmd.sender, Room: r, Hub: r.hub}, cmd.params)
			}
		}
	}
}

func (r *Room) broadcastSlideChanged() {
	type data struct {
		CurrentSlide int `json:"current_slide"`
	}
	type event struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	e, err := json.Marshal(event{
		Event: EventSlideChanged,
		Data:  data{CurrentSlide: r.currentIndex},
	})

	if err != nil {
		return
	}

	for client := range r.clients {
		select {
		case client.send <- e:
		default:
		}
	}
}

func (r *Room) SetHub(h *Hub) {
	r.hub = h
}
