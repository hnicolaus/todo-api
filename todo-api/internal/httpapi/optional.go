package httpapi

import (
	"bytes"
	"encoding/json"
)

// OptionalString lets PATCH distinguish between:
// - field absent (Set=false)
// - field present with string value (Set=true, Value!=nil)
// - field present with null (Set=true, Value=nil)
type OptionalString struct {
	Set   bool
	Value *string
}

func (o *OptionalString) UnmarshalJSON(b []byte) error {
	o.Set = true
	if bytes.Equal(b, []byte("null")) {
		o.Value = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	o.Value = &s
	return nil
}
