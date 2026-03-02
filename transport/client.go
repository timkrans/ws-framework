package transport

import (
    "bufio"
    "encoding/json"
    "net"
    "github.com/timkrans/ws-framework/events"
    "github.com/timkrans/ws-framework/auth"
)

type Client struct {
    Conn    net.Conn
    Reader  *bufio.Reader
    Send    chan []byte
    Room    *Room
    UserID  string
    Auth    *auth.AuthResult
}

func (c *Client) ReadLoop() {
    defer func() {
        c.Room.Unregister <- c
        c.Conn.Close()
    }()

    for {
        op, payload, err := ReadFrame(c.Reader)
        if err != nil {
            return
        }

        if op != 0x1 {
            continue
        }

        var evt events.Event
        if err := json.Unmarshal(payload, &evt); err != nil {
            continue
        }

        if handler := events.GetHandler(evt.Type); handler != nil {
            handler.Handle(c, evt)
        }
    }
}

func (c *Client) WriteLoop() {
    defer c.Conn.Close()

    for msg := range c.Send {
        WriteFrame(c.Conn, 0x1, msg)
    }
}
