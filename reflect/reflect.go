package reflect

import (
	"reflect"
)

type Type = reflect.Type

type StructField struct {
	Name string
	Tag  string
	Type Type
}

func ValueOf(i interface{}) reflect.Value {
	return reflect.ValueOf(i)
}

func PtrValueOf(i interface{}) reflect.Value {
	return reflect.ValueOf(i).Elem()
}

func PtrToValue(v reflect.Value) reflect.Value {
	return v.Elem()
}

func GetType(v reflect.Value) reflect.Type {
	if IsSlice(v) {
		return v.Type().Elem()
	}
	return v.Type()
}

func IsPtr(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr
}

func IsSlice(v reflect.Value) bool {
	return v.Kind() == reflect.Slice
}

func IsMap(v reflect.Value) bool {
	return v.Kind() == reflect.Map
}

func IsStruct(v reflect.Value) bool {
	return GetType(v).Kind() == reflect.Struct
}

func MakeStruct(t reflect.Type) reflect.Value {
	return reflect.New(t).Elem()
}

func GetStructFields(v reflect.Value) []StructField {
	structType := GetType(v)
	structV := MakeStruct(structType)

	fieldLen := structV.NumField()
	names := make([]StructField, fieldLen)

	for i := 0; i < fieldLen; i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("col")

		if tag == "" {
			tag = field.Name
		}

		names[i] = StructField{
			Name: field.Name,
			Tag:  tag,
			Type: field.Type,
		}
	}

	return names
}

func Append(v1, v2 reflect.Value) reflect.Value {
	return reflect.Append(v1, v2)
}
