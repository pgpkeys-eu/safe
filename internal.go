package safe

import (
	"fmt"
	"reflect"
	"unsafe"
)

// optionN represents a tagged-union of None and Some.
// The zero literal optionN[T1, T2, T3, T4, T5, T6, T7, T8]{} represents None.
//
// BEWARE that T1 MUST be the type with the largest size.
// If not, someN() will panic on first use (but emit a friendly error).
// Unit testing is your friend.
//
// BEWARE also that if T1 is (or contains) a pointer, reflect.DeepEquals() will attempt to
// dereference that location in memory as a pointer, regardless of the type currently stored.
// This can cause test suites to panic with a pointer error even if no pointers are being compared,
// and DeepEquals() will emit a very UNfriendly error.
type optionN[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
	t any
	v T1
}

// none is a private type with no members, used when we need a non-nil but unusable type.
type none struct{}

// someN[T0, T1, T2, T3, T4, T5, T6, T7, T8] constructs an optionN[T1, T2, T3, T4, T5, T6, T7, T8] from an existing value of type T0.
// If the value is nil it returns None, otherwise Some.
//
// someN will panic with a friendly error if either:
//  1. T1 is not the larger of {T1, T2, T3, T4, T5, T6, T7, T8}
//  2. T0 is not one of T1 or T2
func someN[T0, T1, T2, T3, T4, T5, T6, T7, T8 any](value T0) (o optionN[T1, T2, T3, T4, T5, T6, T7, T8]) {
	// Check static constraints first
	noneType := reflect.TypeFor[none]()
	var t [9]reflect.Type
	t[0], t[1], t[2], t[3], t[4], t[5], t[6], t[7], t[8] = reflect.TypeFor[T0](),
		reflect.TypeFor[T1](), reflect.TypeFor[T2](), reflect.TypeFor[T3](), reflect.TypeFor[T4](),
		reflect.TypeFor[T5](), reflect.TypeFor[T6](), reflect.TypeFor[T7](), reflect.TypeFor[T8]()
	largest := 1
	typeMatched := (t[1] == t[0])
	for n := 2; n < 9 && t[n] != noneType; n++ {
		if t[n].Size() > t[largest].Size() {
			largest = n
		}
		if t[n] == t[0] {
			typeMatched = true
		}
	}
	if largest != 1 {
		panic(fmt.Sprintf("bad type ordering; you must put the largest type (%s) first", t[largest].String()))
	}
	if !typeMatched {
		panic(fmt.Sprintf("type %s is not in [%s, %s, %s, %s, %s, %s, %s, %s]", t[0].String(), t[1].String(), t[2].String(), t[3].String(), t[4].String(), t[5].String(), t[6].String(), t[7].String(), t[8].String()))
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
		return optionN[T1, T2, T3, T4, T5, T6, T7, T8]{}
	} else {
		// pointer mangle o into an Option[T] so we can write its value
		switch t[0] {
		case t[1]:
			m := (*Option[T1])(unsafe.Pointer(&o))
			m.v = *(*T1)(unsafe.Pointer(&value))
			m.t = m.v
		case t[2]:
			m := (*Option[T2])(unsafe.Pointer(&o))
			m.v = *(*T2)(unsafe.Pointer(&value))
			m.t = m.v
		case t[3]:
			m := (*Option[T3])(unsafe.Pointer(&o))
			m.v = *(*T3)(unsafe.Pointer(&value))
			m.t = m.v
		case t[4]:
			m := (*Option[T4])(unsafe.Pointer(&o))
			m.v = *(*T4)(unsafe.Pointer(&value))
			m.t = m.v
		case t[5]:
			m := (*Option[T5])(unsafe.Pointer(&o))
			m.v = *(*T5)(unsafe.Pointer(&value))
			m.t = m.v
		case t[6]:
			m := (*Option[T6])(unsafe.Pointer(&o))
			m.v = *(*T6)(unsafe.Pointer(&value))
			m.t = m.v
		case t[7]:
			m := (*Option[T7])(unsafe.Pointer(&o))
			m.v = *(*T7)(unsafe.Pointer(&value))
			m.t = m.v
		case t[8]:
			m := (*Option[T8])(unsafe.Pointer(&o))
			m.v = *(*T8)(unsafe.Pointer(&value))
			m.t = m.v
		}
		return
	}
}

// match() converts an optionN into a MatchSome or none.
func (o optionN[T1, T2, T3, T4, T5, T6, T7, T8]) match() any {
	// use pointer type mangling to avoid making copies
	// this relies on the memory layouts of Option and MatchSome being identical
	switch o.t.(type) {
	case nil:
		return none{} // check none first for efficiency
	case T1:
		return *(*MatchSome[T1])(unsafe.Pointer(&o))
	case T2:
		return *(*MatchSome[T2])(unsafe.Pointer(&o))
	case T3:
		return *(*MatchSome[T3])(unsafe.Pointer(&o))
	case T4:
		return *(*MatchSome[T4])(unsafe.Pointer(&o))
	case T5:
		return *(*MatchSome[T5])(unsafe.Pointer(&o))
	case T6:
		return *(*MatchSome[T6])(unsafe.Pointer(&o))
	case T7:
		return *(*MatchSome[T7])(unsafe.Pointer(&o))
	case T8:
		return *(*MatchSome[T8])(unsafe.Pointer(&o))
	default:
		return none{} // should never get here but just in case
	}
}
