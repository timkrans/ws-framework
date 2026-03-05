package userprofile

import (
    "encoding/json"
	"github.com/timkrans/ws-framework/auth"
    "github.com/timkrans/ws-framework/events"
    "github.com/timkrans/ws-framework/transport"
)

var Config Config

type Config struct {
    BroadcastChanges bool
}

type Handler struct{}

func (h Handler) Handle(c interface{}, evt events.Event) {
    client := c.(*transport.Client)
	switch client.Source.(type) {
	case *auth.RemoteAuth:
		if evt.User != client.UserID {
			return
		}
	default:
		evt.User = client.UserID
	}



    out, _ := json.Marshal(evt)
    client.Room.Broadcast <- out
}

func Init(cfg Config) {
    Config = cfg
    events.Register("user.update", Handler{})
}
