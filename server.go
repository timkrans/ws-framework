package ws

import (
    "net/http"

    "github.com/timkrans/ws-framework/auth"
    "github.com/timkrans/ws-framework/modules/chat"
    "github.com/timkrans/ws-framework/modules/call"
    "github.com/timkrans/ws-framework/transport"
)

type Server struct {
    Hub    *transport.RoomHub
    Auth   auth.Authenticator
    Config ServerConfig
}

type ServerConfig struct {
    ChatPersistence ChatPersistenceConfig
    Authenticator   auth.Authenticator 
}

type ChatPersistenceConfig struct {
    Mode    string
    RESTURL string
}

func NewServer(cfg ServerConfig) *Server {
    chat.Init(chat.ChatPersistenceConfig{
        Mode:    cfg.ChatPersistence.Mode,
        RESTURL: cfg.ChatPersistence.RESTURL,
    })

    authenticator := cfg.Authenticator
    if authenticator == nil {
        authenticator = &auth.NoAuth{} 
    }

    return &Server{
        Hub:    transport.NewRoomHub(),
        Auth:   authenticator,
        Config: cfg,
    }
}

func NewCallServer() *Server {
    call.Init()

    return &Server{
        Hub:  transport.NewRoomHub(),
        Auth: &auth.NoAuth{},
        Config: ServerConfig{
            ChatPersistence: ChatPersistenceConfig{
                Mode: "none",
            },
        },
    }
}

func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
    transport.HandleWebSocket(s.Hub, s.Auth, w, r)
}
