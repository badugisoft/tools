package lib

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// Error error information struct
type Error struct {
	// error message
	Message string

	// error position, format: 'file:line'
	Position string

	// nested error
	InnerError error
}

// NewError create Error
func NewError(msg, pos string, err error) Error {
	return Error{
		Message:    msg,
		Position:   pos,
		InnerError: err,
	}
}

func (e Error) String() string {
	return fmt.Sprintf("message: %v\nposition: %v\ninnerError: %v",
		e.Message, e.Position, e.InnerError)
}

func (e Error) Error() string {
	return e.String()
}

func getPosition() string {
	s := string(debug.Stack())
	t := strings.Split(s, "\n")
	if len(t) < 9 {
		return s
	}

	c := t[8]
	i := strings.Index(c, ".")
	if i < 0 {
		return c
	}

	c = c[i+1:]
	t = strings.SplitN(c, "/", 4)
	if len(t) < 4 {
		return c
	}

	return t[3]
}

// Panic panic
func Panic(msg string) {
	panic(NewError(msg, getPosition(), nil))
}

// Panicf panic with format
func Panicf(format string, args ...interface{}) {
	panic(NewError(fmt.Sprintf(format, args...), getPosition(), nil))
}

// PanicIf panic if err is not nil
func PanicIf(err error) {
	if err != nil {
		panic(NewError(err.Error(), getPosition(), nil))
	}
}

// PanicfIf panic with format if err is not nil
func PanicfIf(err error, format string, args ...interface{}) {
	if err != nil {
		panic(NewError(fmt.Sprintf(format, args...), getPosition(), err))
	}
}
