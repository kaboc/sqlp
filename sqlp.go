package sqlp

import (
	"context"
	"database/sql"

	"github.com/kaboc/sqlp/placeholder"
)

type sqler interface {
	sqlExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	sqlQueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	sqlPrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

func execContext(ctx context.Context, sq sqler, query string, args ...interface{}) (Result, error) {
	var result Result

	query, bind, err := placeholder.Convert(query, args...)
	if err != nil {
		return result, err
	}

	res, err := sq.sqlExecContext(ctx, query, bind...)
	if err != nil {
		return result, err
	}

	affectedRows, _ := res.RowsAffected()
	insertId, _ := res.LastInsertId()

	result = Result{
		affectedRows: affectedRows,
		insertId:     insertId,
	}

	return result, err
}

func queryContext(ctx context.Context, sq sqler, query string, args ...interface{}) (*Rows, error) {
	var sqlRows *sql.Rows

	query, bind, err := placeholder.Convert(query, args...)
	if err == nil {
		sqlRows, err = sq.sqlQueryContext(ctx, query, bind...)
	}

	return &Rows{Rows: sqlRows}, err
}

func queryRowContext(ctx context.Context, sq sqler, query string, args ...interface{}) *Row {
	rows, err := queryContext(ctx, sq, query, args...)
	return &Row{rows: rows, err: err}
}

func prepareContext(ctx context.Context, sq sqler, query string) (*Stmt, error) {
	queryUnnamed, err := placeholder.ConvertSQL(query)
	if err != nil {
		return nil, err
	}

	stmt, err := sq.sqlPrepareContext(ctx, queryUnnamed)

	return &Stmt{
		SqlStmt: stmt,
		query:   query, // Stores original query.
	}, err
}
