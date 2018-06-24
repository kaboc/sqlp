package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

func trace() []string {
	pc := make([]uintptr, 20)
	n := runtime.Callers(3, pc)

	var msg []string
	for _, v := range pc[:n] {
		f := runtime.FuncForPC(v)
		if strings.Index(f.Name(), "runtime.") == 0 {
			break
		}

		file, line := f.FileLine(v)
		msg = append(msg, fmt.Sprintf("%s:%d", file, line))
	}

	return msg
}

func New(text string) error {
	return errors.New("sqlp: " + text + "\n" + strings.Join(trace(), "\n"))
}

func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf("sqlp: "+format+"\n"+strings.Join(trace(), "\n"), a...)
}
