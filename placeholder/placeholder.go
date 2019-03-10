package placeholder

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	er "github.com/kaboc/sqlp/errors"
	ref "github.com/kaboc/sqlp/reflect"
)

type rep struct {
	query string
	orgs  []string
}

const TEMP_REPLACEMENT = "/**SQLP_REPLACE**/"

func replace(query string) *rep {
	p1 := `'(\\'|[^'])*?'`
	p2 := `"(\\"|[^"])*?"`
	p3 := "`[^`]*?`"
	p4 := `/\*.*?\*/`
	p5 := `(#|\s+--).*?([\r\n]|$)`
	exp := regexp.MustCompile(`(?s)(` + p1 + `|` + p2 + `|` + p3 + `|` + p4 + `|` + p5 + `)`)

	return &rep{
		query: exp.ReplaceAllString(query, TEMP_REPLACEMENT),
		orgs:  exp.FindAllString(query, -1),
	}
}

func (r *rep) restore() string {
	for _, v := range r.orgs {
		r.query = strings.Replace(r.query, TEMP_REPLACEMENT, v, 1)
	}
	return r.query
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

	r := replace(query)
	r.unnamedToStd()

	return r.restore(), bind, nil
}

func (r *rep) unnamedToStd() {
	exp := regexp.MustCompile(`(?im)(\s+in)\s+\?\[(\d*)\]([\s\),]|` + regexp.QuoteMeta(TEMP_REPLACEMENT) + `|$)`)
	matches := exp.FindAllStringSubmatch(r.query, -1)

	for _, v := range matches {
		num, _ := strconv.Atoi(v[2])
		r.query = exp.ReplaceAllString(r.query, "$1 ("+strings.Repeat("?,", num)[:num*2-1]+")$3")
	}

	if convertFunc != nil {
		convertFunc(&r.query)
	}
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

	return r.restore(), bind, nil
}

func (r *rep) namedToStd() map[string]interface{} {
	bindModel := make(map[string]interface{})

	exp := regexp.MustCompile(`(?im)(\s+in)\s+:([^\s\[\),]+)\[(\d*)\]([\s\)/,]|` + regexp.QuoteMeta(TEMP_REPLACEMENT) + `|$)`)
	matches := exp.FindAllStringSubmatch(r.query, -1)

	for _, v := range matches {
		num, _ := strconv.Atoi(v[3])
		r.query = exp.ReplaceAllString(r.query, "$1 ("+strings.Repeat("?,", num)[:num*2-1]+")$4")
		bindModel[v[2]] = make([]interface{}, num)
	}

	exp = regexp.MustCompile(`(?m):([^\s\[\)/,]+)([\s\)/,]|$)`)
	matches = exp.FindAllStringSubmatch(r.query, -1)

	for _, v := range matches {
		r.query = strings.Replace(r.query, ":"+v[1], "?", 1)
		bindModel[v[1]] = nil
	}

	if convertFunc != nil {
		convertFunc(&r.query)
	}

	return bindModel
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
