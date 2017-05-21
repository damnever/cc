package cc

import (
	"fmt"
	"reflect"
	"time"
)

// TODO(damnever): cache used value

// Value implements the Valuer interface.
type Value struct {
	v interface{}
}

// NewValue creates a new Value.
func NewValue(v interface{}) *Value {
	return &Value{v: v}
}

// Exist returns true is value is a valid value, otherwise false.
func (v *Value) Exist() bool {
	if v.v == nil {
		return false
	}
	return reflect.ValueOf(v.v).IsValid()
}

// Pattern returns a Patterner.
func (v *Value) Pattern() Patterner {
	return NewPattern(v.String())
}

// Raw returns the raw value.
func (v *Value) Raw() interface{} {
	return v.v
}

// Map returns the value as a map.
func (v *Value) Map() map[string]Valuer {
	if x, ok := v.v.(map[string]interface{}); ok {
		ms := make(map[string]Valuer, len(x))
		for kx, vx := range x {
			ms[kx] = NewValue(vx)
		}
		return ms
	}
	if x, ok := v.v.(map[interface{}]interface{}); ok {
		ms := make(map[string]Valuer, len(x))
		for kx, vx := range x {
			ms[fmt.Sprintf("%v", kx)] = NewValue(vx)
		}
		return ms
	}
	return map[string]Valuer{}
}

// List returns the value as slice.
func (v *Value) List() []Valuer {
	if x, ok := v.v.([]interface{}); ok {
		vs := make([]Valuer, len(x))
		for i, e := range x {
			vs[i] = NewValue(e)
		}
		return vs
	}
	return []Valuer{}
}

// String returns the string value, returns "" if not exists.
func (v *Value) String() string {
	return v.StringOr("")
}

// StringOr returns the string value, returns the deflt if not exists.
func (v *Value) StringOr(deflt string) string {
	return toString(v.v, deflt)
}

// StringAnd returns the (string value, true) if pattern matched,
// otherwise returns ("", false).
func (v *Value) StringAnd(pattern string) (string, bool) {
	if !v.Exist() {
		return "", false
	}
	p := NewPattern(pattern)
	if s := v.String(); p.ValidateString(s) {
		return s, true
	}
	return "", false
}

// Bool returns the bool value, returns false if not exists.
func (v *Value) Bool() bool {
	return v.BoolOr(false)
}

// BoolOr returns the bool value, returns the deflt if not exists.
func (v *Value) BoolOr(deflt bool) bool {
	return toBool(v.v, deflt)
}

// Int returns the int value, returns 0 if not exists.
func (v *Value) Int() int {
	return v.IntOr(0)
}

// IntOr returns the int value, returns the deflt if not exists.
func (v *Value) IntOr(deflt int) int {
	return toInt(v.v, deflt)
}

// IntAnd returns the (int value, true) if pattern matched,
// otherwise returns (0, false).
func (v *Value) IntAnd(pattern string) (int, bool) {
	if !v.Exist() {
		return 0, false
	}
	p := NewPattern(pattern)
	if n := v.Int(); p.ValidateInt(n) {
		return n, true
	}
	return 0, false
}

// Float returns the float64 value, returns 0.0 if not exists.
func (v *Value) Float() float64 {
	return v.FloatOr(0.0)
}

// FloatOr returns the float64 value, return the deflt if not exists.
func (v *Value) FloatOr(deflt float64) float64 {
	return toFloat64(v.v, deflt)
}

// FloatAnd returns the (float64 value, true) if pattern matched,
// otherwise returns (0.0, false).
func (v *Value) FloatAnd(pattern string) (float64, bool) {
	if !v.Exist() {
		return 0.0, false
	}
	p := NewPattern(pattern)
	if n := v.Float(); p.ValidateFloat(n) {
		return n, true
	}
	return 0.0, false
}

// Duration returns the time.Duration value, returns time.Duration(0) if not exists.
func (v *Value) Duration() time.Duration {
	return v.DurationOr(0)
}

// Duration returns the time.Duration value, returns time.Duration(deflt)
// if not exists.
func (v *Value) DurationOr(deflt int64) time.Duration {
	return time.Duration(toInt64(v.v, deflt))
}

// GoString implements the native format for Value
func (v *Value) GoString() string {
	return fmt.Sprintf("%v", v.v)
}
