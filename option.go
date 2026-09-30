package safe

import (
	"reflect"
	"unsafe"
)

// Option uses optionN as its underlying type, but optionN has to handle
// special cases such as type and size mismatches that do not apply to Option.
// Some and Match are therefore implemented separately here and not aliased.

// Option represents a tagged-union of None and Some.
type Option[T any] optionN[T, xx, xx, xx, xx, xx, xx, xx]

// MatchSome is a public match type with one public member Some, used only for type matching.
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

// Match() converts an Option into a MatchSome or an unusable type (representing None).
func (o Option[T]) Match() any {
	// use pointer type mangling to avoid making copies
	// this relies on the memory layouts of Option and MatchSome being identical
	switch o.t.(type) {
	case T:
		return *(*MatchSome[T])(unsafe.Pointer(&o))
	default:
		return xx{}
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
