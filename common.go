package sqlp

import (
	er "github.com/kaboc/sqlp/errors"
	ref "github.com/kaboc/sqlp/reflect"
)

func isStructPtr(ptr interface{}) error {
	rv := ref.ValueOf(ptr)
	if ref.IsPtr(rv) && ref.IsStruct(ref.PtrToValue(rv)) {
		return nil
	}
	return er.New("expected a pointer to a struct")
}

func isStructSlice(slice interface{}) error {
	rv := ref.ValueOf(slice)
	if ref.IsSlice(rv) && ref.IsStruct(rv) {
		return nil
	}
	return er.New("expected a slice of struct")
}

func isStructSlicePtr(ptr interface{}) error {
	rv := ref.ValueOf(ptr)
	if ref.IsPtr(rv) {
		rv = ref.PtrToValue(rv)
		if ref.IsSlice(rv) && ref.IsStruct(rv) {
			return nil
		}
	}
	return er.New("expected a pointer to a slice of struct")
}
