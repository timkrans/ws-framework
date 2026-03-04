package reactions

import (
    "encoding/json"

    "github.com/timkrans/ws-framework/auth"
    "github.com/timkrans/ws-framework/events"
    "github.com/timkrans/ws-framework/transport"
)

type Config struct {
    AllowReactions bool
}

var ModuleConfig Config

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
    ModuleConfig = cfg

    if cfg.AllowReactions {
        events.Register("reaction.add", Handler{})
        events.Register("reaction.remove", Handler{})
    }
}
