package ccdb_test

import (
	"context"
	"net/mail"
	"testing"
	"tic-tac-toe/internal/users/repo/ccdb"
	"tic-tac-toe/pkg/containers"

	"github.com/google/uuid"
	"github.com/hyphengolang/prelude/testing/is"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestRepository(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	c, err := containers.NewCockroachDBContainer(ctx)
	is.NoErr(err) // test-container should start

	defer c.Terminate(ctx)

	// #0 -- initialize the repo
	rwc, err := pgxpool.New(ctx, c.URI+"/testing")
	is.NoErr(err) // pgxpool should connect

	repo := ccdb.New(rwc)
	// TODO -- run migration
	{
		const migration = `
		create database testing;
		create table testing.profile (
			id uuid primary key not null,
			email string not null,
			username string not null
		);`
		_, err := rwc.Exec(ctx, migration)
		is.NoErr(err) // migration should run
	}

	// #1 -- create a new profile with a unique email
	{
		e, _ := mail.ParseAddress("adoublef@mail.com")

		args := ccdb.CreateProfile{
			Email:    e,
			Username: "adoublef",
		}

		created, err := repo.CreateProfile(ctx, args)
		is.NoErr(err)                // create should run
		is.True(created != uuid.Nil) // created should not be nil
	}

	// #2 -- get the profile by username
	{
		args := ccdb.GetProfile{
			Username: "adoublef",
		}

		found, err := repo.GetProfile(ctx, args)
		is.NoErr(err)         // get should run
		is.True(found != nil) // found should not be nil
		is.Equal(found.Email.Address, "adoublef@mail.com")
	}
}
