package http

import (
	"errors"
	"fmt"
	"net/http"

	"tic-tac-toe/internal/games"

	service "tic-tac-toe/internal/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Service struct {
	// m is the multiplexer for the service.
	m service.Service
	// br is the websocket broker.
	br games.Broker
}

func New() *Service {
	s := &Service{
		m:  service.New(),
		br: games.NewBroker(),
	}
	s.routes()
	return s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

func (s *Service) routes() {
	s.m.Post("/", s.handleCreateGame())
	s.m.Get("/", s.handleListGames())
	s.m.Get("/{id}", s.handleFindGame())
	s.m.Put("/{id}", s.handleUpdateGameInfo())
	s.m.Delete("/{id}", s.handleCloseGame())
	s.m.Get("/ws", s.handleP2P())

}

func (s *Service) handleFindGame() http.HandlerFunc {
	parseID := func(r *http.Request) (id uuid.UUID, err error) {
		id, err = uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			return id, errors.New("id must be a valid uuid")
		}

		return id, nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r)
		if err != nil {
			http.Error(w, "id must be an UUID", http.StatusBadRequest)
			return
		}

		g, loaded := s.br.Load(id)
		if !loaded {
			http.Error(w, "game not found", http.StatusNotFound)
			return
		}

		fmt.Printf("game: %v", g)

		s.m.Respond(w, r, g, http.StatusOK)
	}
}

func (s *Service) handleCreateGame() http.HandlerFunc {
	type P struct {
		ID string `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// create a new game
		g := games.New(2)
		// add the game to the broker
		s.br.Store(g.ID, g)
		// return the game info
		s.m.Respond(w, r, P{ID: g.ID.String()}, http.StatusCreated)
	}
}

func (s *Service) handleCloseGame() http.HandlerFunc {
	parseID := func(r *http.Request) (id uuid.UUID, err error) {
		id, err = uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			return id, errors.New("id must be a valid uuid")
		}

		return id, nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r)
		if err != nil {
			http.Error(w, "id must be an UUID", http.StatusBadRequest)
			return
		}

		g, loaded := s.br.LoadAndDelete(id)
		if !loaded {
			http.Error(w, "game not found", http.StatusNotFound)
			return
		}
		// should notify players that the client is closed
		defer g.Close()
		s.m.Respond(w, r, nil, http.StatusNoContent)
	}
}

func (s *Service) handleP2P() http.HandlerFunc {
	parseID := func(r *http.Request) (id uuid.UUID, err error) {
		id, err = uuid.Parse(r.URL.Query().Get("id"))
		if err != nil {
			return id, errors.New("id must be a valid uuid")
		}

		return id, nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r)
		if err != nil {
			http.Error(w, "id must be an UUID", http.StatusBadRequest)
			return
		}

		loaded, found := s.br.Load(id)
		if !found {
			http.Error(w, "game not found", http.StatusNotFound)
			return
		}

		loaded.ServeHTTP(w, r)
	}
}

func (s *Service) handleListGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.m.Respond(w, r, "list games", http.StatusOK)
	}
}

func (s *Service) handleUpdateGameInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.m.Respond(w, r, "update game info", http.StatusOK)
	}
}
