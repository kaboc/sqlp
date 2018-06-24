package placeholder

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	er "github.com/kaboc/sqlp/errors"
	ref "github.com/kaboc/sqlp/reflect"
)

func replaceQuotes(query string) (string, []string) {
	exp := regexp.MustCompile(`(?m)'(\\'|[^'])*'`)
	return exp.ReplaceAllString(query, "''"), exp.FindAllString(query, -1)
}

func restoreQuotes(query string, orgQuotes []string) string {
	for _, v := range orgQuotes {
		query = strings.Replace(query, "''", v, 1)
	}
	return query
}

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

	query, quotes := replaceQuotes(query)
	query = unnamedToStd(query)
	query = restoreQuotes(query, quotes)

	return query, bind, nil
}

func unnamedToStd(query string) string {
	exp := regexp.MustCompile(`(?im)(\s+in)\s+\?\[(\d*)\]([\s\),]|$)`)
	matches := exp.FindAllStringSubmatch(query, -1)

	for _, v := range matches {
		num, _ := strconv.Atoi(v[2])
		query = exp.ReplaceAllString(query, "$1 ("+strings.Repeat("?,", num)[:num*2-1]+")$3")
	}

	if convertFunc != nil {
		convertFunc(&query)
	}

	return query
}

func convertNamed(query string, arg map[string]interface{}) (string, []interface{}, error) {
	query, quotes := replaceQuotes(query)

	queryUnnamed, bindModel := namedToStd(query)

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

	mapLen := len(arg)
	orderSort := make(map[int]string, mapLen)
	indexes := make([]int, mapLen)

	i := 0
	for name := range arg {
		index := strings.Index(query, ":"+name)
		orderSort[index] = name
		indexes[i] = index
		i++
	}
	sort.Ints(indexes)

	var bind []interface{}
	for _, index := range indexes {
		name := orderSort[index]

		v := ref.ValueOf(arg[name])
		if ref.IsSlice(v) {
			for i = 0; i < v.Len(); i++ {
				bind = append(bind, v.Index(i).Interface())
			}
		} else {
			bind = append(bind, arg[name])
		}
	}

	queryUnnamed = restoreQuotes(queryUnnamed, quotes)

	return queryUnnamed, bind, nil
}

func namedToStd(query string) (string, map[string]interface{}) {
	bindModel := make(map[string]interface{})

	exp := regexp.MustCompile(`(?im)(\s+in)\s+:([^\s\[\),]+)\[(\d*)\]([\s\),]|$)`)
	matches := exp.FindAllStringSubmatch(query, -1)

	for _, v := range matches {
		num, _ := strconv.Atoi(v[3])
		query = exp.ReplaceAllString(query, "$1 ("+strings.Repeat("?,", num)[:num*2-1]+")$4")
		bindModel[v[2]] = make([]interface{}, num)
	}

	exp = regexp.MustCompile(`(?m):([^\s\[\),]+)([\s\),]|$)`)
	matches = exp.FindAllStringSubmatch(query, -1)

	for _, v := range matches {
		query = strings.Replace(query, ":"+v[1], "?", 1)
		bindModel[v[1]] = nil
	}

	if convertFunc != nil {
		convertFunc(&query)
	}

	return query, bindModel
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
	q, quotes := replaceQuotes(query)

	var queryUnnamed string

	if strings.Contains(q, ":") {
		queryUnnamed, _ = namedToStd(q)
	} else {
		queryUnnamed = unnamedToStd(query)
	}

	queryUnnamed = restoreQuotes(queryUnnamed, quotes)

	return queryUnnamed, nil
}
