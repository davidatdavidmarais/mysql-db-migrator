package migrator

import (
	"context"

	"github.com/davidatdavidmarais/mysql-db-migrator/controller"
	"github.com/davidatdavidmarais/mysql-db-migrator/parser"
)

func (m *Migrator) Run(ctx context.Context, b Backends, file []byte) error {
	return RunMigration(ctx, b, file)
}

func RunMigration(ctx context.Context, b Backends, file []byte) error {
	ql, err := parser.Parse(string(file))
	if err != nil {
		return err
	}

	for _, q := range ql {
		err := processQuery(ctx, b, q.Query, q.Hash)
		if err != nil {
			return err
		}
	}

	return nil
}

func processQuery(ctx context.Context, b Backends, query, hash string) error {
	e, err := controller.HasQueryBeenExecuted(ctx, b, hash)
	if err != nil {
		return err
	}

	if e {
		return nil
	}

	err = executeQuery(ctx, b, query, hash)
	if err != nil {
		return err
	}

	return nil
}

func executeQuery(ctx context.Context, b Backends, query, hash string) error {
	tx, err := b.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = controller.InsertExecutedQuery(ctx, tx, hash)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
