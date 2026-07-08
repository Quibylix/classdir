package hub

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"

	"classdir/api/internal/presentation"
)

const (
	CmdInitPresentation = "init_presentation"
	CmdJoinRoom         = "join_room"
	CmdNextSlide        = "next_slide"
	CmdPrevSlide        = "prev_slide"
	CmdGoToSlide        = "go_to_slide"
	CmdAnnotation       = "annotation"
)

const (
	EventSlideChanged     = "slide_changed"
	EventPresentationInit = "presentation_initialized"
	EventAnnotationAdded  = "annotation_added"
	EventAnnotationsBatch = "annotations_batch"
)

const channelBuffer = 256

type Hub struct {
	mu          sync.RWMutex
	rooms       map[string]*Room
	roomsByCode map[string]*Room
	store       presentation.Store
}

func NewHub(store presentation.Store) *Hub {
	return &Hub{
		rooms:       make(map[string]*Room),
		roomsByCode: make(map[string]*Room),
		store:       store,
	}
}

func (h *Hub) GetOrCreateRoom(id string) (*Room, error) {
	pres, err := h.store.GetByID(context.Background(), id)

	if err != nil {
		return nil, err
	}

	if pres == nil {
		return nil, presentation.ErrNotFound
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	if existingRoom, ok := h.rooms[id]; ok {
		return existingRoom, nil
	}

	room := NewRoom(id, h.generateCodeLocked(), h, pres.Slides)
	h.rooms[id] = room
	h.roomsByCode[room.Code] = room
	go room.Run()
	return room, nil
}

func (h *Hub) GetRoomByCode(code string) *Room {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.roomsByCode[code]
}

func (h *Hub) RemoveRoom(room *Room) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.rooms, room.ID)
	delete(h.roomsByCode, room.Code)
}

func (h *Hub) Store() presentation.Store {
	return h.store
}

func (h *Hub) generateCodeLocked() string {
	for {
		n, err := rand.Int(rand.Reader, big.NewInt(100000000))
		if err != nil {
			continue
		}
		code := fmt.Sprintf("%08d", n.Int64())
		if _, ok := h.roomsByCode[code]; !ok {
			return code
		}
	}
}
