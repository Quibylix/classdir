package hub

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"

	"classdir/api/internal/presentation"
	"classdir/api/internal/shared/cfg"
)

type mockConn struct {
	readCh  chan []byte
	writeCh chan []byte
	doneCh  chan struct{}
}

func newMockConn() *mockConn {
	return &mockConn{
		readCh:  make(chan []byte, 256),
		writeCh: make(chan []byte, 256),
		doneCh:  make(chan struct{}),
	}
}

func (m *mockConn) Read(ctx context.Context) (websocket.MessageType, []byte, error) {
	select {
	case msg := <-m.readCh:
		return websocket.MessageText, msg, nil
	case <-m.doneCh:
		return 0, nil, ctx.Err()
	case <-ctx.Done():
		return 0, nil, ctx.Err()
	}
}

func (m *mockConn) Write(ctx context.Context, typ websocket.MessageType, p []byte) error {
	select {
	case m.writeCh <- p:
		return nil
	case <-m.doneCh:
		return ctx.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *mockConn) Close(code websocket.StatusCode, reason string) error {
	close(m.doneCh)
	return nil
}

func (m *mockConn) SetReadLimit(n int64) {}

type mockAcceptor struct {
	conn wsConn
}

func (a *mockAcceptor) Accept(w http.ResponseWriter, r *http.Request, opts *websocket.AcceptOptions) (wsConn, error) {
	return a.conn, nil
}

type mockStore struct {
	presentation.Store
	getByIDFunc func(ctx context.Context, id string) (*presentation.Presentation, error)
}

func (m *mockStore) GetByID(ctx context.Context, id string) (*presentation.Presentation, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func startClient(hub *Hub, conn *mockConn) {
	WSHandler(hub, &mockAcceptor{conn: conn})(httptest.NewRecorder(), httptest.NewRequest("GET", "/ws/v1", nil))
}

func validUUID() string { return "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f" }

func testSlides() []presentation.Slide {
	return []presentation.Slide{
		{ID: "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b70", Content: "# Slide 1"},
		{ID: "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b71", Content: "# Slide 2"},
		{ID: "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b72", Content: "# Slide 3"},
	}
}

func newTestHub() *Hub {
	return NewHub(&mockStore{
		getByIDFunc: func(ctx context.Context, id string) (*presentation.Presentation, error) {
			return &presentation.Presentation{
				ID:     validUUID(),
				Title:  "Test Pres",
				Slides: testSlides(),
			}, nil
		},
	})
}

type dataResponse struct {
	Data json.RawMessage `json:"data"`
}

type slideChangedEvent struct {
	Event string `json:"event"`
	Data  struct {
		CurrentSlide int `json:"current_slide"`
	} `json:"data"`
}

type errorResponse struct {
	Error struct {
		Code string `json:"code"`
	} `json:"error"`
}

func sendCommand(t *testing.T, conn *mockConn, cmd, params string) {
	t.Helper()
	var raw string
	if params != "" {
		raw = `{"command":"` + cmd + `","parameters":{` + params + `}}`
	} else {
		raw = `{"command":"` + cmd + `"}`
	}
	conn.readCh <- []byte(raw)
}

func recvData(t *testing.T, conn *mockConn) json.RawMessage {
	t.Helper()
	select {
	case msg := <-conn.writeCh:
		var resp dataResponse
		if err := json.Unmarshal(msg, &resp); err != nil {
			t.Fatalf("failed to unmarshal data response: %v", err)
		}
		if resp.Data == nil {
			t.Fatal("expected data in response")
		}
		return resp.Data
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for data response")
		return nil
	}
}

func recvEvent(t *testing.T, conn *mockConn, expectedEvent string) slideChangedEvent {
	t.Helper()
	select {
	case msg := <-conn.writeCh:
		var ev slideChangedEvent
		if err := json.Unmarshal(msg, &ev); err != nil {
			t.Fatalf("failed to unmarshal event: %v (raw: %s)", err, string(msg))
		}
		if ev.Event != expectedEvent {
			t.Fatalf("expected event %s, got %s", expectedEvent, ev.Event)
		}
		return ev
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for event " + expectedEvent)
		return slideChangedEvent{}
	}
}

func recvError(t *testing.T, conn *mockConn, expectedCode string) {
	t.Helper()
	select {
	case msg := <-conn.writeCh:
		var errResp errorResponse
		if err := json.Unmarshal(msg, &errResp); err != nil {
			t.Fatalf("failed to unmarshal error: %v", err)
		}
		if errResp.Error.Code != expectedCode {
			t.Fatalf("expected error code %s, got %s", expectedCode, errResp.Error.Code)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for error")
	}
}

// --- broadcast tests ---

func TestRoom_BroadcastsNextSlideToAllClients(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	viewerA := newMockConn()
	viewerB := newMockConn()

	startClient(hub, controller)
	startClient(hub, viewerA)
	startClient(hub, viewerB)

	sendCommand(t, controller, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)
	recvData(t, controller)

	sendCommand(t, viewerA, CmdJoinRoom, `"presentation_id":"`+validUUID()+`"`)
	recvData(t, viewerA)

	sendCommand(t, viewerB, CmdJoinRoom, `"presentation_id":"`+validUUID()+`"`)
	recvData(t, viewerB)

	sendCommand(t, controller, CmdNextSlide, "")

	evC := recvEvent(t, controller, EventSlideChanged)
	evA := recvEvent(t, viewerA, EventSlideChanged)
	evB := recvEvent(t, viewerB, EventSlideChanged)

	if evC.Data.CurrentSlide != 1 {
		t.Fatalf("controller expected current_slide 1, got %d", evC.Data.CurrentSlide)
	}
	if evA.Data.CurrentSlide != evC.Data.CurrentSlide {
		t.Fatal("viewer A got different current_slide than controller")
	}
	if evB.Data.CurrentSlide != evC.Data.CurrentSlide {
		t.Fatal("viewer B got different current_slide than controller")
	}
}

func TestRoom_BroadcastsPrevSlideToAllClients(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	viewer := newMockConn()

	startClient(hub, controller)
	startClient(hub, viewer)

	sendCommand(t, controller, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)
	recvData(t, controller)

	sendCommand(t, viewer, CmdJoinRoom, `"presentation_id":"`+validUUID()+`"`)
	recvData(t, viewer)

	sendCommand(t, controller, CmdNextSlide, "")
	recvEvent(t, controller, EventSlideChanged)
	recvEvent(t, viewer, EventSlideChanged)

	sendCommand(t, controller, CmdNextSlide, "")
	recvEvent(t, controller, EventSlideChanged)
	recvEvent(t, viewer, EventSlideChanged)

	sendCommand(t, controller, CmdPrevSlide, "")

	evC := recvEvent(t, controller, EventSlideChanged)
	evV := recvEvent(t, viewer, EventSlideChanged)

	if evC.Data.CurrentSlide != 1 {
		t.Fatalf("controller expected current_slide 1, got %d", evC.Data.CurrentSlide)
	}
	if evV.Data.CurrentSlide != 1 {
		t.Fatalf("viewer expected current_slide 1, got %d", evV.Data.CurrentSlide)
	}
}

func TestRoom_BroadcastsGoToSlideToAllClients(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	viewer := newMockConn()

	startClient(hub, controller)
	startClient(hub, viewer)

	sendCommand(t, controller, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)
	recvData(t, controller)

	sendCommand(t, viewer, CmdJoinRoom, `"presentation_id":"`+validUUID()+`"`)
	recvData(t, viewer)

	sendCommand(t, controller, CmdGoToSlide, `"slide_number":1`)

	evC := recvEvent(t, controller, EventSlideChanged)
	evV := recvEvent(t, viewer, EventSlideChanged)

	if evC.Data.CurrentSlide != 1 {
		t.Fatalf("controller expected current_slide 1, got %d", evC.Data.CurrentSlide)
	}
	if evV.Data.CurrentSlide != 1 {
		t.Fatalf("viewer expected current_slide 1, got %d", evV.Data.CurrentSlide)
	}
}

func TestRoom_ViewerCommandDoesNotBroadcast(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	viewer := newMockConn()

	startClient(hub, controller)
	startClient(hub, viewer)

	sendCommand(t, controller, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)
	recvData(t, controller)

	sendCommand(t, viewer, CmdJoinRoom, `"presentation_id":"`+validUUID()+`"`)
	recvData(t, viewer)

	sendCommand(t, viewer, CmdGoToSlide, `"slide_number":2`)
	sendCommand(t, controller, CmdNextSlide, "")

	evC := recvEvent(t, controller, EventSlideChanged)
	evV := recvEvent(t, viewer, EventSlideChanged)

	if evC.Data.CurrentSlide != 1 {
		t.Fatalf("expected current_slide 1 (not 2 — viewer's cmd was ignored), got %d", evC.Data.CurrentSlide)
	}
	if evV.Data.CurrentSlide != 1 {
		t.Fatalf("expected current_slide 1, got %d", evV.Data.CurrentSlide)
	}
}

// --- init/join tests ---

func TestClient_HandleInit_Valid(t *testing.T) {
	hub := newTestHub()
	conn := newMockConn()

	startClient(hub, conn)

	sendCommand(t, conn, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)

	data := recvData(t, conn)
	var init struct {
		PresentationID string               `json:"presentation_id"`
		Slides         []presentation.Slide `json:"slides"`
		CurrentIndex   int                  `json:"current_index"`
	}
	if err := json.Unmarshal(data, &init); err != nil {
		t.Fatalf("failed to unmarshal init data: %v", err)
	}
	if init.PresentationID != validUUID() {
		t.Fatalf("expected presentation_id %s, got %s", validUUID(), init.PresentationID)
	}
	if len(init.Slides) != 3 {
		t.Fatalf("expected 3 slides, got %d", len(init.Slides))
	}
	if init.CurrentIndex != 0 {
		t.Fatalf("expected current_index 0, got %d", init.CurrentIndex)
	}
}

func TestClient_HandleInit_InvalidUUID(t *testing.T) {
	hub := newTestHub()
	conn := newMockConn()

	startClient(hub, conn)

	sendCommand(t, conn, CmdInitPresentation, `"presentation_id":"not-a-uuid"`)
	recvError(t, conn, cfg.ErrInvalidUUID)
}

func TestClient_HandleInit_NilPresentation(t *testing.T) {
	hub := NewHub(&mockStore{
		getByIDFunc: func(ctx context.Context, id string) (*presentation.Presentation, error) {
			return nil, nil
		},
	})
	conn := newMockConn()

	startClient(hub, conn)

	sendCommand(t, conn, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)
	recvError(t, conn, cfg.ErrNotFound)
}

func TestClient_HandleJoin_Valid(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	viewer := newMockConn()

	startClient(hub, controller)
	startClient(hub, viewer)

	sendCommand(t, controller, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)
	recvData(t, controller)

	sendCommand(t, viewer, CmdJoinRoom, `"presentation_id":"`+validUUID()+`"`)
	data := recvData(t, viewer)

	var join struct {
		PresentationID string               `json:"presentation_id"`
		Slides         []presentation.Slide `json:"slides"`
		CurrentIndex   int                  `json:"current_index"`
	}
	if err := json.Unmarshal(data, &join); err != nil {
		t.Fatalf("failed to unmarshal join data: %v", err)
	}
	if len(join.Slides) != 3 {
		t.Fatalf("expected 3 slides, got %d", len(join.Slides))
	}
	if join.CurrentIndex != 0 {
		t.Fatalf("expected current_index 0, got %d", join.CurrentIndex)
	}
}

func TestClient_HandleJoin_MissingRoom(t *testing.T) {
	hub := newTestHub()
	conn := newMockConn()

	startClient(hub, conn)

	sendCommand(t, conn, CmdJoinRoom, `"presentation_id":"`+validUUID()+`"`)
	recvError(t, conn, cfg.ErrNotFound)
}

func TestClient_HandleJoin_InvalidUUID(t *testing.T) {
	hub := newTestHub()
	conn := newMockConn()

	startClient(hub, conn)

	sendCommand(t, conn, CmdJoinRoom, `"presentation_id":"not-a-uuid"`)
	recvError(t, conn, cfg.ErrInvalidUUID)
}

func TestClient_HandleInit_DuplicateInit(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	newController := newMockConn()

	startClient(hub, controller)
	startClient(hub, newController)

	sendCommand(t, controller, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)
	recvData(t, controller)

	sendCommand(t, newController, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)
	recvData(t, newController)

	sendCommand(t, controller, CmdGoToSlide, `"slide_number":2`)
	sendCommand(t, newController, CmdGoToSlide, `"slide_number":1`)

	evC := recvEvent(t, controller, EventSlideChanged)
	evI := recvEvent(t, newController, EventSlideChanged)

	if evC.Data.CurrentSlide != 1 {
		t.Fatalf("expected current_slide 1 (not 2 — new controller's cmd was ignored), got %d", evC.Data.CurrentSlide)
	}
	if evI.Data.CurrentSlide != 1 {
		t.Fatalf("expected current_slide 1, got %d", evI.Data.CurrentSlide)
	}
}
