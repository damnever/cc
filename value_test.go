package cc

import (
	"fmt"
	"testing"
	"time"

	"github.com/damnever/cc/assert"
)

func TestValueExists(t *testing.T) {
	v := NewValue(nil)
	assert.Check(t, v.Exist(), false)
	v = NewValue(0)
	assert.Check(t, v.Exist(), true)
}

func TestValueToPattern(t *testing.T) {
	v := NewValue("")
	if x, ok := v.Pattern().(Patterner); !ok {
		t.Fatalf("expected Patterner, got %#v\n", x)
	}
}

func TestValueToMap(t *testing.T) {
	v := NewValue(map[string]interface{}{"foo": "bar"})
	m := v.Map()
	assert.Check(t, len(m), 1)
	assert.Check(t, m["foo"].String(), "bar")
}

func TestValueToList(t *testing.T) {
	v := NewValue([]interface{}{"bar", "baz"})
	l := v.List()
	assert.Check(t, len(l), 2)
	assert.Check(t, l[0].String(), "bar")
	assert.Check(t, l[1].String(), "baz")
}

func TestValueToString(t *testing.T) {
	v := NewValue("wow")
	assert.Check(t, v.String(), "wow")
	assert.Check(t, v.StringOr("xx"), "wow")
	res, ok := v.StringAnd("^w")
	assert.Check(t, ok, true)
	assert.Check(t, res, "wow")
	res, ok = v.StringAnd("^o")
	assert.Check(t, ok, false)
	assert.Check(t, res, "")

	v = NewValue(1)
	assert.Check(t, v.String(), "")
	assert.Check(t, v.StringOr("bad"), "bad")
}

func TestValueToBool(t *testing.T) {
	v := NewValue(true)
	assert.Check(t, v.Bool(), true)
	assert.Check(t, v.BoolOr(false), true)

	v = NewValue(false)
	assert.Check(t, v.Bool(), false)
	assert.Check(t, v.BoolOr(true), false)

	v = NewValue("")
	assert.Check(t, v.Bool(), false)
	assert.Check(t, v.BoolOr(true), true)
}

func TestValueToInt(t *testing.T) {
	v := NewValue(1)
	assert.Check(t, v.Int(), 1)
	assert.Check(t, v.IntOr(2), 1)
	res, ok := v.IntAnd("N>0")
	assert.Check(t, ok, true)
	assert.Check(t, res, 1)

	v = NewValue(1.0)
	assert.Check(t, v.Int(), 1)
	assert.Check(t, v.IntOr(2), 1)
	res, ok = v.IntAnd("N<0")
	assert.Check(t, ok, false)
	assert.Check(t, res, 0)

	v = NewValue("")
	assert.Check(t, v.Int(), 0)
	assert.Check(t, v.IntOr(1), 1)
}

func TestValueToFloat(t *testing.T) {
	v := NewValue(3.0)
	assert.Check(t, v.Float(), 3.0)
	assert.Check(t, v.FloatOr(4), 3.0)
	res, ok := v.FloatAnd("N>=3.0")
	assert.Check(t, ok, true)
	assert.Check(t, res, 3.0)

	v = NewValue(3)
	assert.Check(t, v.Float(), 3.0)
	assert.Check(t, v.FloatOr(4), 3.0)

	v = NewValue("")
	assert.Check(t, v.Float(), 0.0)
	assert.Check(t, v.FloatOr(0.5), 0.5)
}

func TestValueToDuration(t *testing.T) {
	v := NewValue(23)
	assert.Check(t, v.Duration(), time.Duration(23))
	assert.Check(t, v.DurationOr(32), time.Duration(23))

	v = NewValue("")
	assert.Check(t, v.Duration(), time.Duration(0))
	assert.Check(t, v.DurationOr(32), time.Duration(32))
}

func TestValueGoString(t *testing.T) {
	v := NewValue(12345)
	s := fmt.Sprintf("%#v", v)
	assert.Check(t, s, "12345")
}
