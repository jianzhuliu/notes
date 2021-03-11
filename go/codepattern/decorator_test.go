package codepattern

import (
	"testing"
)

func TestTimedSumFunc(t *testing.T) {
	sum1 := TimedSumFunc(Sum1)
	sum2 := TimedSumFunc(Sum2)

	if r1 := sum1(-10000, 10001); r1 != 10001 {
		t.Fatalf("expect 10001, but get %d\n", r1)
	}

	if r2 := sum2(-10000, 10001); r2 != 10001 {
		t.Fatalf("expect 10001, but get %d\n", r2)
	}
}

func TestTimedFuncDecorator(t *testing.T) {
	var sum1 SumFunc
	_ = TimedFuncDecorator(&sum1, Sum1)

	if r1 := sum1(-10000, 10001); r1 != 10001 {
		t.Fatalf("expect 10001, but get %d\n", r1)
	}

	sum2 := Sum2
	_ = TimedFuncDecorator(&sum2, Sum2)

	if r1 := sum2(-10000, 10002); r1 != 20003 {
		t.Fatalf("expect 20003, but get %d\n", r1)
	}
}
