package ccdb_test

import (
	"context"
	"fmt"
	"net/mail"
	"testing"
	"tic-tac-toe/internal/users/repo"
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

	testRepo := ccdb.New(rwc)
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
	type testcaseA struct {
		email    string
		username string
	}

	tta := []testcaseA{
		{
			email:    "adoublef@mail.com",
			username: "adoublef",
		},
		// TODO -- add cases where this should fail
	}

	for _, ta := range tta {
		t.Run(fmt.Sprintf("email=%s,username=%s", ta.email, ta.username), func(t *testing.T) {
			e, _ := mail.ParseAddress(ta.email)

			args := repo.CreateProfileArgs{
				Email:    e,
				Username: ta.username,
			}

			created, err := testRepo.CreateProfile(ctx, &args)
			is.NoErr(err)                // create should run
			is.True(created != uuid.Nil) // created should not be nil
		})
	}

	// #2 -- get the profile by username
	type testcaseB struct {
		username string
	}

	ttb := []testcaseB{
		{
			username: "adoublef",
		},
		// TODO -- add cases where this should fail
	}

	for _, tb := range ttb {
		t.Run(fmt.Sprintf("username=%s", tb.username), func(t *testing.T) {
			args := repo.GetProfileArgs{
				Username: "adoublef",
			}

			found, err := testRepo.GetProfile(ctx, &args)
			is.NoErr(err)         // get should run
			is.True(found != nil) // found should not be nil
			is.Equal(found.Username, tb.username)
		})
	}
}
