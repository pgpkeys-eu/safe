package safe

import (
	"errors"
)

// Result represents a tagged-union of T and error.
// Its default value represents an OK containing the default value of T.
// Result values SHOULD be created by calling either OK() or Err().
// DO NOT construct non-zero Result literals directly.
type Result[T any] result[T]

// Hide the internals slightly by defining Result as a type alias, like Option.
// TODO: why can we not alias Result to optionN? Unwraping throws a wobbler.
// error is an interface, any is an interface, are they not the same size?
type result[T any] struct {
	tag   error
	value T
}

// Unwrap() converts a Result into a T or an error.
func (r Result[T]) Unwrap() any {
	if r.tag == nil {
		return r.value
	} else {
		return r.tag
	}
}

// OK() constructs an OK Result from an existing value.
func OK[T any](value T) (r Result[T]) {
	r.value = value
	return
}

// Err() constructs an Err Result from an existing error.
// If the error is nil, an error containing the empty string is returned.
func Err[T any](err error) (r Result[T]) {
	if err == nil {
		r.tag = errors.New("")
	} else {
		r.tag = err
	}
	return
}

// IsOK checks directly if the Result is OK.
func (r Result[T]) IsOK() bool {
	return r.tag == nil
}

// IsErr checks directly if the Result is Err.
func (r Result[T]) IsErr() bool {
	return r.tag != nil
}

// UnwrapOr returns the value of the Result if it is OK.
// It returns a static value if the Result is Err.
func (r Result[T]) UnwrapOr(d *T) *T {
	if r.tag == nil {
		return &r.value
	} else {
		return d
	}
}

// UnwrapOrElse returns the value of the Result if it is OK.
// It returns a calculated value if the Result is Err.
func (r Result[T]) UnwrapOrElse(f func() *T) *T {
	if r.tag == nil {
		return &r.value
	} else {
		return f()
	}
}
