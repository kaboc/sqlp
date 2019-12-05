package placeholder

import (
	"regexp"
	"strconv"
	"strings"
)

type repl struct {
	query string
	orgs  []string
}

const tempReplacement = "/**SQLP_REPLACE**/"

func replace(query string) *repl {
	p1 := `'(\\'|[^'])*?'`
	p2 := `"(\\"|[^"])*?"`
	p3 := "`[^`]*?`"
	p4 := `/\*.*?\*/`
	p5 := `(#|\s+--).*?([\r\n]|$)`
	exp := regexp.MustCompile(`(?s)(` + p1 + `|` + p2 + `|` + p3 + `|` + p4 + `|` + p5 + `)`)

	return &repl{
		query: exp.ReplaceAllString(query, tempReplacement),
		orgs:  exp.FindAllString(query, -1),
	}
}

func (r *repl) restore() string {
	for _, v := range r.orgs {
		r.query = strings.Replace(r.query, tempReplacement, v, 1)
	}
	return r.query
}

func (r *repl) unnamedToStd() {
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

func (r *repl) namedToStd() map[string]interface{} {
	bindModel := make(map[string]interface{})

	// Replacement of in :xxxx[xx]
	exp := regexp.MustCompile(`(?im)(\s+in)\s+:([a-z0-9_]+)\[(\d*)\]([\s\)/,]|` + regexp.QuoteMeta(tempReplacement) + `|$)`)
	matches := exp.FindAllStringSubmatch(r.query, -1)

	for _, v := range matches {
		num, _ := strconv.Atoi(v[3])
		r.query = exp.ReplaceAllString(r.query, "$1 ("+strings.Repeat("?,", num)[:num*2-1]+")$4")
		bindModel[v[2]] = make([]interface{}, num)
	}

	// Replacement of :xxxx
	exp = regexp.MustCompile(`(?m):([a-zA-Z0-9_]+)([\s\)/,]|$)`)
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
