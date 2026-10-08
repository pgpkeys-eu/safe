package safe

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// optionN represents a tagged-union of None and N other types.
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
	tag   any
	value T1
}

// xx is a private type with no members, used when we need a non-nil but unusable type.
// The name is chosen purely for its visual distinctiveness, particularly in large blocks of boilerplate.
type xx struct{}

// someN[T0, T1, T2, T3, T4, T5, T6, T7, T8] constructs an optionN[T1, T2, T3, T4, T5, T6, T7, T8] from an existing value of type T0.
// If the value is nil it returns None, otherwise Some.
//
// someN will panic with a friendly error if either:
//  1. T1 is not the larger of {T1, T2, T3, T4, T5, T6, T7, T8}
//  2. T0 is not one of T1 or T2
func someN[T0, T1, T2, T3, T4, T5, T6, T7, T8 any](value T0) (o optionN[T1, T2, T3, T4, T5, T6, T7, T8]) {
	// Check static constraints first
	xxType := reflect.TypeFor[xx]()
	var t [9]reflect.Type
	t[0], t[1], t[2], t[3], t[4], t[5], t[6], t[7], t[8] = reflect.TypeFor[T0](),
		reflect.TypeFor[T1](), reflect.TypeFor[T2](), reflect.TypeFor[T3](), reflect.TypeFor[T4](),
		reflect.TypeFor[T5](), reflect.TypeFor[T6](), reflect.TypeFor[T7](), reflect.TypeFor[T8]()
	largest := 1
	typeMatched := (t[1] == t[0])
	// Scan the type parameters until we find xx; that's how we know which concrete Some<n> we were called as
	var n int
	for n = 2; n < 9 && t[n] != xxType; n++ {
		if t[n].Size() > t[largest].Size() {
			largest = n
		}
		if t[n] == t[0] {
			typeMatched = true
		}
	}
	// Panic if static constraints were breached
	if largest != 1 {
		panic(fmt.Sprintf("bad type ordering; you must put the largest type (%s) first", t[largest].String()))
	}
	if !typeMatched {
		var names [8]string
		for i := 0; i < n; i++ {
			names[i] = t[i].String()
		}
		panic(fmt.Sprintf("type %s is not in [%s]", names[0], strings.Join(names[1:n], ", ")))
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
		// pointer mangle o and value so that the types match
		switch t[0] {
		case t[1]:
			// no need to pointer mangle o, o.value is T1 by default
			o.value = *(*T1)(unsafe.Pointer(&value))
			o.tag = o.value
		case t[2]:
			m := (*Option[T2])(unsafe.Pointer(&o))
			m.value = *(*T2)(unsafe.Pointer(&value))
			m.tag = m.value
		case t[3]:
			m := (*Option[T3])(unsafe.Pointer(&o))
			m.value = *(*T3)(unsafe.Pointer(&value))
			m.tag = m.value
		case t[4]:
			m := (*Option[T4])(unsafe.Pointer(&o))
			m.value = *(*T4)(unsafe.Pointer(&value))
			m.tag = m.value
		case t[5]:
			m := (*Option[T5])(unsafe.Pointer(&o))
			m.value = *(*T5)(unsafe.Pointer(&value))
			m.tag = m.value
		case t[6]:
			m := (*Option[T6])(unsafe.Pointer(&o))
			m.value = *(*T6)(unsafe.Pointer(&value))
			m.tag = m.value
		case t[7]:
			m := (*Option[T7])(unsafe.Pointer(&o))
			m.value = *(*T7)(unsafe.Pointer(&value))
			m.tag = m.value
		case t[8]:
			m := (*Option[T8])(unsafe.Pointer(&o))
			m.value = *(*T8)(unsafe.Pointer(&value))
			m.tag = m.value
		}
		return
	}
}
