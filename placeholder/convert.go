package placeholder

import (
	"regexp"
	"sort"
	"strings"

	er "github.com/kaboc/sqlp/errors"
	ref "github.com/kaboc/sqlp/reflect"
)

func isNamed(args ...interface{}) bool {
	return len(args) == 1 && ref.IsMap(ref.ValueOf(args[0]))
}

func convertUnnamed(query string, args ...interface{}) (string, []interface{}, error) {
	var bind []interface{}

	if len(args) > 0 {
		for _, arg := range args {
			rv := ref.ValueOf(arg)
			if ref.IsMap(rv) {
				return "", nil, er.New("map cannot be used as arguments for unnamed placeholder parameters")
			} else if ref.IsSlice(rv) {
				bind = append(bind, arg.([]interface{})...)
			} else {
				bind = append(bind, arg)
			}
		}
	}

	r := replace(query)
	r.unnamedToStd()

	return r.restore(), bind, nil
}

func convertNamed(query string, arg map[string]interface{}) (string, []interface{}, error) {
	r := replace(query)

	bindModel := r.namedToStd()

	for name := range arg {
		_, ok := bindModel[name]
		if !ok {
			return "", nil, er.Errorf("binding value is set for an unknown placeholder :%s", name)
		}
	}

	for name, s := range bindModel {
		a, ok := arg[name]
		if !ok {
			return "", nil, er.Errorf("binding value is not set for a placeholder :%s", name)
		}

		rvS := ref.ValueOf(s)
		rvA := ref.ValueOf(a)

		if ref.IsSlice(rvS) {
			if !ref.IsSlice(rvA) {
				return "", nil, er.Errorf("slice is required for :%s[%d]", name, rvS.Len())
			} else if rvA.Len() != rvS.Len() {
				return "", nil, er.Errorf("a slice with %d elements was given for :%s[%d]", rvA.Len(), name, rvS.Len())
			}
		} else if ref.IsSlice(rvA) {
			return "", nil, er.Errorf("a non-slice type is required for :%s", name)
		}
	}

	names := sortNameByIndex(query, arg)
	var bind []interface{}

	for _, name := range names {
		v := ref.ValueOf(arg[name])
		if ref.IsSlice(v) {
			for i := 0; i < v.Len(); i++ {
				bind = append(bind, v.Index(i).Interface())
			}
		} else {
			bind = append(bind, arg[name])
		}
	}

	return r.restore(), bind, nil
}

func sortNameByIndex(query string, arg map[string]interface{}) []string {
	// Add colon to the end of each placeholder name in query string to prevent
	// strings.Index(query, ":foo") from returning the index of ":foobar".
	exp := regexp.MustCompile(`(?m):([a-zA-Z0-9_]+)`)
	if exp.MatchString(query) {
		query = exp.ReplaceAllString(query, "$0:")
	}

	mapLen := len(arg)
	names := make(map[int]string, mapLen)
	indexes := make([]int, mapLen)
	i := 0
	for name := range arg {
		index := strings.Index(query, ":"+name+":")
		names[index] = name
		indexes[i] = index
		i++
	}
	sort.Ints(indexes)

	sorted := make([]string, mapLen)
	for j, index := range indexes {
		sorted[j] = names[index]
	}

	return sorted
}

func Convert(query string, args ...interface{}) (string, []interface{}, error) {
	if isNamed(args...) {
		mp, ok := args[0].(map[string]interface{})
		if ok {
			return convertNamed(query, mp)
		}
	}
	return convertUnnamed(query, args...)
}

func ConvertSQL(query string) (string, error) {
	r := replace(query)

	if strings.Contains(r.query, ":") {
		r.namedToStd()
	} else {
		r.unnamedToStd()
	}

	return r.restore(), nil
}
