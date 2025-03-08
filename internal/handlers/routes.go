package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"simplecrm/internal/db"
)

func MountRoutes(r chi.Router, dbc *sqlx.DB, querier db.Querier) {
	r.Route("/api/v1/query", func(r chi.Router) {
		r.Get("/user/{id}", GetUser())
		r.Get("/lead/{id}", GetLead())
		r.Get("/contact/{id}", GetContact())
		r.Get("/task/{id}", GetTask())
	})

	r.Route("/api/v1/user", func(r chi.Router) {
		r.Post("/create", CreateUser(dbc, querier))
		r.Post("/update/{id}", UpdateUser())
		r.Post("/command", HandleUserCommand())
	})

	r.Route("/api/v1/lead", func(r chi.Router) {
		r.Post("/create", CreateLead())
		r.Patch("/update/{id}", UpdateLead())
		r.Post("/command", HandleLeadCommand())
	})

	r.Route("/api/v1/contact", func(r chi.Router) {
		r.Post("/create", CreateContact())
		r.Patch("/update/{id}", UpdateContact())
		r.Post("/command", HandleContactCommand())
	})

	r.Route("/api/v1/task", func(r chi.Router) {
		r.Post("/create", CreateTask())
		r.Patch("/update/{id}", UpdateTask())
		r.Post("/command", HandleTaskCommand())
	})
}
