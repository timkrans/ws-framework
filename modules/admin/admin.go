package admin

import (
    "encoding/json"

    "github.com/timkrans/ws-framework/auth"
    "github.com/timkrans/ws-framework/events"
    "github.com/timkrans/ws-framework/transport"
)

var ModuleConfig Config

type Config struct {
    AllowRemoveUser bool
    AllowOffboard   bool
    AllowReactivate bool
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

    events.Register("channel.remove_user", Handler{})
    events.Register("user.offboard", Handler{})
    events.Register("user.reactivate", Handler{})
}
