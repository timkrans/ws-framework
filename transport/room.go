package transport

import "sync"

type Room struct {
    Name       string
    Clients    map[*Client]bool
    Register   chan *Client
    Unregister chan *Client
    Broadcast  chan []byte
    mu         sync.Mutex
}

func NewRoom(name string) *Room {
    return &Room{
        Name:       name,
        Clients:    make(map[*Client]bool),
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
        Broadcast:  make(chan []byte),
    }
}

func (r *Room) Run() {
    for {
        select {
        case c := <-r.Register:
            r.mu.Lock()
            r.Clients[c] = true
            r.mu.Unlock()

        case c := <-r.Unregister:
            r.mu.Lock()
            if _, ok := r.Clients[c]; ok {
                delete(r.Clients, c)
                close(c.Send)
            }
            r.mu.Unlock()

        case msg := <-r.Broadcast:
            r.mu.Lock()
            for c := range r.Clients {
                select {
                case c.Send <- msg:
                default:
                    delete(r.Clients, c)
                    close(c.Send)
                }
            }
            r.mu.Unlock()
        }
    }
}
