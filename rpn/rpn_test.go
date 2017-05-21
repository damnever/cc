package rpn

import (
	"testing"

	"github.com/damnever/cc/assert"
)

// TODO(damnever): bad cases

func TestBasicCalculateRPN(t *testing.T) {
	{
		rpn, err := New("N+2")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", "+"})
	}
	{
		rpn, err := New("N-2")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", "-"})
	}
	{
		rpn, err := New("N*2")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", "*"})
	}
	{
		rpn, err := New("N/2")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", "/"})
	}
	{
		rpn, err := New("N%(2+3)")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", "3", "+", "%"})
	}
	{
		rpn, err := New("5+N*2-3")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"5", "N", "2", "*", "+", "3", "-"})
	}
}

func TestBasicConditionRPN(t *testing.T) {
	{
		rpn, err := New("N>2")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", ">"})
	}
	{
		rpn, err := New("N>=2")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", ">="})
	}
	{
		rpn, err := New("N<2")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", "<"})
	}
	{
		rpn, err := New("N<=2")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", "<="})
	}
	{
		rpn, err := New("N==2")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", "=="})
	}
	{
		rpn, err := New("N!=2")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", "!="})
	}
	{
		rpn, err := New("!(N==2)")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", "==", "!"})
	}
	{
		rpn, err := New("(N==2)||(N!=3)")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", "==", "N", "3", "!=", "||"})
	}
	{
		rpn, err := New("!((N!=2)&&(N>=3))")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", "!=", "N", "3", ">=", "&&", "!"})
	}
}

func TestRPN(t *testing.T) {
	{
		rpn, err := New("(N*(N-3)>=10)&&(N<7)&&(N>5)")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "N", "3", "-", "*", "10", ">=", "N", "7", "<", "&&", "N", "5", ">", "&&"})
		res, err := rpn.Calculate(6)
		assert.Must(t, err)
		assert.Check(t, res, true)
		res, err = rpn.Calculate(4)
		assert.Must(t, err)
		assert.Check(t, res, false)
		res, err = rpn.Calculate(5)
		assert.Must(t, err)
		assert.Check(t, res, false)
	}
	{
		rpn, err := New("!(N>3)")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "3", ">", "!"})
		res, err := rpn.Calculate(4)
		assert.Must(t, err)
		assert.Check(t, res, false)
		res, err = rpn.Calculate(3)
		assert.Must(t, err)
		assert.Check(t, res, true)
	}
	{
		rpn, err := New("(N/100>0.3)&&(N/100<=0.8)")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "100", "/", "0.3", ">", "N", "100", "/", "0.8", "<=", "&&"})
		res, err := rpn.Calculate(40)
		assert.Must(t, err)
		assert.Check(t, res, true)
		res, err = rpn.Calculate(80)
		assert.Must(t, err)
		assert.Check(t, res, true)
		res, err = rpn.Calculate(29)
		assert.Must(t, err)
		assert.Check(t, res, false)
		res, err = rpn.Calculate(81)
		assert.Must(t, err)
		assert.Check(t, res, false)
	}
	{
		rpn, err := New("!((N*2>20)||(N<=8&&N%2==0))")
		assert.Must(t, err)
		assertStringList(t, rpn.notation, []string{"N", "2", "*", "20", ">", "N", "8", "<=", "N", "2", "%", "0", "==", "&&", "||", "!"})
		res, err := rpn.Calculate(6)
		assert.Must(t, err)
		assert.Check(t, res, false)
		res, err = rpn.Calculate(9)
		assert.Must(t, err)
		assert.Check(t, res, true)
	}
}

func assertStringList(t *testing.T, l1 []string, l2 []string) {
	if len(l1) != len(l2) {
		t.Fatalf("%v != %v\n", l1, l2)
	}
	for i, s := range l1 {
		if s != l2[i] {
			t.Fatalf("%v != %v at %v\n", l1, l2, i)
		}
	}
}
