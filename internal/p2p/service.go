package p2p

import (
	"net/http"

	"tic-tac-toe/pkg/websocket"

	"github.com/go-chi/chi/v5"
)

type Service struct {
	m chi.Router
}

func New() *Service {
	s := &Service{
		m: chi.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

func (s *Service) routes() {
	s.m.Get("/", s.handleIndex())
	s.m.Get("/ws", websocket.NewClient(2).ServeHTTP)
}

func (s *Service) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}
}
