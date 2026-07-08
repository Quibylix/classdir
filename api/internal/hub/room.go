package hub

import (
	"encoding/json"
	"time"

	"classdir/api/internal/presentation"
)

type roomCommand struct {
	handler CommandHandler
	params  json.RawMessage
	sender  *Client
}

type Room struct {
	ID                string
	Code              string
	clients           map[*Client]bool
	controller        *Client
	register          chan *Client
	unregister        chan *Client
	done              chan struct{}
	commands          chan roomCommand
	currentIndex      int
	slides            []presentation.Slide
	hub               *Hub
	operationsBySlide map[int][]AnnotationOperation
}

const roomDeleteTimeout = 1 * time.Minute

func NewRoom(id string) *Room {
	return &Room{
		ID:                id,
		clients:           make(map[*Client]bool),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		commands:          make(chan roomCommand, channelBuffer),
		done:              make(chan struct{}),
		operationsBySlide: make(map[int][]AnnotationOperation),
	}
}

func (r *Room) Run() {
	var (
		deleteTimer *time.Timer
		deleteCh    <-chan time.Time
	)

	for {
		select {
		case client := <-r.register:
			if deleteTimer != nil {
				deleteTimer.Stop()
				deleteTimer = nil
				deleteCh = nil
			}
			r.clients[client] = true

		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
				if client == r.controller {
					r.controller = nil
				}
				if len(r.clients) == 0 && r.hub != nil && deleteTimer == nil {
					deleteTimer = time.NewTimer(roomDeleteTimeout)
					deleteCh = deleteTimer.C
				}
			}

		case <-deleteCh:
			if len(r.clients) == 0 && r.hub != nil {
				close(r.done)
				r.hub.RemoveRoom(r)
				return
			}
			deleteTimer = nil
			deleteCh = nil

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

func (r *Room) broadcastAnnotationAdded(op AnnotationOperation) {
	e, err := json.Marshal(annotationAddedEvent{
		Event: EventAnnotationAdded,
		Data: annotationAddedData{
			Type:    op.Type,
			ID:      op.ID,
			Payload: op.Payload,
		},
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

func (r *Room) sendAnnotationsBatch(client *Client) {
	if r.operationsBySlide == nil {
		return
	}
	e, err := json.Marshal(annotationsBatchEvent{
		Event: EventAnnotationsBatch,
		Data: annotationsBatchData{
			OperationsBySlide: r.operationsBySlide,
		},
	})
	if err != nil {
		return
	}
	select {
	case client.send <- e:
	default:
	}
}

func (r *Room) SetHub(h *Hub) {
	r.hub = h
}
