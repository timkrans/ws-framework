package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "time"

    "github.com/timkrans/ws-framework/events"
)

var config Config

type Config struct {
    EnableWebhooks bool
    WebhookURL     string
}

type Handler struct{}

func (h Handler) Handle(c interface{}, evt events.Event) {
    if !config.EnableWebhooks || config.WebhookURL == "" {
        return
    }
    body, err := json.Marshal(evt)
    if err != nil {
        return
    }
    req, err := http.NewRequest("POST", config.WebhookURL, bytes.NewBuffer(body))
    if err != nil {
        return
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{
        Timeout: 5 * time.Second,
    }

    //fire and forget
    _, _ = client.Do(req)
}

func Init(cfg Config) {
    config = cfg
    events.Register("user.update", Handler{})
}
