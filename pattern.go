package cc

import (
	"regexp"

	"github.com/damnever/cc/rpn"
)

// Pattern implements the Patterner interface.
type Pattern struct {
	pattern string
	re      *regexp.Regexp
	rpn     *rpn.ReversePolishNotation
	err     error
	badRe   bool
	badCond bool
}

// NewPattern creates a new Pattern, even if pattern is not valid,
// in such case, Validate-like methods always return false.
// The compiled pattern is cached.
func NewPattern(pattern string) *Pattern {
	return &Pattern{
		pattern: pattern,
		re:      nil,
		rpn:     nil,
		err:     nil,
		badRe:   false,
		badCond: false,
	}
}

// Err returns the error if pattern is wrong.
func (p *Pattern) Err() error {
	return p.err
}

// ValidateInt validate the int value n, return true if it is valid.
func (p *Pattern) ValidateInt(n int) bool {
	return p.ValidateFloat(float64(n))
}

// ValidateFloat validate the float64 value n, return true if it is valid.
func (p *Pattern) ValidateFloat(n float64) bool {
	if p.badCond {
		return false
	}

	if p.rpn == nil {
		rpn, err := rpn.New(p.pattern)
		if err != nil {
			p.err = err
			p.badCond = true
			return false
		}
		p.rpn = rpn
	}
	res, err := p.rpn.Calculate(n)
	if err != nil {
		p.err = err
		p.badCond = true
	}
	return res
}

// ValidateString validate the string value n, return true if it is valid.
func (p *Pattern) ValidateString(s string) bool {
	if p.badRe {
		return false
	}

	if p.re == nil {
		re, err := regexp.Compile(p.pattern)
		if err != nil {
			p.err = err
			p.badRe = true
			return false
		}
		p.re = re
	}
	return p.re.MatchString(s)
}
