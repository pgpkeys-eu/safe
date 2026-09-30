package safe

import (
	"reflect"
	"unsafe"
)

// Option uses optionN as its underlying type, but optionN has to handle
// special cases such as type and size mismatches that do not apply to Option.
// Some and Match are therefore implemented separately here and not aliased.

// Option represents a tagged-union of None and Some[T], which defaults to None.
// Option values SHOULD be created by calling either None() or Some().
// DO NOT construct non-zero Option literals directly.
type Option[T any] optionN[T, xx, xx, xx, xx, xx, xx, xx]

// MatchSome is a public match type with one public member Some.
// It MUST be used ONLY for type matching the return value of Match().
// DO NOT construct directly using a struct literal, use Some<n>() to create an Option<n>.
type MatchSome[T any] struct {
	_    any
	Some T
}

// Match() converts an Option into a MatchSome or an unusable type (representing None).
func (o Option[T]) Match() any {
	// use pointer type mangling to avoid making copies
	// this relies on the memory layouts of Option and MatchSome being identical
	switch o.tag.(type) {
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
		o.value = value
		o.tag = o.value
		return
	}
}

// IsNone checks directly if the Option is None.
func (o Option[T]) IsNone() bool {
	return o.tag == nil
}

// IsSome checks directly if the Option is Some.
func (o Option[T]) IsSome() bool {
	return o.tag != nil
}

// UnwrapOr returns the value of the Option if it is Some.
// It returns a static value if the Option is None.
func (o Option[T]) UnwrapOr(d *T) *T {
	if o.tag == nil {
		return d
	}
	return &o.value
}

// UnwrapOrElse returns the value of the Option if it is Some.
// It returns a calculated value if the Option is None.
func (o Option[T]) UnwrapOrElse(f func() *T) *T {
	if o.tag == nil {
		return f()
	}
	return &o.value
}
