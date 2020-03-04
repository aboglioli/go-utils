package errors

import (
	"encoding/json"
	"fmt"
)

type Type string

const (
	Internal   = Type("Internal")
	Status     = Type("Status")
	Validation = Type("Validation")
	Unknown    = Type("Unknown")
)

type Context map[string]string

type Field struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

// Error
type Error struct {
	Type    Type    `json:"type"`
	Code    string  `json:"code"`
	Path    string  `json:"path,omitempty"`
	Message string  `json:"message,omitempty"`
	Status  int     `json:"status,omitempty"`
	Context Context `json:"context,omitempty"`
	Fields  []Field `json:"fields,omitempty"`
	Cause   error   `json:"cause,omitempty"`
}

func (t Type) New(code string) Error {
	return Error{
		Type: t,
		Code: code,
	}
}

func (e Error) P(path string) Error {
	e.Path = path
	return e
}

func (e Error) M(format string, args ...interface{}) Error {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

func (e Error) S(status int) Error {
	e.Status = status
	return e
}

func (e Error) C(key, value string, args ...interface{}) Error {
	if len(args) > 0 {
		value = fmt.Sprintf(value, args...)
	}

	ctx := make(Context)
	for key, value := range e.Context {
		ctx[key] = value
	}
	ctx[key] = value
	e.Context = ctx
	return e
}

func (e Error) F(field string, code string, msgs ...interface{}) Error {
	msg := ""

	if len(msgs) > 0 {
		if format, ok := msgs[0].(string); ok {
			args := msgs[1:]
			msg = fmt.Sprintf(format, args...)
		}
	}

	e.Fields = append(e.Fields, Field{
		Field:   field,
		Code:    code,
		Message: msg,
	})

	return e
}

func (e Error) Wrap(err error) Error {
	e.Cause = err
	return e
}

// Unwrap declared for compatiblity with errors.As(...)
func (e Error) Unwrap() error {
	return e.Cause
}

func (e Error) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(b)
}

func (err1 Error) Equals(err2 Error) bool {
	return err1.Type == err2.Type && err1.Code == err2.Code
}

// Is declared for compatiblity with errors.Is(...)
func (err1 Error) Is(err2 error) bool {
	if err2, ok := err2.(Error); ok {
		return err1.Equals(err2)
	}
	return false
}

// Errors is a collection of errors
type Errors []error

func (e Errors) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(b)
}
