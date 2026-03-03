package call

import (
    "encoding/json"
    "github.com/timkrans/ws-framework/auth"
    "github.com/timkrans/ws-framework/events"
    "github.com/timkrans/ws-framework/transport"
)

type CallHandler struct{}

func (h CallHandler) Handle(c interface{}, evt events.Event) {
    client := c.(*transport.Client)

    switch client.Source.(type) {
    case *auth.RemoteAuth:
        if evt.User != client.UserID {
            return
        }
    default:
        evt.User = client.UserID
    }

    data, _ := json.Marshal(evt)
    client.Room.Broadcast <- data
}

func Init() {
    events.Register("call.offer", CallHandler{})
    events.Register("call.answer", CallHandler{})
    events.Register("call.ice", CallHandler{})
    events.Register("call.end", CallHandler{})
}
