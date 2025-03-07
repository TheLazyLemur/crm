package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"

	"simplecrm/internal/db"
)

func main() {
	dbc, err := sql.Open("sqlite3", "./simplecrm.db")
	if err != nil {
		panic(err)
	}

	querier := db.New()
	querier.GetUser(context.Background(), dbc, "1")

	r := chi.NewRouter()

	mountRoutes(r, dbc, querier)

	http.ListenAndServe(":8080", r)
}

func mountRoutes(r chi.Router, dbc *sql.DB, querier db.Querier) {
	r.Route("/api/v1/query", func(r chi.Router) {
		r.Get("/user/{id}", getUser())
		r.Get("/lead/{id}", getLead())
		r.Get("/contact/{id}", getContact())
		r.Get("/task/{id}", getTask())
	})

	r.Route("/api/v1/user", func(r chi.Router) {
		r.Post("/create", createUser(dbc, querier))
		r.Post("/update/{id}", updateUser())
		r.Post("/command", handleUserCommand())
	})

	r.Route("/api/v1/lead", func(r chi.Router) {
		r.Post("/create", createLead())
		r.Patch("/update/{id}", updateLead())
		r.Post("/command", handleLeadCommand())
	})

	r.Route("/api/v1/contact", func(r chi.Router) {
		r.Post("/create", createContact())
		r.Patch("/update/{id}", updateContact())
		r.Post("/command", handleContactCommand())
	})

	r.Route("/api/v1/task", func(r chi.Router) {
		r.Post("/create", createTask())
		r.Patch("/update/{id}", updateTask())
		r.Post("/command", handleTaskCommand())
	})
}
