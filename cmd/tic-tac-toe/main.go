package main

import (
	"fmt"
	"log"
	"net/http"
	p2pHTTP "tic-tac-toe/internal/p2p"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func run() error {
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	root.Use(cors.AllowAll().Handler)

	root.Mount("/game", p2pHTTP.New())

	s := &http.Server{
		Addr:    ":8080",
		Handler: root,
	}

	fmt.Printf("Listening on %s\n", s.Addr)
	return s.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("main: %v", err)
	}
}
