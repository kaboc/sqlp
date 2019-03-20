package placeholder

import (
	"regexp"
	"strconv"
	"strings"
)

type rep struct {
	query string
	orgs  []string
}

const tempReplacement = "/**SQLP_REPLACE**/"

func replace(query string) *rep {
	p1 := `'(\\'|[^'])*?'`
	p2 := `"(\\"|[^"])*?"`
	p3 := "`[^`]*?`"
	p4 := `/\*.*?\*/`
	p5 := `(#|\s+--).*?([\r\n]|$)`
	exp := regexp.MustCompile(`(?s)(` + p1 + `|` + p2 + `|` + p3 + `|` + p4 + `|` + p5 + `)`)

	return &rep{
		query: exp.ReplaceAllString(query, tempReplacement),
		orgs:  exp.FindAllString(query, -1),
	}
}

func (r *rep) restore() string {
	for _, v := range r.orgs {
		r.query = strings.Replace(r.query, tempReplacement, v, 1)
	}
	return r.query
}

func (r *rep) unnamedToStd() {
	exp := regexp.MustCompile(`(?im)(\s+in)\s+\?\[(\d*)\]([\s\),]|` + regexp.QuoteMeta(tempReplacement) + `|$)`)
	matches := exp.FindAllStringSubmatch(r.query, -1)

	for _, v := range matches {
		num, _ := strconv.Atoi(v[2])
		r.query = exp.ReplaceAllString(r.query, "$1 ("+strings.Repeat("?,", num)[:num*2-1]+")$3")
	}

	if convertFunc != nil {
		convertFunc(&r.query)
	}
}

func (r *rep) namedToStd() map[string]interface{} {
	bindModel := make(map[string]interface{})

	exp := regexp.MustCompile(`(?im)(\s+in)\s+:([^\s\[\),]+)\[(\d*)\]([\s\)/,]|` + regexp.QuoteMeta(tempReplacement) + `|$)`)
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
