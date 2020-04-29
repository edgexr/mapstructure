package mapstructure

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Error implements the error interface and can represents multiple
// errors that occur in the course of a single decode.
type Error struct {
	Errors []error
}

func (e *Error) Error() string {
	points := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		points[i] = fmt.Sprintf("* %s", err)
	}

	sort.Strings(points)
	return fmt.Sprintf(
		"%d error(s) decoding:\n\n%s",
		len(e.Errors), strings.Join(points, "\n"))
}

// WrappedErrors implements the errwrap.Wrapper interface to make this
// return value more useful with the errwrap and go-multierror libraries.
func (e *Error) WrappedErrors() []error {
	if e == nil {
		return nil
	}
	return e.Errors
}

func (e *Error) append(err error) {
	switch inerr := err.(type) {
	case *Error:
		e.Errors = append(e.Errors, inerr.Errors...)
	default:
		e.Errors = append(e.Errors, inerr)
	}
}

type ParseError struct {
	Name string
	Val  interface{}
	To   reflect.Kind
	Err  error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("cannot parse '%s' as %v: %v", e.Name, e.To, e.Err)
}

func parseError(name string, val interface{}, to reflect.Kind, err error) *ParseError {
	return &ParseError{name, val, to, err}
}

type OverflowError struct {
	Name string
	Val  interface{}
	To   reflect.Kind
}

func (e *OverflowError) Error() string {
	return fmt.Sprintf("cannot parse '%s', %v overflows %v",
		e.Name, e.Val, e.To)
}

func overflowError(name string, val interface{}, to reflect.Kind) error {
	return &OverflowError{name, val, to}
}
