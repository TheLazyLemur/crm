package db

import (
	"context"
)

type DBExecutor interface {
	BindNamed(query string, arg any) (string, []any, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
}

type Querier interface {
	GetUser(ctx context.Context, dbc DBExecutor, id string) (User, error)
	InsertAndReturnUser(
		ctx context.Context,
		dbc DBExecutor,
		arg InsertAndReturnUserParams,
	) (User, error)
}

var _ Querier = (*Queries)(nil)
