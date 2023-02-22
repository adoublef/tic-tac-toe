package repo

import (
	"context"
	"tic-tac-toe/internal/users"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	CreateProfile(ctx context.Context, args pgx.QueryRewriter) (*users.Profile, error)
}

type Reader interface {
	// GetProfile gets a profile by ID.
	GetProfile(ctx context.Context, args pgx.QueryRewriter) (*users.Profile, error)
}
