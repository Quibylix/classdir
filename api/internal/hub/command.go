package hub

import "encoding/json"

type CommandContext struct {
	Client *Client
	Room   *Room
	Hub    *Hub
}

type CommandHandler interface {
	Name() string
	Handle(ctx CommandContext, params json.RawMessage)
}

var clientHandlers = map[string]CommandHandler{}
var roomHandlers = map[string]CommandHandler{}

func register(h CommandHandler, roomScoped bool) {
	if roomScoped {
		roomHandlers[h.Name()] = h
	} else {
		clientHandlers[h.Name()] = h
	}
}
