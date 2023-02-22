package repo

import (
	"context"
	"net/mail"
	"tic-tac-toe/internal/users"

	"github.com/google/uuid"
)

type Profile struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}

type Repo interface {
	Reader
	Writer
}

type Writer interface {
	// CreateProfile creates a new profile.
	CreateProfile(ctx context.Context, args *CreateProfileArgs) (uuid.UUID, error)
}

type Reader interface {
	// GetProfile gets a profile by ID.
	GetProfile(ctx context.Context, args *GetProfileArgs) (*users.Profile, error)
}

// CreateProfileArgs is the arguments for creating a new profile.
type CreateProfileArgs struct {
	Email    *mail.Address
	Username string
}

// GetProfileArgs is the arguments for getting a profile.
type GetProfileArgs struct {
	Username string
}
