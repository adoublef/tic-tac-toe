package p2p

import (
	"net/http"
	"strconv"
	"sync"

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
	var store sync.Map

	s.m.Get("/", s.handleIndex())
	s.m.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		// get search params and make sure its a number
		// if not, return an error
		q := r.URL.Query().Get("id")
		if q == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(q)
		if err != nil {
			http.Error(w, "id must be a number", http.StatusBadRequest)
			return
		}

		client, _ := store.LoadOrStore(id, websocket.NewClient(2))

		client.(*websocket.Client).ServeHTTP(w, r)
	})
}

func (s *Service) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}
}
