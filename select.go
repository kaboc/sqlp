package sqlp

import (
	"context"

	ref "github.com/kaboc/sqlp/reflect"
)

func selectToStructContext(ctx context.Context, sq sqler, structSlicePtr interface{}, query string, args ...interface{}) error {
	err := isStructSlicePtr(structSlicePtr)
	if err != nil {
		return err
	}

	rows, err := queryContext(ctx, sq, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	sliceV := ref.PtrValueOf(structSlicePtr)
	structV := ref.MakeStruct(ref.GetType(sliceV))
	structPtr := structV.Addr().Interface()

	newSliceV := sliceV.Slice(0, 0)
	for rows.Next() {
		err = rows.ScanToStruct(structPtr)
		if err != nil {
			return err
		}
		newSliceV = ref.Append(newSliceV, structV)
	}
	sliceV.Set(newSliceV)

	return nil
}

func selectToMapContext(ctx context.Context, sq sqler, query string, args ...interface{}) ([]map[string]string, error) {
	rows, err := queryContext(ctx, sq, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]string

	for rows.Next() {
		data, err := rows.ScanToMap()
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}

	return result, nil
}

func selectToSliceContext(ctx context.Context, sq sqler, query string, args ...interface{}) ([][]string, error) {
	rows, err := queryContext(ctx, sq, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result [][]string

	for rows.Next() {
		data, err := rows.ScanToSlice()
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}

	return result, nil
}
