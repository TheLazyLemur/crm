package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"simplecrm/internal/db"
	"simplecrm/internal/handlers"
	"simplecrm/internal/pubsub"
)

func main() {
	userCreatedEventService := pubsub.NewUserCreatedEventService()

	userCreatedConsumer := func(event pubsub.UserCreatedEvent) {
		slog.Info("User created event received", "user", event.User.FirstName)
	}

	go userCreatedEventService.Consume(context.Background(), userCreatedConsumer)

	dbc, err := sqlx.Connect("sqlite3", "./simplecrm.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer dbc.Close()

	querier := db.NewQueries()

	r := chi.NewRouter()

	handlers.MountRoutes(r, dbc, querier, userCreatedEventService)

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	slog.Info("Server started", "addr", server.Addr)
	log.Fatal(server.ListenAndServe())
}
