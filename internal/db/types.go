package db

import (
	"database/sql"
)

type Entity struct {
	ID          string         `db:"id"`
	FirstName   string         `db:"first_name"`
	LastName    string         `db:"last_name"`
	Email       string         `db:"email"`
	Phone       string         `db:"phone"`
	Status      string         `db:"status"`
	AssignedTo  sql.NullString `db:"assigned_to"`
	CreatedAt   string         `db:"created_at"`
	ConvertedAt string         `db:"converted_at"`
}

type Task struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Description string         `db:"description"`
	DueDate     string         `db:"due_date"`
	AssignedTo  sql.NullString `db:"assigned_to"`
	Status      string         `db:"status"`
}

type User struct {
	ID        string `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	CreatedAt string `db:"created_at"`
}

type InsertAndReturnUserParams struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
}
