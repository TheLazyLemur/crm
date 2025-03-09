package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"simplecrm/internal/db"
	"simplecrm/internal/handlers"
)

func main() {
	dbc, err := sqlx.Connect("sqlite3", "./simplecrm.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer dbc.Close()

	querier := db.NewQueries()
	querier.GetUser(context.Background(), dbc, "1")

	r := chi.NewRouter()

	handlers.MountRoutes(r, dbc, querier)

	http.ListenAndServe(":8080", r)
}
