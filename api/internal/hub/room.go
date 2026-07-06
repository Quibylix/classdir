package hub

import (
	"encoding/json"

	"classdir/api/internal/presentation"
)

type roomCommand struct {
	msg    []byte
	sender *Client
}

type Room struct {
	ID           string
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
					r.hub.RemoveRoom(r.ID)
					return
				}
			}
		case cmd := <-r.commands:
			r.handleCommand(cmd)
		}
	}
}

func (r *Room) handleCommand(cmd roomCommand) {
	if cmd.sender != r.controller {
		return
	}

	var parsed Command
	if err := json.Unmarshal(cmd.msg, &parsed); err != nil {
		return
	}

	switch parsed.Command {
	case CmdNextSlide:
		if r.currentIndex < len(r.slides)-1 {
			r.currentIndex++
		}
	case CmdPrevSlide:
		if r.currentIndex > 0 {
			r.currentIndex--
		}
	case CmdGoToSlide:
		var params SlideParams
		if err := json.Unmarshal(parsed.Parameters, &params); err != nil {
			return
		}
		if params.SlideNumber < 0 || params.SlideNumber >= len(r.slides) {
			return
		}
		r.currentIndex = params.SlideNumber
	default:
		return
	}

	type SlideChangedEventData struct {
		CurrentSlide int `json:"current_slide"`
	}

	type SlideChangedEvent struct {
		Event string                `json:"event"`
		Data  SlideChangedEventData `json:"data"`
	}

	event, _ := json.Marshal(SlideChangedEvent{
		Event: EventSlideChanged,
		Data: SlideChangedEventData{
			CurrentSlide: r.currentIndex,
		},
	})

	for client := range r.clients {
		select {
		case client.send <- event:
		default:
		}
	}
}

func (r *Room) SetHub(h *Hub) {
	r.hub = h
}
