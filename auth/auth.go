package auth

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type AuthResult struct {
    UserID string
    Meta   map[string]any
}

type Authenticator interface {
    VerifyRequest(r *http.Request) (*AuthResult, error)
}

type RemoteAuth struct {
    VerifyURL string
    Client    *http.Client
}

func (a *RemoteAuth) VerifyRequest(r *http.Request) (*AuthResult, error) {
    token := r.Header.Get("Authorization")
    if token == "" {
        return nil, fmt.Errorf("missing Authorization header")
    }

    req, _ := http.NewRequest("GET", a.VerifyURL, nil)
    req.Header.Set("Authorization", token)

    resp, err := a.Client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unauthorized")
    }

    var data struct {
        UserID string                 `json:"user_id"`
        Meta   map[string]interface{} `json:"meta"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, err
    }

    return &AuthResult{
        UserID: data.UserID,
        Meta:   data.Meta,
    }, nil
}
