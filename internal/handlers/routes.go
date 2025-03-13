package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"simplecrm/internal/db"
	"simplecrm/internal/pubsub"
)

type httpError struct {
	Message    string
	StatusCode int
}

type httpResponse[T any] struct {
	Data       T
	StatusCode int
}

type (
	handlerFunc[T Validatable, Resp any] func(w http.ResponseWriter, r *http.Request, params T) (*httpResponse[Resp], *httpError)
	getHandlerFunc[Resp any]             func(w http.ResponseWriter, r *http.Request) (*httpResponse[Resp], *httpError)
)

func JSONDecoderMiddleware[Req Validatable, Resp any](
	handler handlerFunc[Req, Resp],
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

		resp, err := handler(w, r, params)
		if err != nil {
			http.Error(w, err.Message, err.StatusCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp.Data)
	}
}

func JSONDecoderMiddlewareGet[Resp any](
	handler getHandlerFunc[Resp],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isJsonRequest := r.Header.Get("Content-Type") == "application/json"
		if !isJsonRequest {
			http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
			return
		}

		resp, err := handler(w, r)
		if err != nil {
			http.Error(w, err.Message, err.StatusCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp.Data)
	}
}

func MountRoutes(
	r chi.Router,
	dbc *sqlx.DB,
	querier db.Querier,
	userCreatedEventService pubsub.UserCreatedEventServicer,
) {
	r.Route("/api/v1/query", func(r chi.Router) {
		r.Get("/user", JSONDecoderMiddlewareGet(
			GetUser(dbc, querier),
		))
		r.Get("/lead/{id}", GetLead())
		r.Get("/contact/{id}", GetContact())
		r.Get("/task/{id}", GetTask())
	})

	r.Route("/api/v1/user", func(r chi.Router) {
		r.Post("/create", JSONDecoderMiddleware(
			CreateUser(dbc, querier, userCreatedEventService),
		))
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
