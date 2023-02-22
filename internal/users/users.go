package users

import (
	"net/mail"

	"github.com/google/uuid"
)

type Profile struct {
	ID       uuid.UUID     `json:"id"`
	Email    *mail.Address `json:"email"`
	Username string        `json:"username"`
}
