# WS-Framework
WS-Framework is a lightweight WebSocket framework for Go. It is designed to make it easy to build real time applications using WebSocket connections. The framework provides room based messaging, direct messages (DMs), event routing, and pluggable authentication. It is modular, simple to extend, and works for both production and development.

## FEATURES

- WebSocket server with manual handshake
- Dynamic room creation
- Room broadcasting, registration, and unregistration
- Event routing system
- Pluggable authentication (RemoteAuth, MockAuth, NoAuth)
- Chat module with optional persistence (database, REST, or none)
- Direct message (DM) room support
- Anti-spoofing protection in event handlers
- Clean separation of transport, events, auth, and modules

## PROJECT STRUCTURE
```
ws-framework/
    transport/
    events/
    auth/
    chat/
    go.mod
    README.md
```

## CORE CONCEPTS

1. Rooms: A room is a broadcast channel. Any client connected to a room receives all messages sent to it.

2. RoomHub:Manages all rooms and creates them when needed.

3. Client: Represents a WebSocket connection. Each client belongs to one room.

4. Events: Messages sent by clients are decoded into events.Event and routed to handlers.

5. Authentication: 
    
    Three modes are supported:
    1. RemoteAuth: verifies tokens via HTTP
    2. MockAuth: predictable identity for development
    3. NoAuth: no authentication, assigns guest IDs

6. Direct Messages:
DM rooms are rooms with a deterministic name:
dm:<userA>:<userB>

The auth service controls who can join which DM.

## INSTALLATION

go get github.com/timkrans/ws-framework



## MINIMAL SERVER EXAMPLE

```
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

## MOCK AUTH EXAMPLE

```
authenticator := &auth.MockAuth{
    UserID: "dev-user",
    Rooms:  []string{"lobby", "dm:dev-user:alice"},
}
```

## REMOTE AUTH EXAMPLE
```
authenticator := &auth.RemoteAuth{
    VerifyURL: "https://auth.example.com/verify",
    Client:    http.DefaultClient,
}
```
### Expected JSON response from auth service:
```
{
  "user_id": "123",
  "meta": {
    "rooms": ["lobby", "dm:123:456"]
  }
}
```

## CHAT MODULE EXAMPLE

### Initialization:
```
chat.Init(chat.ChatPersistenceConfig{
    Mode:    "none",
    RESTURL: "",
})
```

### Client sends a chat message:
```
{
  "type": "chat.message",
  "room": "lobby",
  "user": "123",
  "data": { "text": "Hello world!" }
}
```
### Typing indicator:
```
{
  "type": "chat.typing",
  "room": "lobby",
  "user": "123",
  "data": {}
}
```
## DM ROOM NAMING

```
func DMRoomID(a, b string) string {
    if a < b {
        return "dm:" + a + ":" + b
    }
    return "dm:" + b + ":" + a
}
```
SECURITY MODEL

The framework enforces:

1. Room access control:
   Only rooms listed in AuthResult.Meta["rooms"] may be joined.
2. Anti spoofing:
   Event handlers must verify:
   if evt.User != client.UserID { return }

3. Externalized authentication:
   All real authentication and authorization happens in your auth service.


## FULL CHAT HANDLER WITH ANTI-SPOOFING
```
func (h ChatHandler) Handle(c interface{}, evt events.Event) {
    client := c.(*transport.Client)

    if evt.User != client.UserID {
        return
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

    out, _ := json.Marshal(evt)
    client.Room.Broadcast <- out
}
```

## CLIENT-SIDE JAVASCRIPT EXAMPLE
```
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

