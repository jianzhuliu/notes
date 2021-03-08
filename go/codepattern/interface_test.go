package codepattern

import (
	"testing"
)

func TestInterface(t *testing.T){
	r := Rect{Width: 3.0, Height: 4.0}
	c := Circle{Radius: 2.0}

	s := []Shape{&r, &c}

	for _, sh := range s {
		t.Logf("%#v\n", sh)
		t.Log("Area() == ", sh.Area())
		t.Log("Perimeter() == ", sh.Perimeter())
		t.Log("---------------------")
	}
}