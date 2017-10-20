package rpn

import (
	"container/list"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var numbers = "0123456789"
var priorities = map[string]uint8{
	"||": 1,
	"&&": 1,
	">":  2,
	"<":  2,
	">=": 2,
	"<=": 2,
	"==": 2,
	"!=": 2,
	"+":  3,
	"-":  3,
	"%":  4,
	"*":  4,
	"/":  4,
}

// ReversePolishNotation represents a reverse polish notation.
type ReversePolishNotation struct {
	notation []string
}

// New creates a new ReversePolishNotation with a string pattern.
func New(s string) (*ReversePolishNotation, error) {
	nop, n := 0, len(s)
	operators := make([]string, n)
	notation := []string{}

	i := 0
	popPushOp := func(op string) {
		priority := priorities[op]
		for nop > 0 && priorities[operators[nop-1]] >= priority {
			nop--
			notation = append(notation, operators[nop])
		}
		operators[nop] = op
		nop++
		i++
	}

	for i < n {
		c := s[i]
		switch c {
		case ' ':
			i++
		case ')':
			for nop > 0 && operators[nop-1] != "(" {
				nop--
				notation = append(notation, operators[nop])
			}
			if nop == 0 || operators[nop-1] != "(" {
				return nil, fmt.Errorf("'%v' has no '(' found for ')' at %v", s, i)
			}
			nop--
			if nop > 0 && operators[nop-1] == "!" {
				notation = append(notation, "!")
				nop--
			}
			i++
		case '(':
			operators[nop] = "("
			nop++
			i++
		case '*', '/', '%', '+', '-':
			popPushOp(string(c))
		case '!':
			next := s[i+1]
			if next == '(' {
				operators[nop] = string(c)
				nop++
				i++
			} else if next == '=' {
				i++
				popPushOp(string([]byte{c, next}))
			} else {
				return nil, fmt.Errorf("'%v' has invalid token at %v: %v", s, i+1, next)
			}
		case '>', '<':
			op := []byte{c}
			if s[i+1] == '=' {
				op = append(op, '=')
				i++
			}
			popPushOp(string(op))
		case '|', '&', '=':
			next := s[i+1]
			if next != c {
				return nil, fmt.Errorf("'%v' has invalid token at %v: %v", s, i+1, next)
			}
			i++
			popPushOp(string([]byte{c, next}))
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			hasDot := false
			num := []byte{}
			for i < n {
				c = s[i]
				if c == '.' {
					if hasDot || len(num) == 0 {
						return nil, fmt.Errorf("'%v' has invalid token at %v: %v", s, i, c)
					}
					hasDot = true
				} else if !strings.Contains(numbers, string(c)) {
					break
				}
				num = append(num, c)
				i++
			}
			notation = append(notation, string(num))
		case 'N':
			notation = append(notation, "N")
			i++
		default:
			return nil, fmt.Errorf("'%v' has invalid token at %v: %v", s, i, c)
		}
	}

	for nop > 0 {
		nop--
		if op := operators[nop]; op != "(" {
			notation = append(notation, op)
		}
	}

	return &ReversePolishNotation{notation: notation}, nil
}

var (
	errInvalidExpr  = errors.New("invalid expression")
	errDivideByZero = errors.New("invalid expression, divide by zero")
	errInvalidCond  = errors.New("invalid condition")
)

// Calculate calculate condition result with float64.
func (rpn *ReversePolishNotation) Calculate(value float64) (bool, error) {
	values := list.New()

	for i, n := 0, len(rpn.notation); i < n; i++ {
		op := rpn.notation[i]
		switch op {
		case "+":
			if values.Len() < 2 {
				return false, errInvalidExpr
			}
			num1, num2, ok := lastTwoNum(values)
			if !ok {
				return false, errInvalidExpr
			}
			values.PushBack((num1 + num2))
		case "-":
			if values.Len() < 2 {
				return false, errInvalidExpr
			}
			num1, num2, ok := lastTwoNum(values)
			if !ok {
				return false, errInvalidExpr
			}
			values.PushBack((num1 - num2))
		case "*":
			if values.Len() < 2 {
				return false, errInvalidExpr
			}
			num1, num2, ok := lastTwoNum(values)
			if !ok {
				return false, errInvalidExpr
			}
			values.PushBack((num1 * num2))
		case "/":
			if values.Len() < 2 {
				return false, errInvalidExpr
			}
			num1, num2, ok := lastTwoNum(values)
			if !ok {
				return false, errInvalidExpr
			}
			if num2 == 0.0 {
				return false, errDivideByZero
			}
			values.PushBack((num1 / num2))
		case "%":
			if values.Len() < 2 {
				return false, errInvalidExpr
			}
			num1, num2, ok := lastTwoNum(values)
			if !ok {
				return false, errInvalidExpr
			}
			if num2 == 0.0 {
				return false, errDivideByZero
			}
			values.PushBack(float64(int(num1) % int(num2)))
		case "<":
			if values.Len() < 2 {
				return false, errInvalidCond
			}
			num1, num2, ok := lastTwoNum(values)
			if !ok {
				return false, errInvalidCond
			}
			values.PushBack((num1 < num2))
		case ">":
			if values.Len() < 2 {
				return false, errInvalidCond
			}
			num1, num2, ok := lastTwoNum(values)
			if !ok {
				return false, errInvalidCond
			}
			values.PushBack((num1 > num2))
		case "<=":
			if values.Len() < 2 {
				return false, errInvalidCond
			}
			num1, num2, ok := lastTwoNum(values)
			if !ok {
				return false, errInvalidCond
			}
			values.PushBack((num1 <= num2))
		case ">=":
			if values.Len() < 2 {
				return false, errInvalidCond
			}
			num1, num2, ok := lastTwoNum(values)
			if !ok {
				return false, errInvalidCond
			}
			values.PushBack((num1 >= num2))
		case "==": // XXX: may someone compare bool values?
			if values.Len() < 2 {
				return false, errInvalidCond
			}
			num1, num2, ok := lastTwoNum(values)
			if !ok {
				return false, errInvalidCond
			}
			values.PushBack((num1 == num2))
		case "!=": // XXX: may someone compare bool values?
			if values.Len() < 2 {
				return false, errInvalidCond
			}
			num1, num2, ok := lastTwoNum(values)
			if !ok {
				return false, errInvalidCond
			}
			values.PushBack((num1 != num2))
		case "||":
			if values.Len() < 2 {
				return false, errInvalidCond
			}
			b1, b2, ok := lastTwoBool(values)
			if !ok {
				return false, errInvalidCond
			}
			values.PushBack((b1 || b2))
		case "&&":
			if values.Len() < 2 {
				return false, errInvalidCond
			}
			b1, b2, ok := lastTwoBool(values)
			if !ok {
				return false, errInvalidCond
			}
			values.PushBack((b1 && b2))
		case "!":
			if values.Len() < 1 {
				return false, errInvalidCond
			}
			e := values.Back()
			b, ok := e.Value.(bool)
			if !ok {
				return false, errInvalidCond
			}
			values.Remove(e)
			values.PushBack((!b))
		case "N":
			values.PushBack(value)
		default:
			num, err := strconv.ParseFloat(op, 64)
			if err != nil {
				return false, fmt.Errorf("invalid number: %v", err)
			}
			values.PushBack(num)
		}
	}

	if values.Len() != 1 {
		return false, errInvalidCond
	}
	e := values.Back()
	if res, ok := e.Value.(bool); ok {
		return res, nil
	}
	return false, errInvalidCond
}

func lastTwoNum(values *list.List) (num1 float64, num2 float64, ok bool) {
	e2 := values.Back()
	if num2, ok = e2.Value.(float64); !ok {
		return
	}
	values.Remove(e2)

	e1 := values.Back()
	if num1, ok = e1.Value.(float64); !ok {
		return
	}
	values.Remove(e1)
	return
}

func lastTwoBool(values *list.List) (v1 bool, v2 bool, ok bool) {
	e2 := values.Back()
	if v2, ok = e2.Value.(bool); !ok {
		return
	}
	values.Remove(e2)

	e1 := values.Back()
	if v1, ok = e1.Value.(bool); !ok {
		return
	}
	values.Remove(e1)
	return
}
