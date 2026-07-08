package hub

import (
	"encoding/json"

	"classdir/api/internal/shared/validate"
)

type AnnotationHandler struct{}

func (h AnnotationHandler) Name() string { return CmdAnnotation }

func (h AnnotationHandler) Handle(ctx CommandContext, params json.RawMessage) {
	var p AnnotationParams
	if err := json.Unmarshal(params, &p); err != nil {
		return
	}
	if p.Type != OpStroke && p.Type != OpClear {
		return
	}
	if !validate.IsValidUUIDv7(p.ID) {
		return
	}
	if p.Type == OpStroke && p.Payload == nil {
		return
	}

	slide := ctx.Room.currentIndex
	for _, op := range ctx.Room.operationsBySlide[slide] {
		if op.ID == p.ID {
			return
		}
	}

	op := AnnotationOperation{
		Type:    p.Type,
		ID:      p.ID,
		Payload: p.Payload,
	}

	ctx.Room.operationsBySlide[slide] = append(ctx.Room.operationsBySlide[slide], op)

	ctx.Room.broadcastAnnotationAdded(op)
}

func init() {
	register(AnnotationHandler{}, true)
}
