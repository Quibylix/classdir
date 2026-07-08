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

type registrationRequest struct {
	client       *Client
	isController bool
}

type Room struct {
	ID                string
	Code              string
	clients           map[*Client]bool
	controller        *Client
	register          chan *registrationRequest
	unregister        chan *Client
	done              chan struct{}
	commands          chan roomCommand
	currentIndex      int
	slides            []presentation.Slide
	hub               *Hub
	operationsBySlide map[int][]AnnotationOperation
}

const roomDeleteTimeout = 1 * time.Minute

func NewRoom(id, code string, hub *Hub, slides []presentation.Slide) *Room {
	return &Room{
		ID:                id,
		Code:              code,
		clients:           make(map[*Client]bool),
		register:          make(chan *registrationRequest),
		unregister:        make(chan *Client),
		commands:          make(chan roomCommand, channelBuffer),
		done:              make(chan struct{}),
		operationsBySlide: make(map[int][]AnnotationOperation),
		currentIndex:      0,
		hub:               hub,
		slides:            slides,
	}
}

func (r *Room) Run() {
	var (
		deleteTimer *time.Timer
		deleteCh    <-chan time.Time
	)

	for {
		select {
		case req := <-r.register:
			if deleteTimer != nil {
				deleteTimer.Stop()
				deleteTimer = nil
				deleteCh = nil
			}
			r.handleClientRegistration(req)

		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
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

func (r *Room) handleClientRegistration(req *registrationRequest) {
	if req.isController {
		r.controller = req.client
	}
	r.clients[req.client] = true

	initResponse := struct {
		Data presentationStatus `json:"data"`
	}{
		Data: presentationStatus{
			PresentationID: r.ID,
			Slides:         r.slides,
			CurrentIndex:   r.currentIndex,
			RoomCode:       "",
		},
	}

	if req.isController {
		initResponse.Data.RoomCode = r.Code
	}

	data, _ := json.Marshal(initResponse)

	select {
	case req.client.send <- data:
	default:
	}

	r.sendAnnotationsBatch(req.client)
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
