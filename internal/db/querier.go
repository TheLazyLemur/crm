package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type DBExecutor interface {
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	sqlx.Ext
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
