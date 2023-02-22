package ccdb

import (
	"context"
	"net/mail"
	"tic-tac-toe/internal/databases/postgres"
	"tic-tac-toe/internal/users"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Profile struct {
	ID       uuid.UUID
	Email    string
	Username string
}

type usersRepo struct {
	rwc postgres.Conn[users.Profile]
}

const getProfile = `SELECT id,email,username FROM "profile" WHERE username = $1`
const getProfileN = `SELECT id,email,username FROM "profile" WHERE username = @username`

type GetProfile struct {
	Username string
}

// GetProfile implements repo.Repo
func (r *usersRepo) GetProfile(ctx context.Context, args GetProfile) (*users.Profile, error) {
	// p := Profile{}
	// err := r.rwc.Conn().QueryRow(ctx, getProfile, args.Username).
	// 	Scan(&p.ID, &p.Email, &p.Username)

	// if err != nil {
	// 	return nil, err
	// }

	// saved := users.Profile{
	// 	ID:       p.ID,
	// 	Email:    &mail.Address{Address: p.Email},
	// 	Username: p.Username,
	// }
	// return &saved, err

	named := pgx.NamedArgs{
		"username": args.Username,
	}

	return r.rwc.QueryRowContext(ctx, func(row pgx.Row, p *users.Profile) error {
		v := Profile{}
		err := row.Scan(&v.ID, &v.Email, &v.Username)
		if err != nil {
			return err
		}

		p.ID = v.ID
		p.Email = &mail.Address{Address: v.Email}
		p.Username = v.Username

		return nil
	}, getProfileN, named)
}

const createProfile = `INSERT INTO "profile" (id,email,username) VALUES ($1,$2,$3)`
const createProfileN = `INSERT INTO "profile" (id,email,username) VALUES (@id,@email,@username)`

type CreateProfile struct {
	Email    *mail.Address
	Username string
}

// CreateProfile implements repo.Repo
func (r *usersRepo) CreateProfile(ctx context.Context, args CreateProfile) (uuid.UUID, error) {
	// id := uuid.New()

	// _, err := r.rwc.Conn().Exec(ctx, createProfile, id, args.Email, args.Username)
	// if err != nil {
	// 	return uuid.Nil, err
	// }

	// return id, nil
	id := uuid.New()

	named := pgx.NamedArgs{
		"id":       id,
		"email":    args.Email.Address,
		"username": args.Username,
	}

	n, err := r.rwc.ExecContext(ctx, createProfileN, named)
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
