package ws

import (
    "net/http"
    "github.com/timkrans/ws-framework/modules/chat"
    "github.com/timkrans/ws-framework/modules/call"
    "github.com/timkrans/ws-framework/transport"
)

type Server struct {
    Hub    *transport.RoomHub
    Config ServerConfig
}

type ServerConfig struct {
    ChatPersistence ChatPersistenceConfig
}

type ChatPersistenceConfig struct {
    Mode    string 
    RESTURL string
}


//for chat modules
func NewServer(cfg ServerConfig) *Server {
    chat.Init(chat.ChatPersistenceConfig{
        Mode:    cfg.ChatPersistence.Mode,
        RESTURL: cfg.ChatPersistence.RESTURL,
    })

    return &Server{
        Hub:    transport.NewRoomHub(),
        Config: cfg,
    }
}

func NewCallServer() *Server {
    //only load the call module
    call.Init()

    return &Server{
        Hub: transport.NewRoomHub(),
        Config: ServerConfig{
            ChatPersistence: ChatPersistenceConfig{
                Mode: "none",
            },
        },
    }
}


func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
    transport.HandleWebSocket(s.Hub, w, r)
}