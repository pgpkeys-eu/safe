package safe

import (
	"reflect"
)

// Option uses optionN as its underlying type, but optionN has to handle
// special cases such as type and size mismatches that do not apply to Option.
// Some and Unwrap are therefore implemented separately here and not aliased.

// Option represents a tagged-union of None and T, which defaults to None.
// Option values SHOULD be created by calling either None() or Some().
// DO NOT construct non-zero Option literals directly.
type Option[T any] optionN[T, xx, xx, xx, xx, xx, xx, xx]

// Unwrap() converts an Option into a T or nil.
func (o Option[T]) Unwrap() any {
	return o.tag
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
