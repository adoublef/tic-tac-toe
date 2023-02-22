package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	usersHTTP "tic-tac-toe/internal/users/http"
	"tic-tac-toe/internal/users/repo"
	"tic-tac-toe/internal/users/repo/ccdb"
	"tic-tac-toe/pkg/containers"

	"github.com/hyphengolang/prelude/testing/is"
	"github.com/jackc/pgx/v5/pgxpool"
)

const applicationJSON = "application/json"

func TestUserService(t *testing.T) {
	is := is.New(t)

	// #0 -- initialize the repo
	testRepo, repoClose := testRepo(is, context.Background())
	defer repoClose()

	// #1 -- initialize the service
	srv := httptest.NewServer(testService(testRepo))
	defer srv.Close()

	// #2 -- create a new profile with a unique email
	type testcaseA struct {
		name       string
		payload    string
		statusCode int
	}

	tta := []testcaseA{
		{
			name:       "create a new profile with a unique email",
			statusCode: http.StatusCreated,
			payload: `
			{
				"email": "adoublef@mail.com",
				"username": "adoublef"
			}`,
		},
	}

	for _, ta := range tta {
		t.Run(ta.name, func(t *testing.T) {
			res, err := srv.Client().Post(srv.URL+"/", applicationJSON, strings.NewReader(ta.payload))
			is.NoErr(err)                           // POST should succeed
			is.Equal(res.StatusCode, ta.statusCode) // POST should return 201
		})
	}

	// #3 -- search for a profile by email
	type testcaseB struct {
		name       string
		param      string
		statusCode int
	}

	ttb := []testcaseB{
		{
			name:       "user found",
			param:      "adoublef",
			statusCode: http.StatusOK,
		},
		{
			name:       "user not found",
			param:      "aboudlef",
			statusCode: http.StatusNotFound,
		},
	}

	for _, tb := range ttb {
		t.Run(fmt.Sprintf("name=%s param=%s", tb.name, tb.param), func(t *testing.T) {
			res, err := srv.Client().Get(srv.URL + "/" + tb.param)
			is.NoErr(err)                           // GET should succeed
			is.Equal(res.StatusCode, tb.statusCode) // GET should return 200
		})
	}
}

func testRepo(is *is.I, ctx context.Context) (repo repo.Repo, close func()) {
	c, err := containers.NewCockroachDBContainer(ctx)
	is.NoErr(err) // test-container should start

	rwc, err := pgxpool.New(ctx, c.URI+"/testing")
	is.NoErr(err) // pgxpool should connect

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

	return ccdb.New(rwc), func() { c.Terminate(ctx) }
}

func testService(repo repo.Repo) http.Handler {
	return usersHTTP.New(usersHTTP.WithRepo(repo))
}
