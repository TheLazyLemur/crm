package ops

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"simplecrm/internal/db"
)

func CreateUser(
	ctx context.Context,
	dbc *sqlx.DB,
	querier db.Querier,
	firstName, lastName, email string,
) (user db.User, err error) {
	tx, err := dbc.Beginx()
	if err != nil {
		return db.User{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	id := uuid.New().String()
	user, err = querier.InsertAndReturnUser(ctx, tx, db.InsertAndReturnUserParams{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	})
	if err != nil {
		return db.User{}, err
	}

	if err := publishUserCreatedEvent(ctx, user); err != nil {
		return db.User{}, err
	}

	return user, nil
}

func publishUserCreatedEvent(
	ctx context.Context,
	user db.User,
) error {
	_ = user
	_ = ctx
	// TODO: Need to make a decision on which messaging system to use, for now I am thinking in a seperate sqlite db and write very specific consumer logic
	return nil
}
