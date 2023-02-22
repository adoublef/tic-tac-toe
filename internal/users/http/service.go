package http

import (
	"net/http"
	"net/mail"

	service "tic-tac-toe/internal/http"
	"tic-tac-toe/internal/users"
	"tic-tac-toe/internal/users/repo"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Option func(*Service)

func WithRepo(r repo.Repo) Option {
	return func(s *Service) {
		s.r = r
	}
}

type Service struct {
	// m is the multiplexer for the service.
	m service.Service
	// r is the repo for the service.
	r repo.Repo
}

func New(opts ...Option) *Service {
	s := &Service{
		m: service.New(),
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.r == nil {
		panic("repo is nil")
	}

	s.routes()
	return s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

func (s *Service) routes() {
	s.m.Post("/", s.handleCreateProfile())
	s.m.Get("/{username}", s.handleGetProfile())
}

func (s *Service) handleCreateProfile() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	parse := func(w http.ResponseWriter, r *http.Request) (*repo.CreateProfileArgs, error) {
		var q request
		if err := s.m.Decode(w, r, &q); err != nil {
			return nil, err
		}

		email, err := mail.ParseAddress(q.Email)
		if err != nil {
			return nil, err
		}

		return &repo.CreateProfileArgs{
			Email:    email,
			Username: q.Username,
		}, nil
	}

	type response struct {
		ID       uuid.UUID `json:"id"`
		Location string    `json:"location"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		args, err := parse(w, r)
		if err != nil {
			s.m.Respond(w, r, err, http.StatusBadRequest)
			return
		}

		// TODO -- this may be an anti-pattern but I find it easier
		created, err := s.r.CreateProfile(r.Context(), args)
		if err != nil {
			s.m.Respond(w, r, err, http.StatusInternalServerError)
			return
		}

		p := response{
			ID: created,
			// FIXME -- should be pointing to client location
			Location: "/v1/users/" + created.String(),
		}

		// TODO -- set location header
		s.m.Respond(w, r, p, http.StatusCreated)
	}
}

func (s *Service) handleGetProfile() http.HandlerFunc {
	parse := func(w http.ResponseWriter, r *http.Request) (*repo.GetProfileArgs, error) {
		username := chi.URLParam(r, "username")

		return &repo.GetProfileArgs{
			Username: username,
		}, nil
	}

	type response struct {
		Profile  *users.Profile `json:"profile"`
		Location string         `json:"location"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		args, err := parse(w, r)
		if err != nil {
			s.m.Respond(w, r, err, http.StatusBadRequest)
			return
		}

		found, err := s.r.GetProfile(r.Context(), args)
		if err != nil {
			s.m.Respond(w, r, err, http.StatusNotFound)
			return
		}

		p := response{
			Profile: found,
			// FIXME -- should be pointing to client location
			Location: "/v1/users/" + found.ID.String(),
		}

		s.m.Respond(w, r, p, http.StatusOK)
	}
}
