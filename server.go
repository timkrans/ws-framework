package ws

import (
    "net/http"

    "github.com/timkrans/ws-framework/auth"
    "github.com/timkrans/ws-framework/modules/chat"
    "github.com/timkrans/ws-framework/modules/call"
    "github.com/timkrans/ws-framework/modules/presence"
    "github.com/timkrans/ws-framework/modules/files"
    "github.com/timkrans/ws-framework/modules/notify"
    "github.com/timkrans/ws-framework/modules/admin"
    "github.com/timkrans/ws-framework/transport"
)

type Server struct {
    Hub    *transport.RoomHub
    Auth   auth.Authenticator
    Config ServerConfig
}

type ServerConfig struct {
    Authenticator auth.Authenticator

    EnableChat      bool
    ChatPersistence ChatPersistenceConfig

    EnableCall bool

    EnablePresence bool
    Presence       PresenceConfig

    EnableFiles bool
    Files       FilesConfig

    EnableNotify bool
    Notify       NotifyConfig

    EnableAdmin bool
    Admin       AdminConfig

}

type ChatPersistenceConfig struct {
    Mode    string
    RESTURL string
}

type PresenceConfig struct {
    BroadcastJoinLeave bool
    BroadcastTyping    bool
}

type FilesConfig struct {
    StorageBaseURL string
}

type NotifyConfig struct {
    Persist bool
}

type AdminConfig struct { 
    AllowRemoveUser bool 
    AllowOffboard bool 
    AllowReactivate bool 
}

func NewServer(cfg ServerConfig) *Server {
    if cfg.EnableChat {
        chat.Init(chat.ChatPersistenceConfig{
            Mode:    cfg.ChatPersistence.Mode,
            RESTURL: cfg.ChatPersistence.RESTURL,
        })
    }

    if cfg.EnableCall {
        call.Init()
    }

    if cfg.EnablePresence {
        presence.Init(presence.Config{
            BroadcastJoinLeave: cfg.Presence.BroadcastJoinLeave,
            BroadcastTyping:    cfg.Presence.BroadcastTyping,
        })
    }

    if cfg.EnableFiles {
        files.Init(files.Config{
            StorageBaseURL: cfg.Files.StorageBaseURL,
        })
    }

    if cfg.EnableNotify {
        notify.Init(notify.Config{
            Persist: cfg.Notify.Persist,
        })
    }

    if cfg.EnableAdmin {
        admin.Init(admin.Config{
            AllowRemoveUser: cfg.Admin.AllowRemoveUser, 
            AllowOffboard: cfg.Admin.AllowOffboard, 
            AllowReactivate: cfg.Admin.AllowReactivate,
        })
    }

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

func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
    transport.HandleWebSocket(s.Hub, s.Auth, w, r)
}
