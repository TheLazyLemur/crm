package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"

	"simplecrm/internal/db"
	"simplecrm/internal/handlers"
)

func main() {
	dbc, err := sql.Open("sqlite3", "./simplecrm.db")
	if err != nil {
		panic(err)
	}

	querier := db.New()
	querier.GetUser(context.Background(), dbc, "1")

	r := chi.NewRouter()

	handlers.MountRoutes(r, dbc, querier)

	http.ListenAndServe(":8080", r)
}
