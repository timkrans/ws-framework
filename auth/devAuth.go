package auth

type DevAuth struct {
    UserID string
}

func (d *DevAuth) VerifyRequest(r *http.Request) (*AuthResult, error) {
    return &AuthResult{
        UserID: d.UserID,
        Meta: map[string]any{
            "rooms": []string{"*"}, //allow all rooms
        },
    }, nil
}
