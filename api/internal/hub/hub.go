package hub

import (
	"sync"

	"classdir/api/internal/presentation"
)

const (
	CmdInitPresentation = "init_presentation"
	CmdJoinRoom         = "join_room"
	CmdNextSlide        = "next_slide"
	CmdPrevSlide        = "prev_slide"
	CmdGoToSlide        = "go_to_slide"
)

const (
	EventSlideChanged     = "slide_changed"
	EventPresentationInit = "presentation_initialized"
)

const channelBuffer = 256

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]*Room
	store presentation.Store
}

func NewHub(store presentation.Store) *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
		store: store,
	}
}

func (h *Hub) GetOrCreateRoom(id string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	if room, ok := h.rooms[id]; ok {
		return room
	}
	room := NewRoom(id)
	room.SetHub(h)
	h.rooms[id] = room
	go room.Run()
	return room
}

func (h *Hub) GetRoom(id string) *Room {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.rooms[id]
}

func (h *Hub) RemoveRoom(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.rooms, id)
}

func (h *Hub) Store() presentation.Store {
	return h.store
}
