package presence

import (
    "encoding/json"

    "github.com/timkrans/ws-framework/auth"
    "github.com/timkrans/ws-framework/events"
    "github.com/timkrans/ws-framework/transport"
)

var ModuleConfig Config

type Config struct {
    BroadcastJoinLeave bool
    BroadcastTyping    bool
    BroadcastStatus    bool
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
    ModuleConfig = cfg

    events.Register("presence.update", Handler{})
    events.Register("presence.typing", Handler{})

    events.Register("presence.away", Handler{})
    events.Register("presence.idle", Handler{})
    events.Register("presence.dnd", Handler{})
    events.Register("presence.back", Handler{})
    events.Register("presence.mobile", Handler{})
    events.Register("presence.status_text", Handler{})
}
