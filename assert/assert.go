package assert

import (
	"runtime"
	"testing"
)

func Must(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Check(t *testing.T, actual interface{}, expect interface{}) {
	_, fileName, line, _ := runtime.Caller(1)
	if actual != expect {
		t.Fatalf("expect %v, got %v at (%v:%v)\n", expect, actual, fileName, line)
	}
}
