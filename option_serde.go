package safe

import (
	"encoding/json"
)

// We need to handle JSON serde specially.
// The special value `{}` is used to represent None, anything else represents Some.
// BEWARE that Some[T] where T is a struct with no public members, or with "omit" on all members
// MAY also serialize to `{}`, and thereby deserialise back to None.
//
// TODO: handle other serde formats.

// MarshalJSON takes a value receiver for maximum genericity
func (o Option[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.t)
}
func (o Option2[T1, T2]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.t)
}
func (o Option3[T1, T2, T3]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.t)
}
func (o Option4[T1, T2, T3, T4]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.t)
}
func (o Option5[T1, T2, T3, T4, T5]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.t)
}
func (o Option6[T1, T2, T3, T4, T5, T6]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.t)
}
func (o Option7[T1, T2, T3, T4, T5, T6, T7]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.t)
}
func (o Option8[T1, T2, T3, T4, T5, T6, T7, T8]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.t)
}

// UnmarshalJSON takes a pointer receiver for maximum genericity
func (o *optionN[T1, T2, T3, T4, T5, T6, T7, T8]) UnmarshalJSON(b []byte) (err error) {
	var try0 struct{}
	var try1 T1
	var try2 T2
	var try3 T3
	var try4 T4
	var try5 T5
	var try6 T6
	var try7 T7
	var try8 T8
	// This will succeed iff b == "{}", which represents None
	err = json.Unmarshal(b, &try0)
	if err == nil {
		*o = optionN[T1, T2, T3, T4, T5, T6, T7, T8]{}
		return
	}
	// Otherwise we read the value directly
	err = json.Unmarshal(b, &try1)
	if err == nil {
		*o = someN[T1, T1, T2, T3, T4, T5, T6, T7, T8](try1)
		return
	}
	err = json.Unmarshal(b, &try2)
	if err == nil {
		*o = someN[T2, T1, T2, T3, T4, T5, T6, T7, T8](try2)
		return
	}
	err = json.Unmarshal(b, &try3)
	if err == nil {
		*o = someN[T3, T1, T2, T3, T4, T5, T6, T7, T8](try3)
		return
	}
	err = json.Unmarshal(b, &try4)
	if err == nil {
		*o = someN[T4, T1, T2, T3, T4, T5, T6, T7, T8](try4)
		return
	}
	err = json.Unmarshal(b, &try5)
	if err == nil {
		*o = someN[T5, T1, T2, T3, T4, T5, T6, T7, T8](try5)
		return
	}
	err = json.Unmarshal(b, &try6)
	if err == nil {
		*o = someN[T6, T1, T2, T3, T4, T5, T6, T7, T8](try6)
		return
	}
	err = json.Unmarshal(b, &try7)
	if err == nil {
		*o = someN[T7, T1, T2, T3, T4, T5, T6, T7, T8](try7)
		return
	}
	err = json.Unmarshal(b, &try8)
	if err == nil {
		*o = someN[T8, T1, T2, T3, T4, T5, T6, T7, T8](try8)
		return
	}
	return
}

func (o *Option[T1]) UnmarshalJSON(b []byte) (err error) {
	err = (*optionN[T1, xx, xx, xx, xx, xx, xx, xx])(o).UnmarshalJSON(b)
	return
}
func (o *Option2[T1, T2]) UnmarshalJSON(b []byte) (err error) {
	err = (*optionN[T1, T2, xx, xx, xx, xx, xx, xx])(o).UnmarshalJSON(b)
	return
}
func (o *Option3[T1, T2, T3]) UnmarshalJSON(b []byte) (err error) {
	err = (*optionN[T1, T2, T3, xx, xx, xx, xx, xx])(o).UnmarshalJSON(b)
	return
}
func (o *Option4[T1, T2, T3, T4]) UnmarshalJSON(b []byte) (err error) {
	err = (*optionN[T1, T2, T3, T4, xx, xx, xx, xx])(o).UnmarshalJSON(b)
	return
}
func (o *Option5[T1, T2, T3, T4, T5]) UnmarshalJSON(b []byte) (err error) {
	err = (*optionN[T1, T2, T3, T4, T5, xx, xx, xx])(o).UnmarshalJSON(b)
	return
}
func (o *Option6[T1, T2, T3, T4, T5, T6]) UnmarshalJSON(b []byte) (err error) {
	err = (*optionN[T1, T2, T3, T4, T5, T6, xx, xx])(o).UnmarshalJSON(b)
	return
}
func (o *Option7[T1, T2, T3, T4, T5, T6, T7]) UnmarshalJSON(b []byte) (err error) {
	err = (*optionN[T1, T2, T3, T4, T5, T6, T7, xx])(o).UnmarshalJSON(b)
	return
}
func (o *Option8[T1, T2, T3, T4, T5, T6, T7, T8]) UnmarshalJSON(b []byte) (err error) {
	err = (*optionN[T1, T2, T3, T4, T5, T6, T7, T8])(o).UnmarshalJSON(b)
	return
}
