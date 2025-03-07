package database

import (
	_ "embed"
)

//go:embed migrations/000001_Initial.sql
var Migrations []byte
