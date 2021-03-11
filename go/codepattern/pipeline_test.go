package codepattern

import (
	"reflect"
	"testing"
)

var nums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func TestPipeline(t *testing.T) {
	out := make([]int, 0, len(nums))
	for n := range Pipeline(nums, Echo, Odd, Sqrt) {
		out = append(out, n)
	}

	expect := []int{1, 9, 25, 49, 81}
	if !reflect.DeepEqual(out, expect) {
		t.Fatalf("should be %v, but got %v", expect, out)
	}
}

func TestEcho(t *testing.T) {
	out := make([]int, 0, len(nums))
	for n := range Echo(nums) {
		out = append(out, n)
	}

	if !reflect.DeepEqual(out, nums) {
		t.Fatalf("should be %v, but got %v", nums, out)
	}
}

func TestSqrt(t *testing.T) {
	out := make([]int, 0, len(nums))

	for n := range Sqrt(Echo(nums)) {
		out = append(out, n)
	}

	expect := []int{1, 4, 9, 16, 25, 36, 49, 64, 81, 100}
	if !reflect.DeepEqual(out, expect) {
		t.Fatalf("should be %v, but got %v", expect, out)
	}
}

func TestOdd(t *testing.T) {
	out := make([]int, 0, len(nums))

	for n := range Odd(Echo(nums)) {
		out = append(out, n)
	}

	expect := []int{1, 3, 5, 7, 9}
	if !reflect.DeepEqual(out, expect) {
		t.Fatalf("should be %v, but got %v", expect, out)
	}
}

func TestSum(t *testing.T) {
	out := <-Sum(Odd(Echo(nums)))
	expect := 25
	if out != expect {
		t.Fatalf("should be %v, but got %v", expect, out)
	}
}
