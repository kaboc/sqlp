package sqlp

import (
	"database/sql"
	"strings"

	er "github.com/kaboc/sqlp/errors"
	ref "github.com/kaboc/sqlp/reflect"
)

type Rows struct {
	Rows        *sql.Rows
	structPtr   interface{}
	destPtr     *[]interface{}
	columnNames []string
}

func (rs *Rows) Close() error {
	return rs.Rows.Close()
}

func (rs *Rows) Columns() ([]string, error) {
	return rs.Rows.Columns()
}

func (rs *Rows) ColumnTypes() ([]*columnTypes, error) {
	return rs.Rows.ColumnTypes()
}

func (rs *Rows) Err() error {
	return rs.Rows.Err()
}

func (rs *Rows) Next() bool {
	return rs.Rows.Next()
}

func (rs *Rows) NextResultSet() bool {
	return rs.Rows.NextResultSet()
}

func (rs *Rows) Scan(dest ...interface{}) error {
	return rs.Rows.Scan(dest...)
}

func prepareScanToStruct(structPtr interface{}, columnNames []string) (*[]interface{}, error) {
	destPtr := make([]interface{}, len(columnNames))
	structV := ref.PtrValueOf(structPtr)
	fields := ref.GetStructFields(structV)

outerLoop:
	for i, v := range columnNames {
		for _, v2 := range fields {
			if v == v2.Tag || (v2.Tag == v2.Name && strings.EqualFold(strings.Replace(v, "_", "", -1), strings.Replace(v2.Name, "_", "", -1))) {
				if v2.Name != strings.Title(v2.Name) {
					return nil, er.New("one or more destination struct fields are unexported")
				}

				destPtr[i] = structV.FieldByName(v2.Name).Addr().Interface()
				continue outerLoop
			}
		}
		return nil, er.Errorf("struct field for column '%s' is missing", v)
	}

	return &destPtr, nil
}

func (rs *Rows) ScanToStruct(structPtr interface{}) error {
	if structPtr != rs.structPtr {
		if err := isStructPtr(structPtr); err != nil {
			return err
		}

		columnNames, err := rs.Columns()
		if err != nil {
			return err
		}

		destPtr, err := prepareScanToStruct(structPtr, columnNames)
		if err != nil {
			return err
		}

		rs.structPtr = structPtr
		rs.destPtr = destPtr
	}

	err := rs.Rows.Scan(*rs.destPtr...)
	if err != nil {
		return err
	}

	return nil
}

func prepareScanToMap(columnNames []string) *[]interface{} {
	columnLen := len(columnNames)
	destPtr := make([]interface{}, columnLen)
	dest := make([]RawBytes, columnLen)

	for i := range columnNames {
		destPtr[i] = &dest[i]
	}

	return &destPtr
}

func (rs *Rows) ScanToMap() (map[string]string, error) {
	if rs.columnNames == nil {
		columnNames, err := rs.Columns()
		if err != nil {
			return nil, err
		}

		rs.columnNames = columnNames
		rs.destPtr = prepareScanToMap(columnNames)
	}

	if err := rs.Scan(*rs.destPtr...); err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for i, v := range rs.columnNames {
		result[v] = string(*((*rs.destPtr)[i].(*RawBytes)))
	}

	return result, nil
}

func (rs *Rows) ScanToSlice() ([]string, error) {
	mp, err := rs.ScanToMap()
	if err != nil {
		return nil, err
	}

	result := make([]string, len(rs.columnNames))
	for i, v := range rs.columnNames {
		result[i] = mp[v]
	}

	return result, nil
}
