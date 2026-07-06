package hub

import "encoding/json"

type SlideParams struct {
	SlideNumber int `json:"slide_number"`
}

type NextSlideHandler struct{}

func (h NextSlideHandler) Name() string { return CmdNextSlide }

func (h NextSlideHandler) Handle(ctx CommandContext, params json.RawMessage) {
	if ctx.Room.currentIndex < len(ctx.Room.slides)-1 {
		ctx.Room.currentIndex++
	}
	ctx.Room.broadcastSlideChanged()
}

type PrevSlideHandler struct{}

func (h PrevSlideHandler) Name() string { return CmdPrevSlide }

func (h PrevSlideHandler) Handle(ctx CommandContext, params json.RawMessage) {
	if ctx.Room.currentIndex > 0 {
		ctx.Room.currentIndex--
	}
	ctx.Room.broadcastSlideChanged()
}

type GoToSlideHandler struct{}

func (h GoToSlideHandler) Name() string { return CmdGoToSlide }

func (h GoToSlideHandler) Handle(ctx CommandContext, params json.RawMessage) {
	var p SlideParams
	if err := json.Unmarshal(params, &p); err != nil {
		return
	}
	if p.SlideNumber < 0 || p.SlideNumber >= len(ctx.Room.slides) {
		return
	}
	ctx.Room.currentIndex = p.SlideNumber
	ctx.Room.broadcastSlideChanged()
}

func init() {
	register(NextSlideHandler{}, true)
	register(PrevSlideHandler{}, true)
	register(GoToSlideHandler{}, true)
}
