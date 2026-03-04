# WS-Framework

WS-Framework is a lightweight, modular WebSocket framework for Go. It is designed for building real‑time applications with room‑based messaging, direct messages (DMs), event routing, and pluggable authentication. The framework emphasizes simplicity, extensibility, and clean separation of concerns.

## FEATURES

- WebSocket server with manual handshake
- Dynamic room creation and automatic cleanup
- Room broadcasting, registration, and unregistration
- Event routing system with anti‑spoofing protection
- Pluggable authentication (RemoteAuth, MockAuth, NoAuth)
- Chat module with optional persistence (database, REST, or none)
- Direct message (DM) room support
- Presence module (online, typing, status)
- File sharing module
- Notifications module
- Admin module (remove user, offboard, reactivate)
- Reactions module (emoji reactions)
- Call module (WebRTC signaling)
- Clean separation of transport, events, auth, and modules

---

## PROJECT STRUCTURE

```
ws-framework/
    transport/      # WebSocket transport layer
    events/         # Event routing system
    auth/           # Authentication providers
    modules/
        chat/
        call/
        presence/
        files/
        notify/
        admin/
        reactions/
    server.go           # Server assembly
    go.mod
    README.md
```

---

## CORE CONCEPTS

### 1. Rooms
A room is a broadcast channel. Any client connected to a room receives all messages sent to it.

### 2. RoomHub
Manages all rooms, creates them on demand, and handles registration/unregistration of clients.

### 3. Client
Represents a WebSocket connection. Each client belongs to exactly one room.

### 4. Events
Incoming WebSocket messages are decoded into `events.Event` and routed to the correct module handler.

### 5. Authentication
Three built‑in modes:

- **RemoteAuth** — verifies tokens via HTTP
- **MockAuth** — predictable identity for development
- **NoAuth** — no authentication, assigns guest IDs

### 6. Direct Messages (DMs)
DM rooms use deterministic naming:

```
dm:<userA>:<userB>
```

The auth service controls who may join which DM.

### 7. Anti‑spoofing
Every module must enforce:

```
if evt.User != client.UserID { return }
```

This prevents clients from impersonating other users.

---

## INSTALLATION

```
go get github.com/timkrans/ws-framework
```

---

## MINIMAL SERVER EXAMPLE

```go
package main

import (
    "net/http"
    "github.com/timkrans/ws-framework/transport"
    "github.com/timkrans/ws-framework/auth"
)

func main() {
    hub := transport.NewRoomHub()
    authenticator := &auth.NoAuth{}

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        transport.HandleWebSocket(hub, authenticator, w, r)
    })

    http.ListenAndServe(":8080", nil)
}
```

---

## AUTHENTICATION EXAMPLES

### Mock Auth

```go
authenticator := &auth.MockAuth{
    UserID: "dev-user",
    Rooms:  []string{"lobby", "dm:dev-user:alice"},
}
```

### Remote Auth

```go
authenticator := &auth.RemoteAuth{
    VerifyURL: "https://auth.example.com/verify",
    Client:    http.DefaultClient,
}
```

Expected JSON response:

```json
{
  "user_id": "123",
  "meta": {
    "rooms": ["lobby", "dm:123:456"]
  }
}
```

---

## CHAT MODULE

### Initialization

```go
chat.Init(chat.ChatPersistenceConfig{
    Mode:    "none",
    RESTURL: "",
})
```

### Sending a chat message

```json
{
  "type": "chat.message",
  "room": "lobby",
  "user": "123",
  "data": { "text": "Hello world!" }
}
```

---

## PRESENCE MODULE

Supports:

- `presence.update` — online/offline
- `presence.typing` — typing indicator
- `presence.away`
- `presence.idle`
- `presence.dnd`
- `presence.back`
- `presence.mobile`
- `presence.status_text`

---

## FILE SHARING MODULE

Clients can broadcast file metadata:

```json
{
  "type": "file.share",
  "room": "lobby",
  "user": "123",
  "data": {
    "name": "image.png",
    "size": 2048,
    "type": "image/png",
    "url": "blob:..."
  }
}
```

---

## NOTIFY MODULE

Simple push notifications:

```json
{
  "type": "notify.send",
  "room": "lobby",
  "user": "123",
  "data": { "message": "Build completed!" }
}
```

---

## ADMIN MODULE

Slack‑style administrative actions:

- `channel.remove_user`
- `user.offboard`
- `user.reactivate`

Example:

```json
{
  "type": "channel.remove_user",
  "room": "lobby",
  "user": "admin",
  "data": { "target": "guest-2" }
}
```

---

## REACTIONS MODULE

Emoji reactions:

- `reaction.add`
- `reaction.remove`

Example:

```json
{
  "type": "reaction.add",
  "room": "lobby",
  "user": "123",
  "data": { "messageId": "msg-1", "emoji": "👍" }
}
```

---

## CALL MODULE (WebRTC Signaling)

Supports:

- `call.offer`
- `call.answer`
- `call.ice`
- `call.end`

Used for peer‑to‑peer audio/video calls.

---

## DM ROOM NAMING

```go
func DMRoomID(a, b string) string {
    if a < b {
        return "dm:" + a + ":" + b
    }
    return "dm:" + b + ":" + a
}
```

---

## SECURITY MODEL

1. **Room access control**  
   Only rooms listed in `AuthResult.Meta["rooms"]` may be joined.

2. **Anti‑spoofing**  
   Handlers must enforce identity:

   ```go
   if evt.User != client.UserID { return }
   ```

3. **Externalized authentication**  
   All real auth happens in your auth service.

---

## CLIENT-SIDE JAVASCRIPT EXAMPLE

```js
const ws = new WebSocket("ws://localhost:8080/ws?room=lobby");

ws.onopen = () => {
    ws.send(JSON.stringify({
        type: "chat.message",
        room: "lobby",
        user: "guest-1",
        data: { text: "Hello!" }
    }));
};

ws.onmessage = (msg) => {
    console.log("Received:", msg.data);
};
```

---

## FUTURE

- Expand the base module system to support:
  - Threads (threaded replies)
  - User profiles
  - Channel management
  - Polls
  - Pins/bookmarks
  - Tasks
  - Search indexing
- Add optional persistence for presence, reactions, and admin actions
- Add cluster support for multi‑node deployments