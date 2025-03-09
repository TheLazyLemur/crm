package db

import (
	"context"
)

type Queries struct{}

func (q *Queries) GetUser(ctx context.Context, dbc DBExecutor, id string) (User, error) {
	panic("not implemented") // TODO: Implement
}

func (q *Queries) InsertAndReturnUser(
	ctx context.Context,
	dbc DBExecutor,
	arg InsertAndReturnUserParams,
) (User, error) {
	query := `
	INSERT INTO users (id, first_name, last_name, email) 
	VALUES (:id, :first_name, :last_name, :email) RETURNING id, first_name, last_name, email, created_at
	`

	query, args, err := dbc.BindNamed(query, map[string]any{
		"id":         arg.ID,
		"first_name": arg.FirstName,
		"last_name":  arg.LastName,
		"email":      arg.Email,
	})
	if err != nil {
		return User{}, err
	}

	var user User
	err = dbc.GetContext(ctx, &user, query, args...)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func NewQueries() Querier {
	return &Queries{}
}
