package ccdb

import (
	"context"
	"net/mail"
	"tic-tac-toe/internal/databases/postgres"
	"tic-tac-toe/internal/users"
	"tic-tac-toe/internal/users/repo"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type usersRepo struct {
	rwc postgres.Conn[users.Profile]
}

const getProfile = `SELECT id,email,username FROM "profile" WHERE username = @username`

// GetProfile implements repo.Repo
func (r *usersRepo) GetProfile(ctx context.Context, args *repo.GetProfileArgs) (*users.Profile, error) {
	named := pgx.NamedArgs{
		"username": args.Username,
	}

	return r.rwc.QueryRowContext(ctx, func(row pgx.Row, p *users.Profile) error {
		v := repo.Profile{}
		err := row.Scan(&v.ID, &v.Email, &v.Username)
		if err != nil {
			return err
		}

		p.ID = v.ID
		p.Email = &mail.Address{Address: v.Email}
		p.Username = v.Username

		return nil
	}, getProfile, named)
}

const createProfile = `INSERT INTO "profile" (id,email,username) VALUES (@id,@email,@username)`

// CreateProfile implements repo.Repo
func (r *usersRepo) CreateProfile(ctx context.Context, args *repo.CreateProfileArgs) (uuid.UUID, error) {
	id := uuid.New()

	named := pgx.NamedArgs{
		"id":       id,
		"email":    args.Email.Address,
		"username": args.Username,
	}

	n, err := r.rwc.ExecContext(ctx, createProfile, named)
	if err != nil {
		return uuid.Nil, err
	}

	if n != 1 {
		return uuid.Nil, pgx.ErrNoRows
	}

	return id, nil
}

func New(rwc *pgxpool.Pool) *usersRepo {
	r := &usersRepo{
		rwc: postgres.NewConn[users.Profile](rwc),
	}

	return r
}
