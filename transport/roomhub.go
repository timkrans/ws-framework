package transport

import "sync"

type RoomHub struct {
    mu    sync.Mutex
    Rooms map[string]*Room
}

func NewRoomHub() *RoomHub {
    return &RoomHub{
        Rooms: make(map[string]*Room),
    }
}

func (h *RoomHub) GetRoom(name string) *Room {
    h.mu.Lock()
    defer h.mu.Unlock()

    room, ok := h.Rooms[name]
    if !ok {
        room = NewRoom(name)
        h.Rooms[name] = room
        go room.Run()
    }
    return room
}
