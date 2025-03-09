package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"simplecrm/internal/db"
)

func JSONDecoderMiddleware[Req Validatable](
	handler handlerFunc[Req],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isJsonRequest := r.Header.Get("Content-Type") == "application/json"
		if !isJsonRequest {
			http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
			return
		}

		var params Req

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		handler(w, r, params)
	}
}

func MountRoutes(r chi.Router, dbc *sqlx.DB, querier db.Querier) {
	r.Route("/api/v1/query", func(r chi.Router) {
		r.Get("/user/{id}", GetUser())
		r.Get("/lead/{id}", GetLead())
		r.Get("/contact/{id}", GetContact())
		r.Get("/task/{id}", GetTask())
	})

	r.Route("/api/v1/user", func(r chi.Router) {
		r.Post("/create", JSONDecoderMiddleware(CreateUser(dbc, querier)))
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
