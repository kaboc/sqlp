package sqlp

import (
	"reflect"
	"testing"

	ref "github.com/kaboc/sqlp/reflect"
)

type sample struct {
	Col1 string `col:"column1"`
	Col2 int    `col:"column2"`
}

func TestCheckStructPtr(t *testing.T) {
	err := isStructPtr(&sample{})
	if err != nil {
		t.Fatal("isStructPtr() failed")
	}
}

func TestCheckStructSlice(t *testing.T) {
	err := isStructSlice([]sample{})
	if err != nil {
		t.Fatal("isStructSlice() failed")
	}
}

func TestCheckStructSlicePtr(t *testing.T) {
	err := isStructSlicePtr(&[]sample{})
	if err != nil {
		t.Fatal("isStructSlicePtr() failed")
	}
}

func TestGetStructFields(t *testing.T) {
	var s []sample
	actual := ref.GetStructFields(ref.ValueOf(s))

	expected := []ref.StructField{
		{
			Name: "Col1",
			Tag:  "column1",
			Type: reflect.TypeOf(""),
		},
		{
			Name: "Col2",
			Tag:  "column2",
			Type: reflect.TypeOf(0),
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatal("got wrong struct fields")
	}
}
