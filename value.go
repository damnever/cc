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

// Config returns the value as a Configer, the modification on returned
// Configer has no affect to the origin value.
func (v *Value) Config() Configer {
	switch x := v.v.(type) {
	case Configer:
		val := NewConfig()
		for kx, kv := range x.KV() {
			val.kv[kx] = kv
		}
		return val
	case map[string]interface{}:
		val := NewConfig()
		for kx, kv := range x {
			val.kv[kx] = kv
		}
		return val
	case map[interface{}]interface{}:
		val := NewConfig()
		val.kv = unknownMapToStringMap(x)
		return val
	}
	return NewConfig()
}

// Map returns the value as a map, the modification on returned
// map has no affect to the origin value.
func (v *Value) Map() map[string]Valuer {
	switch x := v.v.(type) {
	case Configer:
		val := x.KV()
		ms := make(map[string]Valuer, len(val))
		for kx, vx := range val {
			ms[kx] = NewValue(vx)
		}
		return ms
	case map[string]interface{}:
		ms := make(map[string]Valuer, len(x))
		for kx, vx := range x {
			ms[kx] = NewValue(vx)
		}
		return ms
	case map[interface{}]interface{}:
		ms := make(map[string]Valuer, len(x))
		for kx, kv := range x {
			ms[fmt.Sprintf("%v", kx)] = NewValue(kv)
		}
		return ms
	}
	return map[string]Valuer{}
}

// List returns the value as slice, the modification on returned
// slice has no affect to the origin value.
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

// StringAndOr returns the string value if pattern matched,
// otherwise returns the deflt.
func (v *Value) StringAndOr(pattern string, deflt string) string {
	if s, ok := v.StringAnd(pattern); ok {
		return s
	}
	return deflt
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

// IntAndOr returns the int value if pattern matched,
// otherwise returns the deflt.
func (v *Value) IntAndOr(pattern string, deflt int) int {
	if n, ok := v.IntAnd(pattern); ok {
		return n
	}
	return deflt
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

// FloatAndOr returns the float64 value if pattern matched,
// otherwise returns the deflt.
func (v *Value) FloatAndOr(pattern string, deflt float64) float64 {
	if n, ok := v.FloatAnd(pattern); ok {
		return n
	}
	return deflt
}

// Duration returns the time.Duration value, returns time.Duration(0) if not exists.
func (v *Value) Duration() time.Duration {
	return v.DurationOr(0)
}

// DurationOr returns the time.Duration value, returns time.Duration(deflt)
// if not exists.
func (v *Value) DurationOr(deflt int) time.Duration {
	return time.Duration(v.IntOr(deflt))
}

// DurationAnd returns the (time.Duration(value), true) if pattern matched,
// otherwise (time.Duration(0), false) returned.
func (v *Value) DurationAnd(pattern string) (time.Duration, bool) {
	n, ok := v.IntAnd(pattern)
	return time.Duration(n), ok
}

// DurationAndOr returns the time.Duration value if pattern matched,
// otherwise returns the deflt.
func (v *Value) DurationAndOr(pattern string, deflt int) time.Duration {
	if d, ok := v.DurationAnd(pattern); ok {
		return d
	}
	return time.Duration(deflt)
}

// GoString implements the native format for Value
func (v *Value) GoString() string {
	return fmt.Sprintf("%v", v.v)
}
