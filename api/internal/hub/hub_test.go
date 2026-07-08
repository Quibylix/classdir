package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/time/rate"

	"classdir/api/internal/presentation"
	"classdir/api/internal/shared/cfg"
)

type mockConn struct {
	readCh    chan []byte
	writeCh   chan []byte
	doneCh    chan struct{}
	closeOnce sync.Once
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
		return 0, nil, context.Canceled
	case <-ctx.Done():
		return 0, nil, ctx.Err()
	}
}

func (m *mockConn) Write(ctx context.Context, typ websocket.MessageType, p []byte) error {
	select {
	case m.writeCh <- p:
		return nil
	case <-m.doneCh:
		return context.Canceled
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *mockConn) Close(code websocket.StatusCode, reason string) error {
	m.closeOnce.Do(func() {
		close(m.doneCh)
	})
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

func startClient(hub *Hub, conn *mockConn, rlp rateLimitProvider, cookies ...*http.Cookie) {
	req := httptest.NewRequest("GET", "/ws/v1", nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	WSHandler(hub, &mockAcceptor{conn: conn}, rlp)(httptest.NewRecorder(), req)
}

func setJWTSecret() {
	os.Setenv(cfg.EnvJWTSecret, "test-secret")
}

func signTestToken(t *testing.T) string {
	t.Helper()
	claims := jwt.RegisteredClaims{}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}
	return signed
}

func authCookie(t *testing.T) *http.Cookie {
	t.Helper()
	return &http.Cookie{Name: cfg.CookieName, Value: signTestToken(t)}
}

func validUUID() string { return "0192e5a0-7b7f-7b7f-8b7f-0192e5a07b7f" }

func validUUIDN(n int) string {
	return fmt.Sprintf("0192e5a0-7b7f-7b7f-8b7f-%012x", n)
}

func TestMain(m *testing.M) {
	os.Setenv(cfg.EnvJWTSecret, "test-secret")
	m.Run()
}

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

type initResponse struct {
	PresentationID string               `json:"presentation_id"`
	Slides         []presentation.Slide `json:"slides"`
	CurrentIndex   int                  `json:"current_index"`
	RoomCode       string               `json:"room_code"`
}

func initAndGetCode(t *testing.T, conn *mockConn) string {
	t.Helper()
	sendCommand(t, conn, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)
	data := recvData(t, conn)
	recvAnnotationsBatch(t, conn)
	var ir initResponse
	if err := json.Unmarshal(data, &ir); err != nil {
		t.Fatalf("failed to unmarshal init data: %v", err)
	}
	if ir.RoomCode == "" {
		t.Fatal("expected non-empty room_code")
	}
	return ir.RoomCode
}

// --- broadcast tests ---

func TestRoom_BroadcastsNextSlideToAllClients(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	viewerA := newMockConn()
	viewerB := newMockConn()

	startClient(hub, controller, DefaultRateLimitProvider{}, authCookie(t))
	startClient(hub, viewerA, DefaultRateLimitProvider{})
	startClient(hub, viewerB, DefaultRateLimitProvider{})

	code := initAndGetCode(t, controller)

	sendCommand(t, viewerA, CmdJoinRoom, `"room_code":"`+code+`"`)
	recvData(t, viewerA)
	recvAnnotationsBatch(t, viewerA)

	sendCommand(t, viewerB, CmdJoinRoom, `"room_code":"`+code+`"`)
	recvData(t, viewerB)
	recvAnnotationsBatch(t, viewerB)

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

	startClient(hub, controller, DefaultRateLimitProvider{}, authCookie(t))
	startClient(hub, viewer, DefaultRateLimitProvider{})

	code := initAndGetCode(t, controller)

	sendCommand(t, viewer, CmdJoinRoom, `"room_code":"`+code+`"`)
	recvData(t, viewer)
	recvAnnotationsBatch(t, viewer)

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

	startClient(hub, controller, DefaultRateLimitProvider{}, authCookie(t))
	startClient(hub, viewer, DefaultRateLimitProvider{})

	code := initAndGetCode(t, controller)

	sendCommand(t, viewer, CmdJoinRoom, `"room_code":"`+code+`"`)
	recvData(t, viewer)
	recvAnnotationsBatch(t, viewer)

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

	startClient(hub, controller, DefaultRateLimitProvider{}, authCookie(t))
	startClient(hub, viewer, DefaultRateLimitProvider{})

	code := initAndGetCode(t, controller)

	sendCommand(t, viewer, CmdJoinRoom, `"room_code":"`+code+`"`)
	recvData(t, viewer)
	recvAnnotationsBatch(t, viewer)

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

	startClient(hub, conn, DefaultRateLimitProvider{}, authCookie(t))

	sendCommand(t, conn, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)

	data := recvData(t, conn)
	recvAnnotationsBatch(t, conn)
	var ir initResponse
	if err := json.Unmarshal(data, &ir); err != nil {
		t.Fatalf("failed to unmarshal init data: %v", err)
	}
	if ir.PresentationID != validUUID() {
		t.Fatalf("expected presentation_id %s, got %s", validUUID(), ir.PresentationID)
	}
	if len(ir.Slides) != 3 {
		t.Fatalf("expected 3 slides, got %d", len(ir.Slides))
	}
	if ir.CurrentIndex != 0 {
		t.Fatalf("expected current_index 0, got %d", ir.CurrentIndex)
	}
	if ir.RoomCode == "" {
		t.Fatal("expected non-empty room_code")
	}
	if len(ir.RoomCode) != 8 {
		t.Fatalf("expected room_code to be 8 digits, got %q", ir.RoomCode)
	}
}

func TestClient_HandleInit_InvalidUUID(t *testing.T) {
	hub := newTestHub()
	conn := newMockConn()

	startClient(hub, conn, DefaultRateLimitProvider{}, authCookie(t))

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

	startClient(hub, conn, DefaultRateLimitProvider{}, authCookie(t))

	sendCommand(t, conn, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)
	recvError(t, conn, cfg.ErrNotFound)
}

func TestClient_HandleJoin_Valid(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	viewer := newMockConn()

	startClient(hub, controller, DefaultRateLimitProvider{}, authCookie(t))
	startClient(hub, viewer, DefaultRateLimitProvider{})

	code := initAndGetCode(t, controller)

	sendCommand(t, viewer, CmdJoinRoom, `"room_code":"`+code+`"`)
	data := recvData(t, viewer)
	recvAnnotationsBatch(t, viewer)

	var join struct {
		PresentationID string               `json:"presentation_id"`
		Slides         []presentation.Slide `json:"slides"`
		CurrentIndex   int                  `json:"current_index"`
		RoomCode       string               `json:"room_code"`
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
	if join.RoomCode != "" {
		t.Fatal("expected room_code to be empty in join response")
	}
}

func TestClient_HandleJoin_MissingRoom(t *testing.T) {
	hub := newTestHub()
	conn := newMockConn()

	startClient(hub, conn, DefaultRateLimitProvider{})

	sendCommand(t, conn, CmdJoinRoom, `"room_code":"00000000"`)
	recvError(t, conn, cfg.ErrNotFound)
}

func TestClient_HandleInit_Unauthenticated(t *testing.T) {
	hub := newTestHub()
	conn := newMockConn()

	startClient(hub, conn, DefaultRateLimitProvider{})

	sendCommand(t, conn, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)
	recvError(t, conn, cfg.ErrUnauthorized)
}

func TestClient_HandleInit_DuplicateInit(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	newController := newMockConn()

	startClient(hub, controller, DefaultRateLimitProvider{}, authCookie(t))
	startClient(hub, newController, DefaultRateLimitProvider{}, authCookie(t))

	initAndGetCode(t, controller)
	initAndGetCode(t, newController)

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

type mockRateLimitProvider struct{}

func (mockRateLimitProvider) Limits(authenticated bool) (rate.Limit, int) {
	return 0.001, 1
}

func TestReadPump_RateLimit_NonAuthenticated(t *testing.T) {
	hub := newTestHub()
	conn := newMockConn()
	startClient(hub, conn, mockRateLimitProvider{})

	sendCommand(t, conn, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)
	recvError(t, conn, cfg.ErrUnauthorized)

	sendCommand(t, conn, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)
	recvError(t, conn, cfg.ErrRateLimit)
}

func TestReadPump_RateLimit_Authenticated(t *testing.T) {
	hub := newTestHub()
	conn := newMockConn()
	startClient(hub, conn, mockRateLimitProvider{}, authCookie(t))

	sendCommand(t, conn, CmdInitPresentation, `"presentation_id":"not-a-uuid"`)
	recvError(t, conn, cfg.ErrInvalidUUID)

	sendCommand(t, conn, CmdInitPresentation, `"presentation_id":"not-a-uuid"`)
	recvError(t, conn, cfg.ErrRateLimit)
}

// --- annotation test helpers ---

type annotationAddedTestEvent struct {
	Event string `json:"event"`
	Data  struct {
		Type    string          `json:"type"`
		ID      string          `json:"id"`
		Payload json.RawMessage `json:"payload,omitempty"`
	} `json:"data"`
}

func recvAnnotationAdded(t *testing.T, conn *mockConn, expectedType string) annotationAddedTestEvent {
	t.Helper()
	select {
	case msg := <-conn.writeCh:
		var ev annotationAddedTestEvent
		if err := json.Unmarshal(msg, &ev); err != nil {
			t.Fatalf("failed to unmarshal annotation event: %v (raw: %s)", err, string(msg))
		}
		if ev.Event != EventAnnotationAdded {
			t.Fatalf("expected event %s, got %s", EventAnnotationAdded, ev.Event)
		}
		if ev.Data.Type != expectedType {
			t.Fatalf("expected type %s, got %s", expectedType, ev.Data.Type)
		}
		return ev
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for annotation_added event")
		return annotationAddedTestEvent{}
	}
}

type annotationsBatchTestEvent struct {
	Event string `json:"event"`
	Data  struct {
		OperationsBySlide map[string][]json.RawMessage `json:"operations_by_slide"`
	} `json:"data"`
}

func recvAnnotationsBatch(t *testing.T, conn *mockConn) annotationsBatchTestEvent {
	t.Helper()
	select {
	case msg := <-conn.writeCh:
		var ev annotationsBatchTestEvent
		if err := json.Unmarshal(msg, &ev); err != nil {
			t.Fatalf("failed to unmarshal batch event: %v (raw: %s)", err, string(msg))
		}
		if ev.Event != EventAnnotationsBatch {
			t.Fatalf("expected event %s, got %s", EventAnnotationsBatch, ev.Event)
		}
		return ev
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for annotations_batch event")
		return annotationsBatchTestEvent{}
	}
}

func TestAnnotation_StrokeAdded(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	viewer := newMockConn()

	startClient(hub, controller, DefaultRateLimitProvider{}, authCookie(t))
	startClient(hub, viewer, DefaultRateLimitProvider{})

	code := initAndGetCode(t, controller)

	sendCommand(t, viewer, CmdJoinRoom, `"room_code":"`+code+`"`)
	recvData(t, viewer)
	recvAnnotationsBatch(t, viewer)

	sendCommand(t, controller, CmdAnnotation, `"type":"stroke","id":"`+validUUIDN(1)+`","payload":{"points":[{"x":10,"y":20},{"x":30,"y":40}],"color":"#ff0000","thickness":3}`)

	evC := recvAnnotationAdded(t, controller, OpStroke)
	evV := recvAnnotationAdded(t, viewer, OpStroke)

	if evC.Data.ID != validUUIDN(1) {
		t.Fatalf("expected id %s, got %s", validUUIDN(1), evC.Data.ID)
	}
	if evC.Data.ID != evV.Data.ID {
		t.Fatal("controller and viewer got different ids")
	}
	if evC.Data.Payload == nil {
		t.Fatal("expected payload for stroke")
	}
	var payload AnnotationPayload
	if err := json.Unmarshal(evC.Data.Payload, &payload); err != nil {
		t.Fatalf("failed to unmarshal payload: %v", err)
	}
	if len(payload.Points) != 2 {
		t.Fatalf("expected 2 points, got %d", len(payload.Points))
	}
	if payload.Color != "#ff0000" {
		t.Fatalf("expected color #ff0000, got %s", payload.Color)
	}
	if payload.Thickness != 3 {
		t.Fatalf("expected thickness 3, got %f", payload.Thickness)
	}
}

func TestAnnotation_Cleared(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	viewer := newMockConn()

	startClient(hub, controller, DefaultRateLimitProvider{}, authCookie(t))
	startClient(hub, viewer, DefaultRateLimitProvider{})

	code := initAndGetCode(t, controller)

	sendCommand(t, viewer, CmdJoinRoom, `"room_code":"`+code+`"`)
	recvData(t, viewer)
	recvAnnotationsBatch(t, viewer)

	sendCommand(t, controller, CmdAnnotation, `"type":"stroke","id":"`+validUUIDN(2)+`","payload":{"points":[{"x":10,"y":20}],"color":"#00ff00","thickness":2}`)
	recvAnnotationAdded(t, controller, OpStroke)
	recvAnnotationAdded(t, viewer, OpStroke)

	sendCommand(t, controller, CmdAnnotation, `"type":"clear","id":"`+validUUIDN(3)+`"`)

	evC := recvAnnotationAdded(t, controller, OpClear)
	evV := recvAnnotationAdded(t, viewer, OpClear)

	if evC.Data.ID != validUUIDN(3) {
		t.Fatalf("expected id %s, got %s", validUUIDN(3), evC.Data.ID)
	}
	if evC.Data.Payload != nil {
		t.Fatal("expected no payload for clear")
	}
	if evV.Data.ID != evC.Data.ID {
		t.Fatal("controller and viewer got different ids")
	}
}

func TestAnnotation_ViewerIgnored(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	viewer := newMockConn()

	startClient(hub, controller, DefaultRateLimitProvider{}, authCookie(t))
	startClient(hub, viewer, DefaultRateLimitProvider{})

	code := initAndGetCode(t, controller)

	sendCommand(t, viewer, CmdJoinRoom, `"room_code":"`+code+`"`)
	recvData(t, viewer)
	recvAnnotationsBatch(t, viewer)

	sendCommand(t, viewer, CmdAnnotation, `"type":"stroke","id":"`+validUUIDN(4)+`","payload":{"points":[{"x":0,"y":0}],"color":"#000","thickness":1}`)
	sendCommand(t, controller, CmdAnnotation, `"type":"stroke","id":"`+validUUIDN(5)+`","payload":{"points":[{"x":0,"y":0}],"color":"#000","thickness":1}`)

	evC := recvAnnotationAdded(t, controller, OpStroke)
	if evC.Data.ID != validUUIDN(5) {
		t.Fatalf("expected id %s, got %s", validUUIDN(5), evC.Data.ID)
	}
}

func TestAnnotation_BatchOnJoin(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()

	startClient(hub, controller, DefaultRateLimitProvider{}, authCookie(t))

	code := initAndGetCode(t, controller)

	sendCommand(t, controller, CmdAnnotation, `"type":"stroke","id":"`+validUUIDN(10)+`","payload":{"points":[{"x":10,"y":20}],"color":"#ff0000","thickness":3}`)
	recvAnnotationAdded(t, controller, OpStroke)

	sendCommand(t, controller, CmdAnnotation, `"type":"clear","id":"`+validUUIDN(11)+`"`)
	recvAnnotationAdded(t, controller, OpClear)

	sendCommand(t, controller, CmdAnnotation, `"type":"stroke","id":"`+validUUIDN(12)+`","payload":{"points":[{"x":50,"y":60}],"color":"#0000ff","thickness":5}`)
	recvAnnotationAdded(t, controller, OpStroke)

	viewer := newMockConn()
	startClient(hub, viewer, DefaultRateLimitProvider{})

	sendCommand(t, viewer, CmdJoinRoom, `"room_code":"`+code+`"`)
	recvData(t, viewer)
	batch := recvAnnotationsBatch(t, viewer)

	if len(batch.Data.OperationsBySlide["0"]) != 3 {
		t.Fatalf("expected 3 operations for slide 0, got %d", len(batch.Data.OperationsBySlide["0"]))
	}
}

func TestAnnotation_OperationsPerSlide(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	viewer := newMockConn()

	startClient(hub, controller, DefaultRateLimitProvider{}, authCookie(t))
	startClient(hub, viewer, DefaultRateLimitProvider{})

	code := initAndGetCode(t, controller)

	sendCommand(t, viewer, CmdJoinRoom, `"room_code":"`+code+`"`)
	recvData(t, viewer)
	recvAnnotationsBatch(t, viewer)

	sendCommand(t, controller, CmdAnnotation, `"type":"stroke","id":"`+validUUIDN(20)+`","payload":{"points":[{"x":0,"y":0}],"color":"#f00","thickness":1}`)
	recvAnnotationAdded(t, controller, OpStroke)
	recvAnnotationAdded(t, viewer, OpStroke)

	sendCommand(t, controller, CmdGoToSlide, `"slide_number":1`)
	recvEvent(t, controller, EventSlideChanged)
	recvEvent(t, viewer, EventSlideChanged)

	sendCommand(t, controller, CmdAnnotation, `"type":"stroke","id":"`+validUUIDN(21)+`","payload":{"points":[{"x":100,"y":100}],"color":"#0f0","thickness":2}`)
	recvAnnotationAdded(t, controller, OpStroke)
	recvAnnotationAdded(t, viewer, OpStroke)

	viewer2 := newMockConn()
	startClient(hub, viewer2, DefaultRateLimitProvider{})

	sendCommand(t, viewer2, CmdJoinRoom, `"room_code":"`+code+`"`)
	recvData(t, viewer2)
	batch := recvAnnotationsBatch(t, viewer2)

	if len(batch.Data.OperationsBySlide["0"]) != 1 {
		t.Fatalf("expected 1 operation for slide 0, got %d", len(batch.Data.OperationsBySlide["0"]))
	}
	if len(batch.Data.OperationsBySlide["1"]) != 1 {
		t.Fatalf("expected 1 operation for slide 1, got %d", len(batch.Data.OperationsBySlide["1"]))
	}
}

func TestClient_ReconnectPreservesCurrentIndex(t *testing.T) {
	hub := newTestHub()

	controller := newMockConn()
	viewer := newMockConn()
	startClient(hub, controller, DefaultRateLimitProvider{}, authCookie(t))
	startClient(hub, viewer, DefaultRateLimitProvider{})

	code := initAndGetCode(t, controller)

	sendCommand(t, viewer, CmdJoinRoom, `"room_code":"`+code+`"`)
	recvData(t, viewer)
	recvAnnotationsBatch(t, viewer)

	sendCommand(t, controller, CmdNextSlide, "")
	recvEvent(t, controller, EventSlideChanged)
	recvEvent(t, viewer, EventSlideChanged)

	controller.Close(websocket.StatusNormalClosure, "test disconnect")

	reconnector := newMockConn()
	startClient(hub, reconnector, DefaultRateLimitProvider{}, authCookie(t))
	sendCommand(t, reconnector, CmdInitPresentation, `"presentation_id":"`+validUUID()+`"`)

	data := recvData(t, reconnector)
	recvAnnotationsBatch(t, reconnector)
	var s struct {
		CurrentIndex int `json:"current_index"`
	}
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatalf("failed to unmarshal init data: %v", err)
	}
	if s.CurrentIndex != 1 {
		t.Fatalf("expected current_index 1 (preserved), got %d", s.CurrentIndex)
	}

	sendCommand(t, reconnector, CmdNextSlide, "")
	ev := recvEvent(t, reconnector, EventSlideChanged)
	if ev.Data.CurrentSlide != 2 {
		t.Fatalf("expected current_slide 2, got %d", ev.Data.CurrentSlide)
	}
}
