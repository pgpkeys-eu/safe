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
	if o.tagField {
		return json.Marshal(o.unionField)
	} else {
		return []byte("{}"), nil
	}
}

// UnmarshalJSON takes a pointer receiver for maximum genericity
func (o *Option[T]) UnmarshalJSON(b []byte) (err error) {
	var try struct{}
	var value T
	// This will succeed iff b == "{}", which represents None
	err = json.Unmarshal(b, &try)
	if err == nil {
		*o = Option[T]{}
		return
	}
	// Otherwise we read the value directly
	err = json.Unmarshal(b, &value)
	if err == nil {
		*o = Option[T]{true, value}
		return
	}
	return
}
