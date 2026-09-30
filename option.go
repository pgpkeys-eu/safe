package safe

import (
	"reflect"
	"unsafe"
)

// Option represents a tagged-union of None and Some.
// The zero literal Option[T]{} represents None.
type Option[T any] struct {
	//
	// The zero literal Option[T]{} may be safely used
	// as the static equivalent of None[T]()
	//
	// For everything else, use Some[T]()
	//
	// DO NOT SET STRUCT MEMBER VALUES DIRECTLY
	//
	t any
	v T
}

// none is a private type with no members, used when we need a non-nil but unusable type.
type none struct{}

// MatchSome is a match type with one public member Some, used only for type matching.
// DO NOT construct directly using a struct literal, use Some() instead.
type MatchSome[T any] struct {
	//
	// !!!!! DO NOT USE THIS STRUCT LITERAL !!!!!
	//
	// It is not safe to directly create MatchSome values
	// Use Some[T]() instead
	//
	// !!!!! DO NOT USE THIS STRUCT LITERAL !!!!!
	//
	_    any
	Some T
}

// Match() converts an Option into a MatchSome or none.
func (o Option[T]) Match() any {
	// use pointer type mangling to avoid making copies
	// this relies on the memory layouts of Option and MatchSome being identical
	switch o.t.(type) {
	case T:
		return *(*MatchSome[T])(unsafe.Pointer(&o))
	default:
		return none{}
	}
}

// None() constructs a new None Option
func None[T any]() Option[T] {
	return Option[T]{}
}

// Some() constructs an Option from an existing value.
// If the value is nil it returns None, otherwise Some.
func Some[T any](value T) (o Option[T]) {
	v := reflect.ValueOf(value)
	kind := v.Kind()
	if (kind == reflect.Ptr ||
		kind == reflect.Interface ||
		kind == reflect.Slice ||
		kind == reflect.Map ||
		kind == reflect.Chan ||
		kind == reflect.Func) && v.IsNil() {
		return Option[T]{}
	} else {
		o.v = value
		o.t = o.v
		return
	}
}

// IsNone checks directly if the Option is None.
func (o Option[T]) IsNone() bool {
	return o.t == nil
}

// IsSome checks directly if the Option is Some.
func (o Option[T]) IsSome() bool {
	return o.t != nil
}

// UnwrapOr returns the value of the Option if it is Some.
// It returns a static value if the Option is None.
func (o Option[T]) UnwrapOr(d *T) *T {
	if o.t == nil {
		return d
	}
	return &o.v
}

// UnwrapOrElse returns the value of the Option if it is Some.
// It returns a calculated value if the Option is None.
func (o Option[T]) UnwrapOrElse(f func() *T) *T {
	if o.t == nil {
		return f()
	}
	return &o.v
}
