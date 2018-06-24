package sqlp

import (
	"context"
	"strings"

	ref "github.com/kaboc/sqlp/reflect"
)

func insertContext(ctx context.Context, sq sqler, tableName string, structSlice interface{}) (Result, error) {
	var result Result

	if err := isStructSlice(structSlice); err != nil {
		return result, err
	}

	sliceV := ref.ValueOf(structSlice)
	sliceLen := sliceV.Len()

	fields := ref.GetStructFields(sliceV)
	fieldLen := len(fields)

	values := make([]string, sliceLen)
	binds := make([]interface{}, sliceLen*fieldLen)

	for i1 := 0; i1 < sliceLen; i1++ {
		values[i1] = "(" + strings.Repeat("?,", fieldLen)[:fieldLen*2-1] + ")"

		dataV := sliceV.Index(i1)
		for i2, v := range fields {
			binds[i1*fieldLen+i2] = dataV.FieldByName(v.Name).Interface()
		}
	}

	columnNames := make([]string, fieldLen)
	for i3, v := range fields {
		columnNames[i3] = v.Tag
	}

	sqlstmt := `INSERT INTO ` + tableName + `
				(` + strings.Join(columnNames, ",") + `)
				VALUES ` + strings.Join(values, ",")

	res, err := execContext(ctx, sq, sqlstmt, binds...)
	if err != nil {
		return result, err
	}

	return res, nil
}
