package controller

import (
	"context"
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func HasQueryBeenExecuted(ctx context.Context, b Backends, hash string) (bool, error) {
	e, err := controllerExists(ctx, b)
	if err != nil {
		return false, err
	}

	if !e {
		err = createController(ctx, b)
		if err != nil {
			return false, err
		}
	}

	return queryStatus(ctx, b, hash)
}

func queryStatus(ctx context.Context, b Backends, hash string) (bool, error) {
	r, err := b.DB().QueryContext(ctx, "select hash from migration_controller where hash = ?", hash)
	if err != nil {
		return false, err
	}

	return r.Next(), nil
}

func InsertExecutedQuery(ctx context.Context, tx *sql.Tx, hash string) error {
	_, err := tx.ExecContext(ctx, "insert into migration_controller (hash, executed) values (?, ?)", hash, true)
	return err
}

func createController(ctx context.Context, b Backends) error {
	_, err := b.DB().ExecContext(ctx, "create table migration_controller (id int not null auto_increment, hash varchar(255) not null, executed bool not null, primary key (id))")
	return err
}

func controllerExists(ctx context.Context, b Backends) (bool, error) {
	_, err := b.DB().QueryContext(ctx, "select * from migration_controller")
	if err != nil {
		if strings.Contains(err.Error(), "doesn't exist") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
