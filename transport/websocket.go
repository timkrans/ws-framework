package transport

import (
    "bufio"
    "crypto/sha1"
    "encoding/base64"
    "net/http"
    "strings"
    "github.com/timkrans/ws-framework/auth"
)

const wsGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func HandleWebSocket(hub *RoomHub, authenticator auth.Authenticator, w http.ResponseWriter, r *http.Request) {
    info, err := authenticator.VerifyRequest(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    roomName := r.URL.Query().Get("room")
    if roomName == "" {
        roomName = "lobby"
    }
    room := hub.GetRoom(roomName)

    allowedRooms, _ := info.Meta["rooms"].([]any)
    if !roomAllowed(roomName, allowedRooms) {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    if !isWebSocketUpgrade(r) {
        http.Error(w, "Not a WebSocket upgrade", 400)
        return
    }

    hj, _ := w.(http.Hijacker)
    conn, _, _ := hj.Hijack()

    secKey := r.Header.Get("Sec-WebSocket-Key")
    accept := computeAcceptKey(secKey)

    resp := "HTTP/1.1 101 Switching Protocols\r\n" +
        "Upgrade: websocket\r\n" +
        "Connection: Upgrade\r\n" +
        "Sec-WebSocket-Accept: " + accept + "\r\n\r\n"

    conn.Write([]byte(resp))

    client := &Client{
        Conn:   conn,
        Reader: bufio.NewReader(conn),
        Send:   make(chan []byte, 256),
        Room:   room,
    }

    room.Register <- client

    go client.WriteLoop()
    client.ReadLoop()
}

func roomAllowed(room string, allowed []any) bool {
    for _, r := range allowed {
        if rs, ok := r.(string); ok {
            if rs == "*" || rs == room {
                return true
            }
        }
    }
    return false
}

func isWebSocketUpgrade(r *http.Request) bool {
    return strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade") &&
        strings.ToLower(r.Header.Get("Upgrade")) == "websocket"
}

func computeAcceptKey(key string) string {
    h := sha1.New()
    h.Write([]byte(key + wsGUID))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
