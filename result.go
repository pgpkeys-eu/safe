package safe

import (
	"errors"
	"unsafe"
)

// Result represents a tagged-union of OK and Err.
// The zero literal Result[T]{} represents an OK containing the zero value of T.
type Result[T any] struct {
	//
	// The zero literal Result[T]{} may be safely used
	// as the static equivalent of OK[T](nil)
	//
	// For everything else, use OK[T]() or Err[T]() instead
	//
	// DO NOT SET STRUCT MEMBER VALUES DIRECTLY
	//
	t error
	v T
}

// MatchOK is a match type with one public member OK, used only for type matching.
// DO NOT construct directly using a struct literal, use OK() or Result[T]{} instead.
type MatchOK[T any] struct {
	//
	// !!!!! DO NOT USE THIS STRUCT LITERAL !!!!!
	//
	// use OK[T]() or Result[T]{} instead
	//
	// !!!!! DO NOT USE THIS STRUCT LITERAL !!!!!
	//
	_  error
	OK T
}

// MatchErr is a match type with one public member Err, used only for type matching.
// DO NOT construct directly using a struct literal, use Err() instead.
type MatchErr[T any] struct {
	//
	// !!!!! DO NOT USE THIS STRUCT LITERAL !!!!!
	//
	// use Err[T]() instead
	//
	// !!!!! DO NOT USE THIS STRUCT LITERAL !!!!!
	//
	Err error
	_   T
}

// Match() converts a Result into a MatchOK or MatchErr.
func (r Result[T]) Match() any {
	// use pointer type mangling to avoid making copies
	// this relies on the memory layouts of all three types being identical
	if r.t == nil {
		return *(*MatchOK[T])(unsafe.Pointer(&r))
	} else {
		return *(*MatchErr[T])(unsafe.Pointer(&r))
	}
}

// OK() constructs an OK Result from an existing value.
func OK[T any](value T) (r Result[T]) {
	r.v = value
	return
}

// Err() constructs an Err Result from an existing error.
// If the error is nil, an error containing the empty string is returned.
func Err[T any](err error) (r Result[T]) {
	if err == nil {
		r.t = errors.New("")
	} else {
		r.t = err
	}
	return
}

// IsOK checks directly if the Result is OK.
func (r Result[T]) IsOK() bool {
	return r.t == nil
}

// IsErr checks directly if the Result is Err.
func (r Result[T]) IsErr() bool {
	return r.t != nil
}

// UnwrapOr returns the value of the Result if it is OK.
// It returns a static value if the Result is Err.
func (r Result[T]) UnwrapOr(d *T) *T {
	if r.t == nil {
		return &r.v
	} else {
		return d
	}
}

// UnwrapOrElse returns the value of the Result if it is OK.
// It returns a calculated value if the Result is Err.
func (r Result[T]) UnwrapOrElse(f func() *T) *T {
	if r.t == nil {
		return &r.v
	} else {
		return f()
	}
}
