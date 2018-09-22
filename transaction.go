package sqlp

import (
	"context"
	"database/sql"
)

type Tx struct {
	SqlTx *sql.Tx
}

func (db *DB) Begin() (*Tx, error) {
	tx := new(Tx)
	sqlTx, err := db.SqlDB.Begin()
	tx.SqlTx = sqlTx

	return tx, err
}

func (db *DB) BeginTx(ctx context.Context, opts *TxOptions) (*Tx, error) {
	tx := new(Tx)
	sqlTx, err := db.SqlDB.BeginTx(ctx, opts)
	tx.SqlTx = sqlTx

	return tx, err
}

func (tx *Tx) Commit() error {
	return tx.SqlTx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.SqlTx.Rollback()
}

func (tx *Tx) sqlExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return tx.SqlTx.ExecContext(ctx, query, args...)
}

func (tx *Tx) sqlQueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return tx.SqlTx.QueryContext(ctx, query, args...)
}

func (tx *Tx) sqlPrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return tx.SqlTx.PrepareContext(ctx, query)
}

func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return execContext(ctx, tx, query, args...)
}

func (tx *Tx) Exec(query string, args ...interface{}) (Result, error) {
	return tx.ExecContext(context.Background(), query, args...)
}

func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	return queryContext(ctx, tx, query, args...)
}

func (tx *Tx) Query(query string, args ...interface{}) (*Rows, error) {
	return tx.QueryContext(context.Background(), query, args...)
}

func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	return queryRowContext(ctx, tx, query, args...)
}

func (tx *Tx) QueryRow(query string, args ...interface{}) *Row {
	return tx.QueryRowContext(context.Background(), query, args...)
}

func (tx *Tx) InsertContext(ctx context.Context, tableName string, structSlice interface{}) (Result, error) {
	return insertContext(ctx, tx, tableName, structSlice)
}

func (tx *Tx) Insert(tableName string, structSlice interface{}) (Result, error) {
	return tx.InsertContext(context.Background(), tableName, structSlice)
}

func (tx *Tx) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	return prepareContext(ctx, tx, query)
}

func (tx *Tx) Prepare(query string) (*Stmt, error) {
	return tx.PrepareContext(context.Background(), query)
}

func (tx *Tx) SelectToStructContext(ctx context.Context, structSlicePtr interface{}, query string, args ...interface{}) error {
	return selectToStructContext(ctx, tx, structSlicePtr, query, args...)
}

func (tx *Tx) SelectToStruct(structSlicePtr interface{}, query string, args ...interface{}) error {
	return tx.SelectToStructContext(context.Background(), structSlicePtr, query, args...)
}

func (tx *Tx) SelectToMapContext(ctx context.Context, query string, args ...interface{}) ([]map[string]string, error) {
	return selectToMapContext(ctx, tx, query, args...)
}

func (tx *Tx) SelectToMap(query string, args ...interface{}) ([]map[string]string, error) {
	return tx.SelectToMapContext(context.Background(), query, args...)
}

func (tx *Tx) SelectToSliceContext(ctx context.Context, query string, args ...interface{}) ([][]string, error) {
	return selectToSliceContext(ctx, tx, query, args...)
}

func (tx *Tx) SelectToSlice(query string, args ...interface{}) ([][]string, error) {
	return tx.SelectToSliceContext(context.Background(), query, args...)
}
