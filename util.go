package cfm

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"reflect"
)

const blankCutSet = " \t\r\n"

type stack struct {
	l *list.List
}

func newStack() (s *stack) {
	s = new(stack)
	s.l = list.New()
	return
}

func (s *stack) push(ctx *Context) {
	s.l.PushBack(ctx)
}

func (s *stack) pop() (ctx *Context) {
	var ok bool

	if e := s.l.Back(); e != nil {
		v := s.l.Remove(e)
		if ctx, ok = v.(*Context); ok {
			return
		}
	}

	return nil
}

func (s *stack) top() (ctx *Context) {
	var ok bool

	if e := s.l.Back(); e != nil {
		if ctx, ok = e.Value.(*Context); ok {
			return
		}
		ctx = nil
	}
	return
}

func (s *stack) size() int {
	return s.l.Len()
}

func skip(c byte) bool {
	return isBlank(c)
}

func isBlank(c byte) bool {
	if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
		return true
	}
	return false
}

// delete one '"' of string's left and right
func trimString(s string) string {
	if len(s) != 0 && s[0] == '"' {
		s = s[1:]
	}

	if l := len(s); l != 0 && s[l-1] == '"' {
		s = s[:l-1]
	}

	return s
}

func GetStructField(s interface{}, field string, kind reflect.Kind) (v reflect.Value, err error) {
	v = reflect.ValueOf(s)
	if v.IsNil() && !v.IsValid() {
		err = errors.New("Not valid")
		return
	}

	if v = v.Elem(); v.Kind() != reflect.Struct {
		err = fmt.Errorf("%s not struct", v.Type().Name())
		return
	}

	structName := v.Type().Name()

	if v = v.FieldByName(field); v.Kind() != kind {
		err = fmt.Errorf("%s not contain field: %s (or type error)", structName, field)
		return
	}

	if !v.CanSet() {
		err = fmt.Errorf("can't set %s.%s", structName, field)
		return
	}

	return
}

func splitCommandEntry(entry []byte) (fields []string) {

	splitFunc := func(s []byte) (string, []byte) {
		s = bytes.Trim(s, blankCutSet)

		for i := 0; i < len(s); i++ {
			c := s[i]

			if isBlank(c) {
				return string(s[:i]), s[i:]
			}
		}

		return string(s), nil
	}

	fields = make([]string, 0, 3)
	var f string

	for {
		if f, entry = splitFunc(entry); f != "" {
			fields = append(fields, f)
		}

		if entry == nil {
			break
		}
	}
	return
}
