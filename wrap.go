package sqlp

import (
	"database/sql"
	"database/sql/driver"
)

type columnTypes = sql.ColumnType

type NullBool = sql.NullBool
type NullFloat64 = sql.NullFloat64
type NullInt64 = sql.NullInt64
type NullString = sql.NullString

type RawBytes = sql.RawBytes

type IsolationLevel = sql.IsolationLevel
type TxOptions = sql.TxOptions

type DBStats = sql.DBStats

var (
	ErrConnDone = sql.ErrConnDone
	ErrNoRows   = sql.ErrNoRows
	ErrTxDone   = sql.ErrTxDone
)

const (
	LevelDefault IsolationLevel = iota
	LevelReadUncommitted
	LevelReadCommitted
	LevelWriteCommitted
	LevelRepeatableRead
	LevelSnapshot
	LevelSerializable
	LevelLinearizable
)

func Drivers() []string {
	return sql.Drivers()
}

func Register(name string, driver driver.Driver) {
	sql.Register(name, driver)
}
