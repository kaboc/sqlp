package sqlp

import (
	"context"
	"database/sql"
)

type Stmt struct {
	SqlStmt *sql.Stmt
	query   string
}

func (s *Stmt) Close() error {
	return s.SqlStmt.Close()
}

func (s *Stmt) sqlExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.SqlStmt.ExecContext(ctx, args...)
}

func (s *Stmt) sqlQueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.SqlStmt.QueryContext(ctx, args...)
}

func (s *Stmt) sqlPrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return nil, nil
}

func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (Result, error) {
	return execContext(ctx, s, s.query, args...)
}

func (s *Stmt) Exec(args ...interface{}) (Result, error) {
	return s.ExecContext(context.Background(), args...)
}

func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*Rows, error) {
	return queryContext(ctx, s, s.query, args...)
}

func (s *Stmt) Query(args ...interface{}) (*Rows, error) {
	return s.QueryContext(context.Background(), args...)
}

func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *Row {
	return queryRowContext(ctx, s, s.query, args...)
}

func (s *Stmt) QueryRow(args ...interface{}) *Row {
	return s.QueryRowContext(context.Background(), args...)
}
