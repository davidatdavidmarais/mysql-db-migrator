package migrator

import "database/sql"

type Backends interface {
	DB() *sql.DB
}

type Migrator struct{}
