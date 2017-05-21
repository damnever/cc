package cc

import (
	"testing"

	"github.com/damnever/cc/assert"
)

func TestPatternValidateString(t *testing.T) {
	{
		p := NewPattern("^123$")
		assert.Check(t, p.ValidateString("123"), true)
		assert.Check(t, p.ValidateString("x"), false)
	}
	{
		p := NewPattern("^[^12&%*")
		assert.Check(t, p.ValidateString("123"), false)
		if p.Err() == nil {
			t.Fatal("expect error, got nothing")
		}
		assert.Check(t, p.badRe, true)
	}
}

func TestPatternValidateInt(t *testing.T) {
	{
		p := NewPattern("N>20&&N<=80")
		assert.Check(t, p.ValidateInt(40), true)
		assert.Check(t, p.ValidateInt(20), false)
	}
	{
		p := NewPattern("N>20&&N<=")
		assert.Check(t, p.ValidateInt(40), false)
		if p.Err() == nil {
			t.Fatal("expect error, got nothing")
		}
		assert.Check(t, p.badCond, true)
	}
}

func TestPatternValidateFloat(t *testing.T) {
	{
		p := NewPattern("N/100>0.2&&N/100<=0.8")
		assert.Check(t, p.ValidateFloat(40), true)
		assert.Check(t, p.ValidateFloat(20), false)
	}
	{
		p := NewPattern("N/100>0.0.20&&N/100<=0.8")
		assert.Check(t, p.ValidateFloat(40), false)
		if p.Err() == nil {
			t.Fatal("expect error, got nothing")
		}
		assert.Check(t, p.badCond, true)
	}

}
