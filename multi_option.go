package safe

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Option2 represents a tagged-union of None and Some.
// The zero literal Option2[T1, T2]{} represents None.
//
// BEWARE that T1 MUST be the type with the largest size.
// If not, Some2() will panic on first use (but emit a friendly error).
// Unit testing is your friend.
//
// BEWARE also that if T1 is (or contains) a pointer, reflect.DeepEquals() will attempt to
// dereference that location in memory as a pointer, regardless of the type currently stored.
// This can cause test suites to panic with a pointer error even if no pointers are being compared,
// and DeepEquals() will emit a very UNfriendly error.
type Option2[T1, T2 any] struct {
	//
	// The zero literal Option[T1, T2]{} may be safely used
	// as the static equivalent of None2[T1, T2]()
	//
	// For everything else, use Some2[T1, T2](Some[T]())
	//
	// DO NOT SET STRUCT MEMBER VALUES DIRECTLY
	//
	t any
	_ T1
}

// None2() constructs a new None Option2
func None2[T1, T2 any]() Option2[T1, T2] {
	return Option2[T1, T2]{}
}

// Some2[T0, T1, T2] constructs an Option2[T1, T2] from an existing value of type T0.
// If the value is nil it returns None, otherwise Some.
//
// Some2 will panic with a friendly error if either:
//  1. T1 is not the larger of {T1, T2}
//  2. T0 is not one of T1 or T2
func Some2[T0, T1, T2 any](value T0) (o Option2[T1, T2]) {
	// Check static constraints first
	t0, t1, t2 := reflect.TypeFor[T0](), reflect.TypeFor[T1](), reflect.TypeFor[T2]()
	if t2.Size() > t1.Size() {
		panic(fmt.Sprintf("bad type ordering; you must put %s (the largest) first", t2.String()))
	}
	switch t0 {
	case t1, t2:
	default:
		panic(fmt.Sprintf("type %s is not in [%s, %s]", t0.String(), t1.String(), t2.String()))
	}
	// OK, we can continue now
	v := reflect.ValueOf(value)
	kind := v.Kind()
	if (kind == reflect.Ptr ||
		kind == reflect.Interface ||
		kind == reflect.Slice ||
		kind == reflect.Map ||
		kind == reflect.Chan ||
		kind == reflect.Func) && v.IsNil() {
		return Option2[T1, T2]{}
	} else {
		// pointer mangle o into an Option[T] so we can write its value
		switch t0 {
		case t1:
			m := (*Option[T1])(unsafe.Pointer(&o))
			m.v = *(*T1)(unsafe.Pointer(&value))
			m.t = m.v
		case t2:
			m := (*Option[T2])(unsafe.Pointer(&o))
			m.v = *(*T2)(unsafe.Pointer(&value))
			m.t = m.v
		}
		return
	}
}

// Match() converts an Option2 into a MatchSome or MatchNone.
func (o Option2[T1, T2]) Match() any {
	// use pointer type mangling to avoid making copies
	// this relies on the memory layouts of all types being identical
	switch o.t.(type) {
	case T1:
		return *(*MatchSome[T1])(unsafe.Pointer(&o))
	case T2:
		return *(*MatchSome[T2])(unsafe.Pointer(&o))
	default:
		return MatchNone{}
	}
}

// IsNone checks directly if the Option2 is None.
func (o Option2[T1, T2]) IsNone() bool {
	return o.t == nil
}

// IsSome checks directly if the Option2 is Some (of any kind).
// More complex tests should be done using type matching.
func (o Option2[T1, T2]) IsSome() bool {
	return o.t != nil
}
