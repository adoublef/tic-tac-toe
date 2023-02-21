package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	gamesHTTP "tic-tac-toe/internal/games/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

var port int

func init() {
	p, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		p = 8080
	}

	flag.IntVar(&port, "port", p, "port to listen on")
}

func run() error {
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	root.Use(cors.AllowAll().Handler)

	root.Mount("/games", gamesHTTP.New())

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
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
