package auth

import (
    "fmt"
    "net/http"
    "sync/atomic"
)

var noAuthCounter uint64

type NoAuth struct{}

func (n *NoAuth) VerifyRequest(r *http.Request) (*AuthResult, error) {
    id := atomic.AddUint64(&noAuthCounter, 1)

    return &AuthResult{
        UserID: fmt.Sprintf("guest-%d", id),
        Meta: map[string]any{
            "rooms": []string{"*"}, //allow all rooms
        },
    }, nil
}
