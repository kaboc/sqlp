package placeholder

import (
	"strconv"
	"strings"
)

const (
	Question int = iota
	Dollar
)

var (
	convertFunc func(*string)
)

func SetType(t int) {
	switch t {
	case Dollar:
		convertFunc = convertToDollar
	}
}

func SetConvertFunc(f func(*string)) {
	convertFunc = f
}

func convertToDollar(query *string) {
	cnt := strings.Count(*query, "?")
	for i := 1; i <= cnt; i++ {
		*query = strings.Replace(*query, "?", "$"+strconv.Itoa(i), 1)
	}
}
