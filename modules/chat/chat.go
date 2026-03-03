package chat

import (
    "bytes"
    "encoding/json"
    "net/http"
    "time"

    "gorm.io/gorm"
    "github.com/timkrans/ws-framework/events"
    "github.com/timkrans/ws-framework/transport"
    "github.com/timkrans/ws-framework/auth"
)

var DB *gorm.DB
var Config ChatPersistenceConfig

type ChatPersistenceConfig struct {
    Mode    string
    RESTURL string
}

type IncomingChat struct {
    Text string `json:"text"`
}

type ChatHandler struct{}

func (h ChatHandler) Handle(c interface{}, evt events.Event) {
    client := c.(*transport.Client)

    switch client.Source.(type) {
    case *auth.RemoteAuth:
        if evt.User != client.UserID {
            return
        }
    default:
        evt.User = client.UserID
    }

    var in IncomingChat
    json.Unmarshal(evt.Data, &in)

    if evt.Type == "chat.typing" {
        out, _ := json.Marshal(evt)
        client.Room.Broadcast <- out
        return
    }

    msg := Message{
        Room:      evt.Room,
        User:      evt.User,
        Text:      in.Text,
        CreatedAt: time.Now().Unix(),
    }

    switch Config.Mode {
    case "db":
        DB.Create(&msg)
    case "rest":
        body, _ := json.Marshal(msg)
        http.Post(Config.RESTURL, "application/json", bytes.NewBuffer(body))
    case "none":
    }

    out, _ := json.Marshal(evt)
    client.Room.Broadcast <- out
}

func Init(cfg ChatPersistenceConfig) {
    Config = cfg
    events.Register("chat.message", ChatHandler{})
    events.Register("chat.typing", ChatHandler{})
}
