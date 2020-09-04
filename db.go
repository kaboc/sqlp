package sqlp

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"
)

type DB struct {
	SqlDB *sql.DB
}

func newDB(sqlDB *sql.DB) *DB {
	db := new(DB)
	db.SqlDB = sqlDB

	return db
}

func Open(driverName string, dataSourceName string) (*DB, error) {
	sqlDB, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	db := newDB(sqlDB)

	return db, nil
}

func OpenDB(c driver.Connector) *DB {
	return newDB(sql.OpenDB(c))
}

func Init(sqlDB *sql.DB) *DB {
	return newDB(sqlDB)
}

func (db *DB) Close() error {
	return db.SqlDB.Close()
}

func (db *DB) sqlExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.SqlDB.ExecContext(ctx, query, args...)
}

func (db *DB) sqlQueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.SqlDB.QueryContext(ctx, query, args...)
}

func (db *DB) sqlPrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return db.SqlDB.PrepareContext(ctx, query)
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return execContext(ctx, db, query, args...)
}

func (db *DB) Exec(query string, args ...interface{}) (Result, error) {
	return db.ExecContext(context.Background(), query, args...)
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	return queryContext(ctx, db, query, args...)
}

func (db *DB) Query(query string, args ...interface{}) (*Rows, error) {
	return db.QueryContext(context.Background(), query, args...)
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	return queryRowContext(ctx, db, query, args...)
}

func (db *DB) QueryRow(query string, args ...interface{}) *Row {
	return db.QueryRowContext(context.Background(), query, args...)
}

func (db *DB) InsertContext(ctx context.Context, tableName string, structSlice interface{}) (Result, error) {
	return insertContext(ctx, db, tableName, structSlice)
}

func (db *DB) Insert(tableName string, structSlice interface{}) (Result, error) {
	return db.InsertContext(context.Background(), tableName, structSlice)
}

func (db *DB) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	return prepareContext(ctx, db, query)
}

func (db *DB) Prepare(query string) (*Stmt, error) {
	return db.PrepareContext(context.Background(), query)
}

func (db *DB) SelectToStructContext(ctx context.Context, structSlicePtr interface{}, query string, args ...interface{}) error {
	return selectToStructContext(ctx, db, structSlicePtr, query, args...)
}

func (db *DB) SelectToStruct(structSlicePtr interface{}, query string, args ...interface{}) error {
	return db.SelectToStructContext(context.Background(), structSlicePtr, query, args...)
}

func (db *DB) SelectToMapContext(ctx context.Context, query string, args ...interface{}) ([]map[string]string, error) {
	return selectToMapContext(ctx, db, query, args...)
}

func (db *DB) SelectToMap(query string, args ...interface{}) ([]map[string]string, error) {
	return db.SelectToMapContext(context.Background(), query, args...)
}

func (db *DB) SelectToSliceContext(ctx context.Context, query string, args ...interface{}) ([][]string, error) {
	return selectToSliceContext(ctx, db, query, args...)
}

func (db *DB) SelectToSlice(query string, args ...interface{}) ([][]string, error) {
	return db.SelectToSliceContext(context.Background(), query, args...)
}

func (db *DB) Driver() driver.Driver {
	return db.SqlDB.Driver()
}

func (db *DB) Ping() error {
	return db.SqlDB.Ping()
}

func (db *DB) PingContext(ctx context.Context) error {
	return db.SqlDB.PingContext(ctx)
}

func (db *DB) SetConnMaxIdleTime(d time.Duration) {
	db.SqlDB.SetConnMaxIdleTime(d)
}

func (db *DB) SetConnMaxLifetime(d time.Duration) {
	db.SqlDB.SetConnMaxLifetime(d)
}

func (db *DB) SetMaxIdleConns(n int) {
	db.SqlDB.SetMaxIdleConns(n)
}

func (db *DB) SetMaxOpenConns(n int) {
	db.SqlDB.SetMaxOpenConns(n)
}

func (db *DB) Stats() DBStats {
	return db.SqlDB.Stats()
}
